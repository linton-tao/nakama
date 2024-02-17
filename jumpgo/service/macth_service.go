package service

import (
	"context"
	"math/rand"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/heroiclabs/nakama/v3/jumpgo/api"
	"github.com/heroiclabs/nakama/v3/jumpgo/model"
	"gorm.io/gorm"
)

func InitCard(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, matchType int, tx *gorm.DB) (cardList api.CardList, err error) {
	var list []*model.JgCard
	err = tx.Find(&list).Error
	for _, card := range list {
		cardNum := 0
		switch matchType {
		case 4:
			cardNum = card.Match4
			break
		}
		for i := 0; i < cardNum; i++ {
			cardList.List = append(cardList.List, card)
			cardList.End++
		}
	}
	randList := cardList.List[cardList.Start:cardList.End]
	rand.Shuffle(len(randList), func(i, j int) {
		randList[i], randList[j] = randList[j], randList[i]
	})
	return
}
