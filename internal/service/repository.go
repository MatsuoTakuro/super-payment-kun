package service

import (
	"context"
	"super-payment-kun/internal/model"
	"time"
)

// TODO: Write test code using mock.
//
//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=service
type CreateInvoiceRepo interface {
	GetVendorAndBankAccount(ctx context.Context, vendorBankAccount model.VendorBankAccount) (*model.Vendor, *model.VendorBankAccount, error)
	CreateInvoice(ctx context.Context, invoice *model.Invoice) (*model.Invoice, error)
}

type GetInvoicesRepo interface {
	// GetInvoicesBetweenDueDates retrieves all invoices with due dates between the specified range for the specified company
	GetInvoicesBetweenDueDates(ctx context.Context,
		company model.Company,
		fromDueDate, toDueDate time.Time,
		limit int,
		lastID string, lastDueDate time.Time,
		direction string,
	) (
		invoices []model.Invoice,
		newLastID string, newLastDueDate time.Time,
		hasNext, hasPrev bool,
		err error,
	)
}
