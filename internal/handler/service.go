package handler

import (
	"context"
	"super-payment-kun/internal/model"
	"time"
)

// TODO: Write test code using mock.
//
//go:generate mockgen -source=service.go -destination=service_mock.go -package=handler
type CreateInvoiceSvc interface {
	CreateInvoice(ctx context.Context,
		user model.User, vendor model.Vendor, vendorBankAcount model.VendorBankAccount,
		paymentAmount float64, dueDate time.Time,
	) (*model.Invoice, *model.Vendor, *model.VendorBankAccount, error)
}

type GetInvoicesSvc interface {
	GetInvoices(ctx context.Context,
		user model.User,
		fromDueDate, toDueDate time.Time,
		limit int,
		lastID string, lastDueDate time.Time,
		direction string,
	) (
		invoinces []model.Invoice,
		newLastID string, newLastDueDate time.Time,
		hasNext, hasPrev bool,
		err error,
	)
}
