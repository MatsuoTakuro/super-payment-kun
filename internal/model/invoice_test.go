package model

import (
	"super-payment-kun/internal/pkg"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewInvoice(t *testing.T) {
	type args struct {
		user              User
		vendor            Vendor
		vendorBankAccount VendorBankAccount
		paymentAmount     float64
		dueDate           time.Time
	}
	company := Company{
		ID: "00000000-0000-0000-0000-000000000000",
	}
	admin := User{
		ID:      "00000000-0000-0000-0000-000000000001",
		Company: company,
		Role:    Admin,
	}
	normal := User{
		ID:      "00000000-0000-0000-0000-000000000001",
		Company: company,
		Role:    Normal,
	}
	vendor := Vendor{
		ID:      "00000000-0000-0000-0000-000000000002",
		Company: company,
	}
	vendorBankAccount := VendorBankAccount{
		ID:     "00000000-0000-0000-0000-000000000003",
		Vendor: vendor,
	}

	today := time.Now()
	tomorrow := today.Add(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	sixtyDaysLater := today.Add(60 * 24 * time.Hour)
	sixtyOneDaysLater := today.Add(61 * 24 * time.Hour)
	tests := []struct {
		name    string
		args    args
		want    *Invoice
		wantErr bool
	}{
		{"user is not admin",
			args{
				user:              normal,
				vendor:            vendor,
				vendorBankAccount: vendorBankAccount,
				paymentAmount:     1000,
				dueDate:           yesterday,
			},
			nil,
			true,
		},
		{"due date is yesterday",
			args{
				user:              admin,
				vendor:            vendor,
				vendorBankAccount: vendorBankAccount,
				paymentAmount:     1000,
				dueDate:           yesterday,
			},
			nil,
			true,
		},
		{"due date is today",
			args{
				user:              admin,
				vendor:            vendor,
				vendorBankAccount: vendorBankAccount,
				paymentAmount:     1000,
				dueDate:           today,
			},
			&Invoice{
				Company:            company,
				User:               admin,
				Vendor:             vendor,
				VendorBankAccount:  vendorBankAccount,
				PaymentAmount:      1000,
				Status:             Unprocessed,
				Fee:                40,
				FeeRate:            0.04,
				ConsumptionTax:     4.0,
				ConsumptionTaxRate: 0.10,
				BilledAmount:       1044.0,
				DueDate:            pkg.TruncateToDate(today),
				IssueDate:          pkg.TruncateToDate(today),
				CreatedBy:          admin.ID,
				UpdatedBy:          admin.ID,
			},
			false,
		},
		{"due date is tomorrow",
			args{
				user:              admin,
				vendor:            vendor,
				vendorBankAccount: vendorBankAccount,
				paymentAmount:     1000,
				dueDate:           tomorrow,
			},
			&Invoice{
				Company:            company,
				User:               admin,
				Vendor:             vendor,
				VendorBankAccount:  vendorBankAccount,
				PaymentAmount:      1000,
				Status:             Unprocessed,
				Fee:                40,
				FeeRate:            0.04,
				ConsumptionTax:     4.0,
				ConsumptionTaxRate: 0.10,
				BilledAmount:       1044.0,
				DueDate:            pkg.TruncateToDate(tomorrow),
				IssueDate:          pkg.TruncateToDate(today),
				CreatedBy:          admin.ID,
				UpdatedBy:          admin.ID,
			},
			false,
		},
		{"due date is 60 days later",
			args{
				user:              admin,
				vendor:            vendor,
				vendorBankAccount: vendorBankAccount,
				paymentAmount:     1000,
				dueDate:           sixtyDaysLater,
			},
			&Invoice{
				Company:            company,
				User:               admin,
				Vendor:             vendor,
				VendorBankAccount:  vendorBankAccount,
				PaymentAmount:      1000,
				Status:             Unprocessed,
				Fee:                40,
				FeeRate:            0.04,
				ConsumptionTax:     4.0,
				ConsumptionTaxRate: 0.10,
				BilledAmount:       1044.0,
				DueDate:            pkg.TruncateToDate(sixtyDaysLater),
				IssueDate:          pkg.TruncateToDate(today),
				CreatedBy:          admin.ID,
				UpdatedBy:          admin.ID,
			},
			false,
		},
		{
			"due date is more than 60 days later",
			args{
				user:              admin,
				vendor:            vendor,
				vendorBankAccount: vendorBankAccount,
				paymentAmount:     1000,
				dueDate:           sixtyOneDaysLater,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInvoice(tt.args.user, tt.args.vendor, tt.args.vendorBankAccount, tt.args.paymentAmount, tt.args.dueDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInvoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			opts := []cmp.Option{
				cmpopts.IgnoreFields(Invoice{}, "ID", "CreatedAt", "UpdatedAt"),
			}

			if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
				t.Errorf("diff: -want, +got:\n%s", diff)
			}
		})
	}
}
