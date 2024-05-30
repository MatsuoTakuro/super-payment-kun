package service

import (
	"context"
	"super-payment-kun/internal/model"
	"super-payment-kun/internal/pkg"
	"time"
)

type createInvoice struct {
	repo CreateInvoiceRepo
}

func NewCreateInvoice(repo CreateInvoiceRepo) *createInvoice {
	return &createInvoice{
		repo: repo,
	}
}

func (s *createInvoice) CreateInvoice(ctx context.Context,
	user model.User, vendor model.Vendor, vendorBankAccount model.VendorBankAccount,
	paymentAmount float64, dueDate time.Time,
) (*model.Invoice, *model.Vendor, *model.VendorBankAccount, error) {

	invoice, err := model.NewInvoice(user, vendor, vendorBankAccount, paymentAmount, dueDate)
	if err != nil {
		return nil, nil, nil, err
	}

	gotVendor, gotVBankAccount, err := s.repo.GetVendorAndBankAccount(ctx, invoice.VendorBankAccount)
	if err != nil {
		return nil, nil, nil, err
	}
	if gotVendor == nil || gotVBankAccount == nil {
		return nil, nil, nil, pkg.NewAPIError(pkg.InvalidParams, nil, "vendor bank account does not exist")
	}

	created, err := s.repo.CreateInvoice(ctx, invoice)
	if err != nil {
		return nil, nil, nil, err
	}

	return created, gotVendor, gotVBankAccount, nil
}
