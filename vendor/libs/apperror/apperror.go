package apperror

import (
	"errors"
	"fmt"
)

// AppError представляет расширенную структуру ошибки
type AppError struct {
	ErrorDef
	err      error                  // приватное поле для хранения ошибок
	Metadata map[string]interface{} `json:"metadata,omitempty"` // метаданные ошибки
}

// Builder для создания ошибок
func New(errDef ErrorDef) *AppError {
	return &AppError{
		ErrorDef: errDef,
		Metadata: make(map[string]interface{}),
	}
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%d]: %s - %v", e.ErrorDef.Code, e.ErrorDef.Message, e.err)
	}
	return fmt.Sprintf("[%d]: %s", e.ErrorDef.Code, e.ErrorDef.Message)
}

// Unwrap реализует интерфейс для errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.err
}

// Is делает проверку на равенство кодов ошибок, если target является AppError
// или делегирует проверку вложенной ошибки
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.ErrorDef.Code == t.ErrorDef.Code
	}

	// Проверяем вложенные ошибки
	return errors.Is(e.err, target)
}

// As преобразует ошибку в AppError, если target является указателем на AppError
// или делегирует преобразование вложенной ошибки
func (e *AppError) As(target interface{}) bool {
	if appErr, ok := target.(**AppError); ok {
		*appErr = e
		return true
	}

	// Проверяем вложенные ошибки
	return errors.As(e.err, target)
}

// WithMessage устанавливает сообщение об ошибке
func (e *AppError) WithMessage(message string) *AppError {
	e.ErrorDef.Message = message
	return e
}

// WithMetadata добавляет метаданные к ошибке
func (e *AppError) WithMetadata(key string, value interface{}) *AppError {
	e.Metadata[key] = value
	return e
}

// Wrap оборачивает ошибку с сохранением цепочки
func (e *AppError) Wrap(err error) *AppError {
	if err == nil {
		return e
	}
	if e.err == nil {
		e.err = err
	} else {
		e.err = errors.Join(e.err, err)
	}
	return e
}

// IsErrorCode проверяет код ошибки
func IsErrorCode(err error, errDef ErrorDef) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.ErrorDef.Code == errDef.Code
	}
	return false
}

// AsError преобразует ошибку в AppError
func AsError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return nil
}
