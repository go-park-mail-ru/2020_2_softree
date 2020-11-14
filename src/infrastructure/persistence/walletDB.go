package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"server/src/domain/entity"
	"server/src/infrastructure/config"
	"time"
)

type WalletDBManager struct {
	DB *sql.DB
}

func NewWalletDBManager() (*WalletDBManager, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.UserDatabaseConfig.Host,
		5432,
		config.UserDatabaseConfig.User,
		config.UserDatabaseConfig.Password,
		config.UserDatabaseConfig.Schema)

	db, err := sql.Open("postgres", psqlInfo)

	db.SetMaxOpenConns(10)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		return nil, err
	}

	return &WalletDBManager{DB: db}, nil
}

func (wm *WalletDBManager) GetWallets(id uint64) ([]entity.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := wm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query(
		"SELECT title, value FROM accounts WHERE user_id = $1",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	wallets := make([]entity.Wallet, 0)
	for result.Next() {
		var wallet entity.Wallet
		if err := result.Scan(&wallet.Title, &wallet.Value); err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if len(wallets) == 0 {
		if err = wm.createInitialWallet(id); err != nil {
			return nil, nil
		}
		wallets = append(wallets, entity.Wallet{Title: "USD", Value: decimal.New(1000, 0)})
	}

	return wallets, nil
}

func (wm *WalletDBManager) createInitialWallet(id uint64) error  {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := wm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		id,
		"USD",
		decimal.New(1000, 0),
	)

	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (wm *WalletDBManager) CreateWallet(id uint64, title string) error  {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := wm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		id,
		title,
		decimal.New(0, 0),
	)

	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (wm *WalletDBManager) CheckWallet(id uint64, title string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := wm.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT COUNT(user_id) FROM accounts WHERE EXISTS(select * FROM accounts WHERE user_id = $1 AND title = $2)",
		id,
		title,
		)

	var exists int
	if err = row.Scan(&exists); err != nil {
		return false, err
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}

	return exists != 0, nil
}

func (wm *WalletDBManager) SetWallet(id uint64, wallet entity.Wallet) error  {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := wm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"INSERT INTO accounts (user_id, title, value) VALUES ($1, $2, $3)",
		id,
		wallet.Title,
		wallet.Value,
	)

	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (wm *WalletDBManager) GetWallet(id uint64, title string) (entity.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := wm.DB.BeginTx(ctx, nil)
	if err != nil {
		return entity.Wallet{}, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		"SELECT value FROM accounts WHERE user_id = $1 AND title = $2",
		id,
		title,
	)

	var wallet = entity.Wallet{Title: title}
	if err = row.Scan(&wallet.Value);err != nil {
		return entity.Wallet{}, err
	}

	if err = tx.Commit(); err != nil {
		return entity.Wallet{}, err
	}

	return wallet, nil
}

func (wm *WalletDBManager) UpdateWallet(id uint64, wallet entity.Wallet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := wm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		"UPDATE accounts SET value = value + $1 WHERE user_id = $2 AND title = $3",
		wallet.Value,
		id,
		wallet.Title,
	)

	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
