package model

import (
	"errors"
	"fmt"
	"math"
	"time"

	"super-payment-kun/internal/pkg"
)

type InvoiceStatus int8

const (
	Unprocessed InvoiceStatus = 0
	Prcoessing  InvoiceStatus = 1
	Paid        InvoiceStatus = 2
	Error       InvoiceStatus = 3
)

func (s InvoiceStatus) ToJPName() string {
	switch s {
	case Unprocessed:
		return "未処理"
	case Prcoessing:
		return "処理中"
	case Paid:
		return "支払済"
	case Error:
		return "エラー"
	default:
		return "不明" // TODO: handle unknown status
	}
}

const (
	feeRate            = 0.04 // TODO: fetch the latest rate from config table
	consumptionTaxRate = 0.10 // TODO: fetch the latest rate from config table
	minPaymentAmount   = 250  // TODO: fetch the latest rate from config table
)

var roundYenAmount = math.Floor // round down

const (
	MaxDeferDays = 60
)

type Invoice struct {
	ID                 string            `db:"id"` // uuid
	Company            Company           `db:"company"`
	User               User              `db:"user"`
	Vendor             Vendor            `db:"vendor"`
	VendorBankAccount  VendorBankAccount `db:"vendor_bank_account"`
	PaymentAmount      float64           `db:"payment_amount"` // PaymentAmount is amount without fee and tax. Unit is yen.
	Status             InvoiceStatus     `db:"status"`
	Fee                float64           `db:"fee"`
	FeeRate            float64           `db:"fee_rate"`
	ConsumptionTax     float64           `db:"consumption_tax"` // Unit is yen.
	ConsumptionTaxRate float64           `db:"consumption_tax_rate"`
	BilledAmount       float64           `db:"billed_amount"` // Unit is yen.
	IssueDate          time.Time         `db:"issue_date"`    // no timestamp like 2021-01-01 00:00:00
	DueDate            time.Time         `db:"due_date"`      // no timestamp like 2021-01-01 00:00:00
	CreatedBy          string            `db:"created_by"`    // user id (uuid v4) or system id
	CreatedAt          time.Time         `db:"created_at"`
	UpdatedBy          string            `db:"updated_by"` // user id (uuid v4) or system id
	UpdatedAt          time.Time         `db:"updated_at"`
	DeletedAt          *time.Time        `db:"deleted_at"`
}

func NewInvoice(
	user User, vendor Vendor, vendorBankAcount VendorBankAccount,
	paymentAmount float64,
	dueDatetime time.Time,
) (*Invoice, error) {

	if !user.CanCreateInvoice() {
		err := errors.New("user does not have permission to create invoice")
		return nil, pkg.NewAPIError(pkg.PermissionDenied, err, err.Error())
	}

	if paymentAmount < minPaymentAmount {
		err := fmt.Errorf("payment amount is too small: %.1f yen (min: %d yen)", paymentAmount, minPaymentAmount)
		return nil, pkg.NewAPIError(pkg.PaymentAmountTooSmall, err, err.Error())
	}
	fee := roundYenAmount(paymentAmount * feeRate)
	consumptionTax := roundYenAmount(fee * consumptionTaxRate)
	billedAmount := paymentAmount + fee + consumptionTax

	dueDate := pkg.TruncateToDate(dueDatetime)
	nowTime := time.Now()
	today := pkg.TruncateToDate(nowTime)

	if dueDate.Before(today) {
		err := fmt.Errorf("due date is past: %s (today: %s)", pkg.FormatToDate(dueDate), pkg.FormatToDate(today))
		return nil, pkg.NewAPIError(pkg.DueDateIsPassed, err, err.Error())
	}
	maxDeferDate := pkg.TruncateToDate(today.AddDate(0, 0, MaxDeferDays))
	if dueDate.After(maxDeferDate) {
		err := fmt.Errorf("due date is too far: %s (today: %s): max defer days are %d", pkg.FormatToDate(dueDate), pkg.FormatToDate(today), MaxDeferDays)
		return nil, pkg.NewAPIError(pkg.DueDateExceedMaxDeferDate, err, err.Error())
	}

	id := pkg.NewUUID()
	return &Invoice{
		ID:                 id,
		Company:            Company{ID: user.Company.ID},
		User:               user,
		Vendor:             vendor,
		VendorBankAccount:  vendorBankAcount,
		PaymentAmount:      paymentAmount,
		Status:             Unprocessed, // default status
		Fee:                fee,
		FeeRate:            feeRate,
		ConsumptionTax:     consumptionTax,
		ConsumptionTaxRate: consumptionTaxRate,
		BilledAmount:       billedAmount,
		IssueDate:          today,
		DueDate:            dueDate,
		CreatedBy:          user.ID,
		CreatedAt:          nowTime,
		UpdatedBy:          user.ID,
		UpdatedAt:          nowTime,
		DeletedAt:          nil,
	}, nil
}
