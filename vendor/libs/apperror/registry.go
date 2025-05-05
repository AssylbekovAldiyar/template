package apperror

var errorRegistry = map[uint32]ErrorDef{
	// common errors: 0-999
	0: CommonErrInternal,
	1: CommonErrUnknown,
	2: CommonErrTimeout,

	// backend/platform 1000-1999
	1001: PlatformErrValidation,
	1002: PlatformErrInvalidFormat,
	1003: PlatformErrUserNotFound,
	1004: PlatformErrUserAlreadyExists,
	1005: PlatformErrUnauthorized,
	1006: PlatformErrTokenExpired,
	1007: PlatformErrDBConnection,

	1008: MsgErrSendFailed,
	1009: MsgErrAuthFailed,
	1010: MsgErrInvalidRequest,
	1011: MsgErrStatusCheckFailed,
	1012: MsgErrProviderUnavailable,
	1013: MsgErrContentInvalid,
	1014: MsgErrSessionExpired,

	// backend/sales 2000-2999
	// backend/core 3000-3999
}

func GetError(code uint32) (*AppError, bool) {
	err, ok := errorRegistry[code]
	return New(err), ok
}
