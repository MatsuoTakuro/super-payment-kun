package pkg

type APICode string

const (
	// 0xxxxx: common error
	Success          APICode = "000000"
	Unauthorized     APICode = "000001"
	PermissionDenied APICode = "000002"
	Unknown          APICode = "099999"
	// 1xxxxx: request/response error
	DecodeReqBodyFailed APICode = "100000"

	// 2xxxxx: business logic error
	DueDateIsPassed APICode = "200000"

	// 3xxxxx: data access error
	RegisterDataFailed APICode = "300000"
)
