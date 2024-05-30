package repository

import (
	"context"
	"errors"
	"fmt"
	"super-payment-kun/internal/model"
	"super-payment-kun/internal/pkg"
	"time"
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

func (r *repository) GetInvoicesBetweenDueDates(ctx context.Context,
	company model.Company,
	fromDueDate, toDueDate time.Time,
	limit int,
	lastID string,
	lastDueDate time.Time,
	direction string,
) (invoices []model.Invoice, newLastID string, newLastDueDate time.Time, hasNext, hasPrev bool, err error) {

	order := "ASC" // Default order (forward = oldest due date first)
	compare := ">" // after last due date
	if direction == "bwd" {
		order = "DESC" // Reverse order (backward = newest due date first)
		compare = "<"  // before last due date
	}

	query := fmt.Sprintf(`
		SELECT
			i.id,
			i.company_id,
			i.user_id,
			i.vendor_id,
			v.name AS vendor_name,
			i.vendor_bank_account_id,
			vba.bank_name,
			vba.branch_name,
			vba.account_number,
			vba.account_name,
			i.payment_amount,
			i.status,
			i.fee,
			i.fee_rate,
			i.consumption_tax,
			i.consumption_tax_rate,
			i.billed_amount,
			i.issue_date,
			i.due_date
		FROM
			invoices i
		JOIN
			vendors v ON i.vendor_id = v.id
		JOIN
			vendor_bank_accounts vba ON i.vendor_bank_account_id = vba.id
		WHERE
			i.company_id = ? AND
			i.due_date BETWEEN ? AND ? AND i.deleted_at IS NULL
			AND (? = '' OR i.due_date %s ? OR (i.due_date = ? AND i.id %s ?))
		ORDER BY
			i.due_date %s, i.id %s
		LIMIT ?
	`, compare, compare, order, order)

	rows, err := r.db.QueryxContext(ctx, query,
		company.ID,
		fromDueDate, toDueDate,
		lastID, lastDueDate, lastDueDate, lastID,
		limit+1) // Fetch one more to check if there are more invoices
	if err != nil {
		return nil, "", time.Time{}, false, false, err
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	invoices = make([]model.Invoice, 0, limit)

	for rows.Next() {
		var invoice model.Invoice
		var vendor model.Vendor
		var vbankAccount model.VendorBankAccount

		err = rows.Scan(
			&invoice.ID,
			&invoice.Company.ID,
			&invoice.User.ID,
			&vendor.ID,
			&vendor.Name,
			&vbankAccount.ID,
			&vbankAccount.BankName,
			&vbankAccount.BranchName,
			&vbankAccount.AccountNumber,
			&vbankAccount.AccountName,
			&invoice.PaymentAmount,
			&invoice.Status,
			&invoice.Fee,
			&invoice.FeeRate,
			&invoice.ConsumptionTax,
			&invoice.ConsumptionTaxRate,
			&invoice.BilledAmount,
			&invoice.IssueDate,
			&invoice.DueDate,
		)
		if err != nil {
			return nil, "", time.Time{}, false, false, err
		}

		invoice.Vendor = vendor
		invoice.VendorBankAccount = vbankAccount
		invoices = append(invoices, invoice)
	}

	if err = rows.Err(); err != nil {
		return nil, "", time.Time{}, false, false, err
	}

	if len(invoices) > limit {
		hasNext = true
		invoices = invoices[:limit]
	}

	if len(invoices) > 0 {
		lastInvoice := invoices[len(invoices)-1]
		newLastID = lastInvoice.ID
		newLastDueDate = lastInvoice.DueDate
	}

	hasPrev = lastID != "" // If lastID is not empty, it means there is a previous invoice

	return invoices, newLastID, newLastDueDate, hasNext, hasPrev, nil
}
