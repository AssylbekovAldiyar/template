package apperror

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// ErrorDef содержит определение ошибки и её представление в разных протоколах
type ErrorDef struct {
	Code    uint32     `json:"code"`    // Уникальный код ошибки, также используется как код ошибки в Kafka header
	HTTP    int        `json:"http"`    // HTTP статус код
	GRPC    codes.Code `json:"grpc"`    // gRPC код
	Message string     `json:"message"` // Человекочитаемое сообщение
}

// NewErrorDef Функция-хелпер для создания ErrorDef
func NewErrorDef(code uint32, http int, grpc codes.Code, message string) ErrorDef {
	return ErrorDef{
		Code:    code,
		HTTP:    http,
		GRPC:    grpc,
		Message: message,
	}
}

// Определение всех ошибок системы
// Каждому домену будет выделяться диапазон кодов ошибок: 1000-1999, 2000-2999 и т.д.

// common errors: 0-999
var (
	CommonErrUnknown  = NewErrorDef(0, http.StatusInternalServerError, codes.Unknown, "UNKNOWN_ERROR")
	CommonErrInternal = NewErrorDef(1, http.StatusInternalServerError, codes.Internal, "INTERNAL_ERROR")
	CommonErrTimeout  = NewErrorDef(2, http.StatusGatewayTimeout, codes.DeadlineExceeded, "TIMEOUT_ERROR")
)

// backend/platform 1000-1999
var (
	PlatformErrValidation        = NewErrorDef(1001, http.StatusBadRequest, codes.InvalidArgument, "VALIDATION_ERROR")
	PlatformErrInvalidFormat     = NewErrorDef(1002, http.StatusBadRequest, codes.InvalidArgument, "INVALID_FORMAT")
	PlatformErrUserNotFound      = NewErrorDef(1003, http.StatusNotFound, codes.NotFound, "USER_NOT_FOUND")
	PlatformErrUserAlreadyExists = NewErrorDef(1004, http.StatusConflict, codes.AlreadyExists, "USER_ALREADY_EXISTS")
	PlatformErrUnauthorized      = NewErrorDef(1005, http.StatusUnauthorized, codes.Unauthenticated, "UNAUTHORIZED")
	PlatformErrTokenExpired      = NewErrorDef(1006, http.StatusUnauthorized, codes.Unauthenticated, "TOKEN_EXPIRED")
	PlatformErrDBConnection      = NewErrorDef(1007, http.StatusServiceUnavailable, codes.Unavailable, "DB_CONNECTION_ERROR")

	//notification gateways error
	MsgErrSendFailed          = NewErrorDef(1008, http.StatusServiceUnavailable, codes.Unavailable, "MSG_SEND_FAILED")
	MsgErrAuthFailed          = NewErrorDef(1009, http.StatusUnauthorized, codes.Unauthenticated, "MSG_AUTH_FAILED")
	MsgErrInvalidRequest      = NewErrorDef(1010, http.StatusBadRequest, codes.InvalidArgument, "MSG_INVALID_REQUEST")
	MsgErrStatusCheckFailed   = NewErrorDef(1011, http.StatusServiceUnavailable, codes.Unavailable, "MSG_STATUS_CHECK_FAILED")
	MsgErrProviderUnavailable = NewErrorDef(1012, http.StatusServiceUnavailable, codes.Unavailable, "MSG_PROVIDER_UNAVAILABLE")
	MsgErrContentInvalid      = NewErrorDef(1013, http.StatusBadRequest, codes.InvalidArgument, "MSG_CONTENT_INVALID")
	MsgErrSessionExpired      = NewErrorDef(1014, http.StatusUnauthorized, codes.Unauthenticated, "MSG_SESSION_EXPIRED")
)

// backend/sales 2000-2999
// backend/core 3000-3999
