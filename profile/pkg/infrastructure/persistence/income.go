package persistence

import (
	"context"
	"errors"
	profile "server/profile/pkg/profile/gen"
)

func (managerDB *UserDBManager) GetIncome(c context.Context, in *profile.IncomeParameters) (*profile.Income, error) {
	return &profile.Income{}, errors.New("implement me")
}
