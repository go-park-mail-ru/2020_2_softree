package persistence

import (
	"context"
	profile "server/profile/pkg/profile/gen"
)

func (managerDB *UserDBManager) GetIncome(c context.Context, in *profile.IncomeParameters) (*profile.Income, error) {
	return &profile.Income{}, nil
}
