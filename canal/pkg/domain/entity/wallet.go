package entity

import (
	json "github.com/mailru/easyjson"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"net/http"
	profile "server/profile/pkg/profile/gen"
)

type Wallet struct {
	Title  string          `json:"title"`
	Value  decimal.Decimal `json:"value"`
	UserId int64
}

type Wallets struct {
	Wallets []Wallet `json:"wallets"`
}

func ConvertToWallets(profileWallets *profile.Wallets) Wallets {
	var entityWallets Wallets
	entityWallets.Wallets = make([]Wallet, 0, len(profileWallets.Wallets))
	for _, wallet := range profileWallets.Wallets {
		entityWallets.Wallets = append(entityWallets.Wallets, Wallet{Title: wallet.Title, Value: decimal.NewFromFloat(wallet.Value)})
	}

	return entityWallets
}

func GetWalletFromBody(body io.ReadCloser) (Wallet, Description, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return Wallet{}, Description{Action: "ReadAll", Status: http.StatusInternalServerError}, err
	}
	defer body.Close()

	var wallet Wallet
	err = json.Unmarshal(data, &wallet)
	if err != nil {
		return Wallet{}, Description{Action: "Unmarshal", Status: http.StatusInternalServerError}, err
	}
	return wallet, Description{}, nil
}
