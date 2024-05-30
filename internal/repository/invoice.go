package repository

import (
	"context"
	"super-payment-kun/internal/model"
	"super-payment-kun/internal/pkg"
)

func (r *repository) CreateInvoice(
	ctx context.Context, invoice *model.Invoice,
) (*model.Invoice, error) {

	sql := `INSERT INTO invoices
		(
			id, company_id, user_id, vendor_id, vendor_bank_account_id,
			payment_amount, status,
			fee, fee_rate,
			consumption_tax, consumption_tax_rate,
			billed_amount,
			issue_date, due_date,
			created_by, created_at,
			updated_by,	updated_at
		) VALUES (
			?, ?, ?, ?, ?,
			?, ?,
			?, ?,
			?, ?,
			?,
			?, ?,
			?, ?,
			?, ?
		)
	`

	_, err := r.db.ExecContext(ctx, sql,
		invoice.ID, invoice.Company.ID, invoice.User.ID, invoice.Vendor.ID, invoice.VendorBankAccount.ID,
		invoice.PaymentAmount, invoice.Status,
		invoice.Fee, invoice.FeeRate,
		invoice.ConsumptionTax, invoice.ConsumptionTaxRate,
		invoice.BilledAmount,
		invoice.IssueDate, invoice.DueDate,
		invoice.CreatedBy, invoice.CreatedAt,
		invoice.UpdatedBy, invoice.UpdatedAt,
	)
	if err != nil {
		if isDuplicateEntryErr(err) {
			return nil, pkg.NewAPIError(pkg.RegisterDuplicateDataRestricted, err,
				"creating an invoice with the same id is restricted. please try it again")
		}
		if isForeignKeyConstraintErr(err) {
			return nil, pkg.NewAPIError(pkg.InvalidParams, err,
				"failed to create invoice due to other invalid ids")
		}

		return nil, pkg.NewAPIError(pkg.RegisterDataFailed, err, "failed to create an invoice")
	}

	return invoice, nil
}
