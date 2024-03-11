package entities

import (
	"time"

	"github.com/uptrace/bun"
)

type UserTasks struct {
	Reps           int64 `json:"reps"`
	Work           int64 `json:"work"`
	Crime          int64 `json:"crime"`
	Bets           int64 `json:"bets"`
	Daily          bool  `json:"daily"`
	RaffleQuantity int64 `json:"raffle_quantity"`
	Completed      bool  `json:"completed"`
}

type Transaction struct {
	bun.BaseModel `bun:"table:\"Transactions\",alias:t"`

	ID            int       `bun:"id"`
	Source        int       `bun:"source"`
	GivenBy       string    `bun:"given_by"`
	ReceivedBy    string    `bun:"received_by"`
	GivenByTag    string    `bun:"given_by_tag"`
	ReceivedByTag string    `bun:"received_by_tag"`
	GivenAt       int64     `bun:"given_at"`
	Quantity      uint64    `bun:"quantity"`
	OtherUsers    []string  `bun:"other_users"`
	CreatedAt     time.Time `bun:"\"createdAt\",nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"\"updatedAt\",nullzero,notnull,default:current_timestamp"`
}

type Daily struct {
	bun.BaseModel `bun:"table:\"Dailys\",alias:d"`

	ID        string    `bun:"id"`
	UserID    string    `bun:"user_id"`
	Cooldown  int64     `bun:"cooldown"`
	IP        string    `bun:"ip"`
	Email     string    `bun:"email"`
	Reward    uint64    `bun:"reward"`
	Country   string    `bun:"country"`
	Region    string    `bun:"region"`
	City      string    `bun:"city"`
	Org       string    `bun:"org"`
	CreatedAt time.Time `bun:"\"createdAt\",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"\"updatedAt\",nullzero,notnull,default:current_timestamp"`
}

type User struct {
	bun.BaseModel `bun:"table:\"Users\",alias:u"`

	ID               string    `bun:"id,pk"`
	UserTasks        UserTasks `bun:"user_tasks,type:json"`
	Money            uint64    `bun:"money"`
	Daily            uint64    `bun:"daily"`
	Tag              string    `bun:"tag"`
	IsPremium        bool      `bun:"is_premium"`
	IsBoosterPremium bool      `bun:"is_booster_premium"`
	CreatedAt        time.Time `bun:"\"createdAt\",nullzero,notnull,default:current_timestamp"`
	UpdatedAt        time.Time `bun:"\"updatedAt\",nullzero,notnull,default:current_timestamp"`
}

type Blacklist struct {
	bun.BaseModel `bun:"table:\"Blacklists\",alias:b"`

	ID        string    `bun:"id"`
	BannedAt  uint64    `bun:"banned_at"`
	Duration  uint64    `bun:"duration"`
	Reason    string    `bun:"reason"`
	BannedBy  string    `bun:"banned_by"`
	Permanent bool      `bun:"permanent"`
	CreatedAt time.Time `bun:"\"createdAt\",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"\"updatedAt\",nullzero,notnull,default:current_timestamp"`
}
