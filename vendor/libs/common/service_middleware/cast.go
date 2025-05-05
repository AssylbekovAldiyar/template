package service_middleware

import "errors"

// SafeCast безопасно приводит interface{} к нужному типу,
// обрабатывая nil значения и возвращая zero value в случае ошибок
func SafeCast[T any](value interface{}, err error) (T, error) {
	var zero T

	if err != nil {
		return zero, err
	}

	if value == nil {
		return zero, nil
	}

	result, ok := value.(T)
	if !ok {
		return zero, errors.New("failed to cast")
	}

	return result, nil
}
