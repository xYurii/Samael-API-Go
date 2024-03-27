package adapter

import (
	"apisamael/entities"
	"apisamael/utils"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type UserAdapter struct {
	Db *bun.DB
}

func (a *UserAdapter) GetDailyReward(ctx context.Context, user utils.DiscordUser, userData entities.User, ipInfo utils.IPDetails, quantity uint64, first bool) {
	now := time.Now()
	nowMs := now.UnixNano() / int64(time.Millisecond)
	createOrUpdate := ""

	if first {
		createOrUpdate = fmt.Sprintf(`INSERT INTO "Dailys" ("id", "user_id", "cooldown", "ip", "email", "country", "region", "city", "org", "createdAt", "updatedAt") VALUES (DEFAULT, '%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', NOW(), NOW())`,
			user.ID, nowMs, ipInfo.Query, user.Email, ipInfo.Country, ipInfo.Region, ipInfo.City, ipInfo.Org)
	} else {
		createOrUpdate = fmt.Sprintf(`UPDATE "Dailys" SET "cooldown" = %d, "email" = '%s', "updatedAt" = NOW() WHERE ip = '%s'`,
			nowMs, user.Email, ipInfo.Query)
	}

	a.Db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.Query(createOrUpdate)
		if err != nil {
			return err
		}
		_, err = tx.NewInsert().Model(&entities.Transaction{
			Source:        3,
			GivenBy:       user.ID,
			ReceivedBy:    user.ID,
			GivenByTag:    user.Username,
			ReceivedByTag: user.Username,
			GivenAt:       nowMs,
			Quantity:      quantity,
			OtherUsers:    []string{user.ID},
			CreatedAt:     now,
			UpdatedAt:     now,
		}).Exec(ctx)

		if err != nil {
			return err
		}

		_, err = tx.NewInsert().
			Model(&entities.User{
				Money:     userData.Money + quantity,
				Daily:     uint64(nowMs),
				Tag:       fmt.Sprintf("%s#%s", user.Username, user.Discriminator),
				ID:        user.ID,
				CreatedAt: now,
				UpdatedAt: now,
			}).
			On("CONFLICT (id) DO UPDATE").
			Set("daily = EXCLUDED.daily").
			Set("money = EXCLUDED.money").
			Exec(ctx)

		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model(&entities.User{ID: user.ID}).
			Set("user_tasks = ?", userData.UserTasks).
			Exec(ctx)

		if err != nil {
			return err
		}
		return nil
	})

}

func (a *UserAdapter) FetchUserByIp(ctx context.Context, u entities.User, ip string) (daily entities.Daily, err error) {
	var dailys []entities.Daily

	err = a.Db.NewSelect().
		Model(&dailys).
		ColumnExpr("d.user_id, d.ip, d.cooldown, d.email").
		Join("JOIN \"Users\" ON \"Users\".id = d.user_id").
		Where("d.ip = ?", ip).
		Order("d.cooldown DESC").
		Scan(ctx)

	if err != nil {
		return entities.Daily{}, err
	}
	return dailys[0], nil
}

func (a *UserAdapter) IsBlacklisted(ctx context.Context, u entities.User) (bool, error) {
	b := entities.Blacklist{
		ID: u.ID,
	}

	count, err := a.Db.NewSelect().Model(&b).Where("id = ?", b.ID).Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (a *UserAdapter) GetUser(ctx context.Context, u entities.User) (user entities.User) {
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

func (a *UserAdapter) InsertUser(ctx context.Context, user entities.User) error {
	return a.Db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock('%s')", user.ID))
		_, err := tx.NewInsert().Model(&user).Exec(ctx)
		return err
	})
}
