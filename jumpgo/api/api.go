package api

import "github.com/heroiclabs/nakama/v3/jumpgo/model"

type RpcInitMacthRequest struct {
	MatchType int `json:"match_type"`
}

type CardList struct {
	Start int             `json:"start"`
	End   int             `json:"end"`
	List  []*model.JgCard `json:"list"`
}
