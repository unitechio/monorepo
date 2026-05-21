package constants

// HTTP Status Codes & Responses
const (
	ResponseSuccess = "SUCCESS"
	ResponseError   = "ERROR"

	MsgSuccess      = "Operation successful"
	MsgCreated      = "Resource created successfully"
	MsgUpdated      = "Resource updated successfully"
	MsgDeleted      = "Resource deleted successfully"
	MsgNotFound     = "Resource not found"
	MsgBadRequest   = "Bad request, invalid data"
	MsgUnauthorized = "Unauthorized access"
	MsgForbidden    = "Forbidden access"
	MsgInternalErr  = "Internal server error"
)

// HTTP Custom Error Codes (optional usage below standard HTTP)
const (
	ErrCodeValidation = 40001
	ErrCodeDuplicate  = 40901
)

// Model constraints and constants
const (
	MaxTitleLength = 500
	MaxNameLength  = 200
	MinPasswordLen = 6
	MaxEmailLength = 255
)

// Common Context Keys
const (
	ContextUserID = "user_id"
	ContextEmail  = "email"
)
