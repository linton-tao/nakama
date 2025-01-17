package jumpgo

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/rtapi"
	"github.com/heroiclabs/nakama-common/runtime"
	jumpgoApi "github.com/heroiclabs/nakama/v3/jumpgo/api"
	"github.com/heroiclabs/nakama/v3/jumpgo/model"
	"github.com/heroiclabs/nakama/v3/jumpgo/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	//rand.Seed(time.Now().UnixNano())

	if err := initializer.RegisterRpc("go_echo_sample", rpcEcho); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("init_card", rpcInitCard); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("rpc_create_match", rpcCreateMatch); err != nil {
		return err
	}
	if err := initializer.RegisterBeforeRt("ChannelJoin", beforeChannelJoin); err != nil {
		return err
	}
	if err := initializer.RegisterAfterGetAccount(afterGetAccount); err != nil {
		return err
	}
	if err := initializer.RegisterMatch("match", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &Match{}, nil
	}); err != nil {
		return err
	}
	if err := initializer.RegisterEventSessionStart(eventSessionStart); err != nil {
		return err
	}
	if err := initializer.RegisterEventSessionEnd(eventSessionEnd); err != nil {
		return err
	}
	if err := initializer.RegisterEvent(func(ctx context.Context, logger runtime.Logger, evt *api.Event) {
		logger.Info("Received event: %+v", evt)
	}); err != nil {
		return err
	}

	return nil
}

func rpcInitCard(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params string) (string, error) {
	req := &jumpgoApi.RpcInitMacthRequest{}
	if err := json.Unmarshal([]byte(params), req); err != nil {
		return "", err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return "gormferr", nil
	}
	resp, err := service.InitCard(ctx, logger, nk, req.MatchType, gormDB)
	res, err := json.Marshal(resp)
	return string(res), err
}
func rpcEcho(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return "gormferr", nil
	}

	list := []model.JgCard{}
	err = gormDB.Find(&list).Error
	logger.Info("RUNNING IN GO")
	res, err := json.Marshal(list)
	return string(res), nil
}

func rpcCreateMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	matchID, err := nk.MatchCreate(ctx, "match", map[string]interface{}{})

	if err != nil {
		return "", err
	}

	return matchID, nil
}

func beforeChannelJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, envelope *rtapi.Envelope) (*rtapi.Envelope, error) {
	logger.Info("Intercepted request to join channel '%v'", envelope.GetChannelJoin().Target)
	return envelope, nil
}

func afterGetAccount(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, in *api.Account) error {
	logger.Info("Intercepted response to get account '%v'", in)
	return nil
}

type MatchState struct {
	debug bool
}

type Match struct{}

func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	var debug bool
	if d, ok := params["debug"]; ok {
		if dv, ok := d.(bool); ok {
			debug = dv
		}
	}
	state := &MatchState{
		debug: debug,
	}

	if state.debug {
		logger.Info("match init, starting with debug: %v", state.debug)
	}
	tickRate := 1
	label := "skill=100-150"

	return state, tickRate, label
}

func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	if state.(*MatchState).debug {
		logger.Info("match join attempt username %v user_id %v session_id %v node %v with metadata %v", presence.GetUsername(), presence.GetUserId(), presence.GetSessionId(), presence.GetNodeId(), metadata)
	}

	return state, true, ""
}

func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	if state.(*MatchState).debug {
		for _, presence := range presences {
			logger.Info("match join username %v user_id %v session_id %v node %v", presence.GetUsername(), presence.GetUserId(), presence.GetSessionId(), presence.GetNodeId())
		}
	}

	return state
}

func (m *Match) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	if state.(*MatchState).debug {
		for _, presence := range presences {
			logger.Info("match leave username %v user_id %v session_id %v node %v", presence.GetUsername(), presence.GetUserId(), presence.GetSessionId(), presence.GetNodeId())
		}
	}

	return state
}

func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	if state.(*MatchState).debug {
		logger.Info("match loop match_id %v tick %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), tick)
		logger.Info("match loop match_id %v message count %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), len(messages))
	}

	if tick >= 10 {
		return nil
	}
	return state
}

func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	if state.(*MatchState).debug {
		logger.Info("match terminate match_id %v tick %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), tick)
		logger.Info("match terminate match_id %v grace seconds %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), graceSeconds)
	}

	return state
}

func (m *Match) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	if state.(*MatchState).debug {
		logger.Info("match signal match_id %v tick %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), tick)
		logger.Info("match signal match_id %v data %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), data)
	}
	return state, "signal received: " + data
}

func eventSessionStart(ctx context.Context, logger runtime.Logger, evt *api.Event) {
	logger.Info("session start %v %v", ctx, evt)
}

func eventSessionEnd(ctx context.Context, logger runtime.Logger, evt *api.Event) {
	logger.Info("session end %v %v", ctx, evt)
}
