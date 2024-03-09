package adapter

import (
	"apisamael/entities"
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
)

type UserAdapter struct {
	Db *bun.DB
}

func (a UserAdapter) GetUser(ctx context.Context, u entities.User) (user entities.User) {
	a.Db.NewSelect().Model(&user).Where("id = ?", u.ID).Scan(ctx)
	if user.ID == "" {
		now := time.Now()

		user.ID = u.ID
		user.UpdatedAt = now
		user.CreatedAt = now
		user.UserTasks = entities.UserTasks{
			Reps:           0,
			Work:           0,
			Crime:          0,
			Bets:           0,
			Daily:          false,
			RaffleQuantity: 0,
			Completed:      false,
		}
		a.InsertUser(ctx, user)
		user = a.GetUser(ctx, u)
	}
	return
}

func (a UserAdapter) InsertUser(ctx context.Context, user entities.User) error {
	return a.Db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock('%s')", user.ID))
		_, err := tx.NewInsert().Model(&user).Exec(ctx)
		return err
	})
}
