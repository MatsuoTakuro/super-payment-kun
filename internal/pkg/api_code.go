package pkg

type APICode string

const (
	// 0xxxxx: common error
	Success          APICode = "000000"
	Unauthorized     APICode = "000001"
	PermissionDenied APICode = "000002"
	Unknown          APICode = "099999"
	// 1xxxxx: request/response error
	DecodeReqBodyFailed  APICode = "100000"
	InvalidParams        APICode = "100001"
	EncodeRespBodyFailed APICode = "100002"
	WriteRespBodyFailed  APICode = "100003"
	// 2xxxxx: business logic error
	DueDateIsPassed           APICode = "200000"
	DueDateExceedMaxDeferDate APICode = "200001"
	PaymentAmountTooSmall     APICode = "200002"
	// 3xxxxx: data access error
	RegisterDataFailed              APICode = "300000"
	RegisterDuplicateDataRestricted APICode = "300001"
	SelectDataFailed                APICode = "300002"
)
