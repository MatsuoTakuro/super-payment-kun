package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"super-payment-kun/internal/model"
	"super-payment-kun/internal/pkg"
	"time"
)

type createInvoice struct {
	svc CreateInvoiceSvc
	v   *pkg.Validator
}

func NewCreateInvoice(svc CreateInvoiceSvc, v *pkg.Validator) *createInvoice {
	return &createInvoice{
		svc: svc,
		v:   v,
	}
}

type createInvoiceReq struct {
	VendorID            string  `json:"vendor_id" validate:"required,uuid"` // TODO: check why `uuid4` tag does not pass (https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Universally_Unique_Identifier_UUID_v4).
	VendorBankAccountID string  `json:"vendor_bank_account_id" validate:"required,uuid"`
	PaymentAmount       float64 `json:"payment_amount" validate:"required"`
	DueDate             string  `json:"due_date" validate:"required,datetime=2006-01-02"`
}

type createInvoiceResp struct {
	ID                 string  `json:"id"`
	CompanyID          string  `json:"company_id"`
	VendorID           string  `json:"vendor_id"`
	VendorName         string  `json:"vendor_name"`
	VendorBankAcountID string  `json:"vendor_bank_account_id"`
	BankName           string  `json:"bank_name"`
	BranchName         string  `json:"branch_name"`
	AccountNumber      string  `json:"account_number"` // TODO: need to mask account number?
	AccountName        string  `json:"account_name"`
	PaymentAmount      float64 `json:"payment_amount"`
	Status             int8    `json:"status"`
	StatusName         string  `json:"status_name"`
	Fee                float64 `json:"fee"`
	FeeRate            float64 `json:"fee_rate"`
	ConsumptionTax     float64 `json:"consumption_tax"`
	ConsumptionTaxRate float64 `json:"consumption_tax_rate"`
	BilledAmount       float64 `json:"billed_amount"`
	IssueDate          string  `json:"issue_date"`
	DueDate            string  `json:"due_date"`
}

func (h *createInvoice) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &createInvoiceReq{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		resp := pkg.NewAPIError(pkg.DecodeReqBodyFailed, err, fmt.Sprintf("failed to decode request body: %q", r.Body))
		RespondWithStatus(ctx, w, resp, http.StatusBadRequest)
		return
	}

	if err, msgs := h.v.Struct(req); err != nil {
		resp := pkg.NewAPIError(pkg.InvalidParams, err, msgs...)
		RespondWithStatus(ctx, w, resp, http.StatusBadRequest)
		return
	}

	if !pkg.IsValidUUID(req.VendorID) {
		err := fmt.Errorf("invalid vendor_id: %s", req.VendorID)
		resp := pkg.NewAPIError(pkg.InvalidParams, err, err.Error())
		RespondWithStatus(ctx, w, resp, http.StatusBadRequest)
		return
	}

	if !pkg.IsValidUUID(req.VendorBankAccountID) {
		err := pkg.NewAPIError(pkg.InvalidParams, fmt.Errorf("invalid vendor_bank_account_id: %s", req.VendorBankAccountID))
		resp := pkg.NewAPIError(pkg.InvalidParams, err, err.Error())
		RespondWithStatus(ctx, w, resp, http.StatusBadRequest)
		return
	}

	dueDate, err := time.Parse(pkg.DateLayout, req.DueDate)
	if err != nil {
		resp := pkg.NewAPIError(pkg.InvalidParams, err, fmt.Sprintf("failed to parse due_date: %s", req.DueDate))
		RespondWithStatus(ctx, w, resp, http.StatusBadRequest)
		return
	}

	userID, err := pkg.GetUserID(ctx)
	if err != nil {
		RespondWithStatus(ctx, w, err, http.StatusInternalServerError)
		return
	}

	companyID, err := pkg.GetCompanyID(ctx)
	if err != nil {
		RespondWithStatus(ctx, w, err, http.StatusInternalServerError)
		return
	}

	role, err := pkg.GetRole(ctx)
	if err != nil {
		RespondWithStatus(ctx, w, err, http.StatusInternalServerError)
		return
	}

	user := model.User{
		ID:      userID,
		Company: model.Company{ID: companyID},
		Role:    model.Role(role),
	}

	vendor := model.Vendor{
		ID:      req.VendorID,
		Company: model.Company{ID: companyID},
	}

	vendorBankAccount := model.VendorBankAccount{
		ID:     req.VendorBankAccountID,
		Vendor: vendor,
	}

	invoice, gotVendor, gotVBankAccount, err := h.svc.CreateInvoice(ctx, user, vendor, vendorBankAccount, req.PaymentAmount, dueDate)
	if err != nil {
		DeriveStatusAndRespond(ctx, w, err)
		return
	}

	resp := NormalResp{
		APICode: pkg.Success,
		Data: &createInvoiceResp{
			ID:                 string(invoice.ID),
			CompanyID:          string(invoice.Company.ID),
			VendorID:           string(invoice.Vendor.ID),
			VendorName:         gotVendor.Name,
			VendorBankAcountID: string(invoice.VendorBankAccount.ID),
			BankName:           gotVBankAccount.BankName,
			BranchName:         gotVBankAccount.BranchName,
			AccountNumber:      gotVBankAccount.AccountNumber,
			AccountName:        gotVBankAccount.AccountName,
			PaymentAmount:      invoice.PaymentAmount,
			Status:             int8(invoice.Status),
			StatusName:         invoice.Status.ToJPName(),
			Fee:                invoice.Fee,
			FeeRate:            invoice.FeeRate,
			ConsumptionTax:     invoice.ConsumptionTax,
			ConsumptionTaxRate: invoice.ConsumptionTaxRate,
			BilledAmount:       invoice.BilledAmount,
			IssueDate:          pkg.FormatToDate(invoice.IssueDate),
			DueDate:            pkg.FormatToDate(invoice.DueDate),
		},
	}

	RespondWithStatus(ctx, w, resp, http.StatusCreated)
}
