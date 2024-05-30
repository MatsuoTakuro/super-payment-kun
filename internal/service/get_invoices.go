package service

import (
	"context"
	"errors"
	"fmt"
	"super-payment-kun/internal/model"
	"super-payment-kun/internal/pkg"
	"time"
)

type getInvoices struct {
	repo GetInvoicesRepo
}

func NewGetInvoices(repo GetInvoicesRepo) *getInvoices {
	return &getInvoices{
		repo: repo,
	}
}

func (s *getInvoices) GetInvoices(ctx context.Context,
	user model.User,
	fromDueDate, toDueDate time.Time,
	limit int,
	lastID string, lastDueDate time.Time,
	direction string,
) (
	invoices []model.Invoice,
	newLastID string, newLastDueDate time.Time,
	hasNext, hasPrev bool,
	err error,
) {

	if !user.CanReadAllInvoicesForCompany() {
		err := errors.New("user does not have permission to read all invoices")
		return nil, "", time.Time{}, false, false, pkg.NewAPIError(pkg.PermissionDenied, err, err.Error())
	}

	from := pkg.TruncateToDate(fromDueDate)
	to := pkg.TruncateToDate(toDueDate)

	if from.After(to) {
		err := fmt.Errorf("start due date (%s) is after end due date (%s)", pkg.FormatToDate(from), pkg.FormatToDate(to))
		return nil, "", time.Time{}, false, false, pkg.NewAPIError(pkg.InvalidParams, err, err.Error())
	}

	maxDeferDate := pkg.TruncateToDate(time.Now().AddDate(0, 0, model.MaxDeferDays))
	if to.After(maxDeferDate) {
		err := fmt.Errorf("end due date (%s) is after max defer date (%s)", pkg.FormatToDate(to), pkg.FormatToDate(maxDeferDate))
		return nil, "", time.Time{}, false, false, pkg.NewAPIError(pkg.InvalidParams, err, err.Error())
	}

	company := model.Company{
		ID: user.Company.ID,
	}

	invoices, newLastID, newLastDueDate, hasNext, hasPrev, err = s.repo.GetInvoicesBetweenDueDates(
		ctx, company,
		from, to,
		limit,
		lastID, lastDueDate,
		direction,
	)
	if err != nil {
		return nil, "", time.Time{}, false, false, err
	}

	return invoices, newLastID, newLastDueDate, hasNext, hasPrev, nil
}
