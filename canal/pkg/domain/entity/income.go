package entity

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/shopspring/decimal"
	profile "server/profile/pkg/profile/gen"
)

//easyjson:json
type (
	Income struct {
		Id     int64 `json:"id"`
		Period string
	}

	WalletState struct {
		Value     decimal.Decimal      `json:"value"`
		UpdatedAt *timestamp.Timestamp `json:"updated_at"`
	}

	WalletStates []WalletState
)

func (in *Income) ConvertToGRPC() *profile.IncomeParameters {
	return &profile.IncomeParameters{Id: in.Id, Period: in.Period}
}

func ConvertToWalletStates(profileStates *profile.WalletStates) WalletStates {
	entityStates := make(WalletStates, 0, len(profileStates.States))
	for _, pflState := range profileStates.States {
		entityStates = append(entityStates, WalletState{
			Value:     decimal.NewFromFloat(pflState.Value),
			UpdatedAt: pflState.UpdatedAt,
		})
	}

	return entityStates
}
