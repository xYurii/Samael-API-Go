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

type User struct {
	bun.BaseModel `bun:"table:\"Users\",alias:u"`

	ID        string    `bun:"id,pk"`
	UserTasks UserTasks `bun:"user_tasks,type:json"`
	Money     uint64    `bun:"money"`
	Daily     uint64    `bun:"daily"`
	Tag       string    `bun:"tag"`
	CreatedAt time.Time `bun:"\"createdAt\",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"\"updatedAt\",nullzero,notnull,default:current_timestamp"`
}
