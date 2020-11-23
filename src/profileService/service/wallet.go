package main

import (
	"context"
	"github.com/shopspring/decimal"
	profile "server/src/profileService/profile/gen"
	"time"
)

func (managerDB *UserDBManager) GetWallets(ctx context.Context, in *profile.UserID) (*profile.Wallets, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query(
		"SELECT title, value FROM accounts WHERE user_id = $1",
		in.Id,
	)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var wallets profile.Wallets
	for result.Next() {
		var wallet profile.Wallet
		var money decimal.Decimal
		if err := result.Scan(&wallet.Title, &money); err != nil {
			return nil, err
		}
		wallet.Value, _ = money.Float64()

		wallets.AllWallets = append(wallets.AllWallets, &wallet)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if len(wallets.AllWallets) == 0 {
		if err = managerDB.createInitialWallet(in.Id); err != nil {
			return nil, err
		}
		money, _ := decimal.New(100000, 0).Float64()
		wallets.AllWallets = append(
			wallets.AllWallets,
			&profile.Wallet{Title: "RUB", Value: money,
			})
	}

	return &wallets, nil
}

func (managerDB *UserDBManager) createInitialWallet(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		id,
		"RUB",
		decimal.New(100000, 0),
	)

	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (managerDB *UserDBManager) CreateWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		in.Id,
		in.Title,
		decimal.New(0, 0),
	)

	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (managerDB *UserDBManager) CheckWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Check, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT COUNT(user_id) FROM accounts WHERE EXISTS(select * FROM accounts WHERE user_id = $1 AND title = $2)",
		in.Id,
		in.Title,
	)

	var exists int
	if err = row.Scan(&exists); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &profile.Check{Existence: exists != 0}, nil
}

func (managerDB *UserDBManager) SetWallet(ctx context.Context, in *profile.ToSetWallet) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		in.Id,
		in.NewWallet.Title,
		in.NewWallet.Value,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (managerDB *UserDBManager) GetWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		"SELECT value FROM accounts WHERE user_id = $1 AND title = $2",
		in.Id,
		in.Title,
	)

	var money decimal.Decimal
	if err = row.Scan(&money); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	wallet := profile.Wallet{Title: in.Title}
	wallet.Value, _ = money.Float64()

	return &wallet, nil
}

func (managerDB *UserDBManager) UpdateWallet(ctx context.Context, in *profile.ToSetWallet) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE accounts SET value = value + $1 WHERE user_id = $2 AND title = $3",
		decimal.NewFromFloat(in.NewWallet.Value),
		in.Id,
		in.NewWallet.Title,
	)

	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
