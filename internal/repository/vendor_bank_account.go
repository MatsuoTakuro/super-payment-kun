package repository

import (
	"context"
	"super-payment-kun/internal/model"
	"super-payment-kun/internal/pkg"
)

func (r *repository) GetVendorAndBankAccount(
	ctx context.Context,
	vendorBankAccount model.VendorBankAccount,
) (*model.Vendor, *model.VendorBankAccount, error) {

	sql := `
		SELECT
			v.id, v.name, v.company_id,
			vba.id, vba.vendor_id, vba.bank_name, vba.branch_name, vba.account_number, vba.account_name
		FROM
			vendor_bank_accounts vba
		INNER JOIN
			vendors v ON v.id = vba.vendor_id
		WHERE
			vba.id = ? AND vba.vendor_id = ? AND v.company_id = ? AND v.deleted_at IS NULL AND vba.deleted_at IS NULL
		LIMIT 1
	`

	var gotVendor model.Vendor
	var gotBankAccount model.VendorBankAccount

	err := r.db.QueryRowContext(ctx, sql, vendorBankAccount.ID, vendorBankAccount.Vendor.ID, vendorBankAccount.Vendor.Company.ID).Scan(
		&gotVendor.ID, &gotVendor.Name, &gotVendor.Company.ID,
		&gotBankAccount.ID, &gotBankAccount.Vendor.ID, &gotBankAccount.BankName, &gotBankAccount.BranchName, &gotBankAccount.AccountNumber, &gotBankAccount.AccountName,
	)
	if err != nil {
		if err == noRowErr {
			return nil, nil, nil
		}
		return nil, nil, pkg.NewAPIError(pkg.SelectDataFailed, err, "failed to retrieve vendor and bank account")
	}

	return &gotVendor, &gotBankAccount, nil
}
