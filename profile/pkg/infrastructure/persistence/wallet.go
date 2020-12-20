package persistence

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	profile "server/profile/pkg/profile/gen"
)

func (managerDB *UserDBManager) GetWallets(ctx context.Context, in *profile.UserID) (*profile.Wallets, error) {
	ctx, cancel := context.WithTimeout(ctx, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Wallets{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "GetWallets",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	result, err := tx.Query(
		"SELECT title, value FROM accounts WHERE user_id = $1",
		in.Id,
	)
	if err != nil {
		return &profile.Wallets{}, err
	}
	defer result.Close()

	var wallets profile.Wallets
	for result.Next() {
		var wallet profile.Wallet
		var money decimal.Decimal
		if err := result.Scan(&wallet.Title, &money); err != nil {
			return &profile.Wallets{}, err
		}
		wallet.Value, _ = money.Float64()

		wallets.Wallets = append(wallets.Wallets, &wallet)
	}

	if err := result.Err(); err != nil {
		return &profile.Wallets{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Wallets{}, err
	}

	return &wallets, nil
}

func (managerDB *UserDBManager) CreateInitialWallet(c context.Context, in *profile.UserID) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Empty{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "createInitialWallet",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		in.Id,
		"USD",
		decimal.New(1000, 0),
	)
	if err != nil {
		return &profile.Empty{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Empty{}, err
	}

	return &profile.Empty{}, err
}

func (managerDB *UserDBManager) CreateWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Empty{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "CreateWallet",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		in.Id,
		in.Title,
		decimal.New(0, 0),
	)

	if err != nil {
		return &profile.Empty{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Empty{}, err
	}

	return &profile.Empty{}, nil
}

func (managerDB *UserDBManager) CheckWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Check, error) {
	ctx, cancel := context.WithTimeout(ctx, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Check{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "CheckWallet",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	row := tx.QueryRow("SELECT COUNT(user_id) FROM accounts WHERE EXISTS(select * FROM accounts WHERE user_id = $1 AND title = $2)",
		in.Id,
		in.Title,
	)

	var exists int
	if err = row.Scan(&exists); err != nil {
		return &profile.Check{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Check{}, err
	}

	return &profile.Check{Existence: exists != 0}, nil
}

func (managerDB *UserDBManager) SetWallet(ctx context.Context, in *profile.ToSetWallet) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Empty{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "SetWallet",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		in.Id,
		in.NewWallet.Title,
		in.NewWallet.Value,
	)
	if err != nil {
		return &profile.Empty{}, err
	}

	if err = tx.Commit(); err != nil {
		return &profile.Empty{}, err
	}

	return &profile.Empty{}, nil
}

func (managerDB *UserDBManager) GetWallet(ctx context.Context, in *profile.ConcreteWallet) (*profile.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Wallet{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "GetWallet",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	row := tx.QueryRow(
		"SELECT value FROM accounts WHERE user_id = $1 AND title = $2",
		in.Id,
		in.Title,
	)

	var money decimal.Decimal
	if err = row.Scan(&money); err != nil {
		return &profile.Wallet{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Wallet{}, err
	}

	wallet := profile.Wallet{Title: in.Title}
	wallet.Value, _ = money.Float64()

	return &wallet, nil
}

func (managerDB *UserDBManager) UpdateWallet(ctx context.Context, in *profile.ToSetWallet) (*profile.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, managerDB.timing)
	defer cancel()

	tx, err := managerDB.DB.BeginTx(ctx, nil)
	if err != nil {
		return &profile.Empty{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.WithFields(logrus.Fields{
				"infrastructure": "profile",
				"function":       "UpdateWallet",
				"action":         "Rollback",
			}).Debug(err)
		}
	}()

	_, err = tx.Exec(
		"UPDATE accounts SET value = value + $1 WHERE user_id = $2 AND title = $3",
		decimal.NewFromFloat(in.NewWallet.Value),
		in.Id,
		in.NewWallet.Title,
	)

	if err != nil {
		return &profile.Empty{}, err
	}
	if err = tx.Commit(); err != nil {
		return &profile.Empty{}, err
	}

	return &profile.Empty{}, nil
}
