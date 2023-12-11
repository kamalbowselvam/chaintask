package service

import (
	"context"

	"github.com/kamalbowselvam/chaintask/db"
	"github.com/kamalbowselvam/chaintask/domain"
)

func (srv *service) PayForATask(ctx context.Context, arg db.CreateTaskPaymentParams) (domain.TaskPayment, error){
	// TODO Verification we don't overflow
	payment, err := srv.globalRepository.PayForATask(ctx, arg)
	if err != nil {
		return domain.TaskPayment{}, err
	}
	// FIXME Should trigger a budget recomputation in a async worker see CHAIN-13
	return payment, nil
}