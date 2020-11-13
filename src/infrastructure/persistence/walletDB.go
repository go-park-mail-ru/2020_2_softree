package persistence

import (
	"context"
	"database/sql"
	"fmt"
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

func (wm *WalletDBManager) GetWallet(id uint64) ([]entity.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := wm.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Query(
		"SELECT title, value FROM accounts WHERE id = $1",
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

	return wallets, nil
}

func (wm *WalletDBManager) SetWallet(id uint64, newWallet entity.Wallet) error  {
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
		newWallet.Title,
		newWallet.Value,
	)

	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

