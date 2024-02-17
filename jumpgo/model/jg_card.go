package model

import (
	"database/sql/driver"
	"encoding/json"
)

type JgCard struct {
	Id         int        `gorm:"column:id;NOT NULL" json:"id"`
	Name       string     `gorm:"column:name;default:::;NOT NULL;comment:'名称'" json:"name"`
	NameCn     string     `gorm:"column:name_cn;comment:'名称'" json:"name_cn"`
	Introduce  string     `gorm:"column:introduce;comment:'介绍'" json:"introduce"`
	IsNextTurn int        `gorm:"column:is_next_turn;default:1;NOT NULL;comment:'回合变更'" json:"is_next_turn"`
	CardType   int        `gorm:"column:card_type;default:1;NOT NULL;comment:'牌类型'" json:"card_type"`
	ZombieInfo ZombieInfo `gorm:"column:zombie_info;comment:'僵尸详情'" json:"zombie_info"`
	GunInfo    ZombieInfo `gorm:"column:gun_info;comment:'枪详情'" json:"gun_info"`
	Match4     int        `gorm:"column:match_4;default:4;NOT NULL;comment:'4人局'" json:"match_4"`
}

type ZombieInfo struct {
	Hp     int `json:"hp"`
	Damage int `json:"damage"`
}

func (zi *ZombieInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return nil
	}
	result := ZombieInfo{}
	err := json.Unmarshal([]byte(str), &result)
	*zi = result
	return err
}

// Value implements the driver Valuer interface.
func (zi ZombieInfo) Value() (driver.Value, error) {
	if zi.Damage == 0 && zi.Hp == 0 {
		return nil, nil
	}
	return json.Marshal(zi)
}

func (j *JgCard) TableName() string {
	return "jg_card"
}
