package deliveries

import (
	"libs/apperror"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleEcho(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	if appErr := apperror.AsError(err); appErr != nil {
		return c.JSON(appErr.ErrorDef.HTTP, appErr)
	}

	return c.JSON(http.StatusInternalServerError, apperror.New(apperror.CommonErrUnknown).Wrap(err))
}
