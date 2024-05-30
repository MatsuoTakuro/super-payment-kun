package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"super-payment-kun/internal/model"
	"super-payment-kun/internal/pkg"
	"time"
)

type getInvoices struct {
	svc GetInvoicesSvc
	v   *pkg.Validator
}

func NewGetInvoices(svc GetInvoicesSvc, v *pkg.Validator) *getInvoices {
	return &getInvoices{
		svc: svc,
		v:   v,
	}
}

const (
	fromDueDateKey = "from_due_date"
	toDueDateKey   = "to_due_date"
	limitKey       = "limit"
	cursorKey      = "cursor"
	directionKey   = "direction"
)

type getInvoicesReq struct {
	FromDueDate string `validate:"required,datetime=2006-01-02"`
	ToDueDate   string `validate:"required,datetime=2006-01-02"`
	Limit       string `validate:"omitempty,number"`
	Cursor      string `validate:"omitempty,base64"`        // Cursor is the last invoice ID and due date in the previous or next direction.
	Direction   string `validate:"omitempty,oneof=bwd fwd"` // Direction is the direction to fetch invoices. "bwd" for backward, "fwd" for forward.
}

type data struct {
	Invoices []getInvoice `json:"invoices"`
	Cursor   string       `json:"cursor,omitempty"`
	HasNext  bool         `json:"has_next"` // HasNext is true if there are more invoices to be fetched in the same direction.
	HasPrev  bool         `json:"has_prev"` // HasPrev is true if there are more invoices to be fetched in the opposite direction.
}

type getInvoice struct {
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

func (h *getInvoices) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := getInvoicesReq{
		FromDueDate: r.URL.Query().Get(fromDueDateKey),
		ToDueDate:   r.URL.Query().Get(toDueDateKey),
		Limit:       r.URL.Query().Get(limitKey),
		Cursor:      r.URL.Query().Get(cursorKey),
		Direction:   r.URL.Query().Get(directionKey),
	}

	if err, msgs := h.v.Struct(req); err != nil {
		resp := pkg.NewAPIError(pkg.InvalidParams, err, msgs...)
		RespondWithStatus(ctx, w, resp, http.StatusBadRequest)
		return
	}

	from, err := time.Parse(pkg.DateLayout, req.FromDueDate)
	if err != nil {
		resp := pkg.NewAPIError(pkg.InvalidParams, err, fmt.Sprintf("failed to parse %s: %s", fromDueDateKey, req.FromDueDate))
		RespondWithStatus(ctx, w, resp, http.StatusBadRequest)
		return
	}

	to, err := time.Parse(pkg.DateLayout, req.ToDueDate)
	if err != nil {
		resp := pkg.NewAPIError(pkg.InvalidParams, err, fmt.Sprintf("failed to parse %s: %s", toDueDateKey, req.ToDueDate))
		RespondWithStatus(ctx, w, resp, http.StatusBadRequest)
		return
	}

	limit := 30 // Default limit
	if req.Limit != "" {
		limit, err = strconv.Atoi(req.Limit)
		if err != nil || limit <= 0 {
			RespondWithStatus(ctx, w, pkg.NewAPIError(pkg.InvalidParams, err, fmt.Sprintf("invalid limit: %s", req.Limit)), http.StatusBadRequest)
			return
		}
	}

	var lastID string
	var lastDueDate time.Time
	if req.Cursor != "" {
		cursor, err := DecodeCursor(req.Cursor)
		if err != nil {
			RespondWithStatus(ctx, w, pkg.NewAPIError(pkg.InvalidParams, err, fmt.Sprintf("invalid cursor: %s", req.Cursor)), http.StatusBadRequest)
			return
		}
		if cursor.LastDueDate.Before(from) {
			err := fmt.Errorf("last due date (%s) is before start due date (%s)", pkg.FormatToDate(cursor.LastDueDate), pkg.FormatToDate(from))
			RespondWithStatus(ctx, w, pkg.NewAPIError(pkg.InvalidParams, err, err.Error()), http.StatusBadRequest)
			return
		}
		lastID = cursor.LastID
		lastDueDate = cursor.LastDueDate
	}

	direction := "fwd" // Default direction
	if req.Direction != "" {
		direction = req.Direction
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
		Company: model.Company{ID: companyID},
		Role:    model.Role(role),
	}

	invoices, newLastID, newLastDueDate, hasNext, hasPrev, err := h.svc.GetInvoices(
		ctx, user,
		from, to,
		limit,
		lastID, lastDueDate,
		direction,
	)
	if err != nil {
		DeriveStatusAndRespond(ctx, w, err)
		return
	}

	if len(invoices) == 0 {
		resp := NormalResp{
			APICode: pkg.Success,
			Data: data{
				Invoices: []getInvoice{},
				HasNext:  hasNext,
				HasPrev:  hasPrev,
			},
		}
		RespondWithStatus(ctx, w, resp, http.StatusOK)
		return
	}

	cursor, err := GetInvsCursor{LastID: newLastID, LastDueDate: newLastDueDate}.Encode()
	if err != nil {
		RespondWithStatus(ctx, w, pkg.NewAPIError(pkg.Unknown, err, err.Error()), http.StatusInternalServerError)
		return
	}

	getInvsResp := make([]getInvoice, 0, len(invoices))
	for _, inv := range invoices {
		getInv := getInvoice{
			ID:                 string(inv.ID),
			CompanyID:          string(inv.Company.ID),
			VendorID:           string(inv.Vendor.ID),
			VendorName:         inv.Vendor.Name,
			VendorBankAcountID: string(inv.VendorBankAccount.ID),
			BankName:           inv.VendorBankAccount.BankName,
			BranchName:         inv.VendorBankAccount.BranchName,
			AccountNumber:      inv.VendorBankAccount.AccountNumber,
			AccountName:        inv.VendorBankAccount.AccountName,
			PaymentAmount:      inv.PaymentAmount,
			Status:             int8(inv.Status),
			StatusName:         inv.Status.ToJPName(),
			Fee:                inv.Fee,
			FeeRate:            inv.FeeRate,
			ConsumptionTax:     inv.ConsumptionTax,
			ConsumptionTaxRate: inv.ConsumptionTaxRate,
			BilledAmount:       inv.BilledAmount,
			IssueDate:          pkg.FormatToDate(inv.IssueDate),
			DueDate:            pkg.FormatToDate(inv.DueDate),
		}
		getInvsResp = append(getInvsResp, getInv)
	}

	resp := NormalResp{
		APICode: pkg.Success,
		Data: data{
			Invoices: getInvsResp,
			Cursor:   cursor,
			HasNext:  hasNext,
			HasPrev:  hasPrev,
		},
	}

	RespondWithStatus(ctx, w, resp, http.StatusOK)
}

type GetInvsCursor struct {
	LastID      string    `json:"last_id"`
	LastDueDate time.Time `json:"last_due_date"`
}

func (c GetInvsCursor) Encode() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}

func DecodeCursor(encoded string) (GetInvsCursor, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return GetInvsCursor{}, err
	}

	var cursor GetInvsCursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return GetInvsCursor{}, err
	}
	return cursor, nil
}
