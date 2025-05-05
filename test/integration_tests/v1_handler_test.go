//go:build integration
// +build integration

package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"libs/common/dbctx"
	"libs/common/service_middleware"

	store2 "template/internal/app/store"
	http2 "template/internal/deliveries/book/http/v1"
	bookpg "template/internal/repositories/book/pg"
	"template/internal/services/book"
	"template/pkg/reqresp"
)

func TestBookHanlder_Integration(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	pgContext := dbctx.New(testDB.DB)
	store := store2.RepositoryStore{PgContext: pgContext, BookRepository: bookpg.NewPgRepository(pgContext)}
	store := store2.ClientStore{} // todo: populate with values
	bookHandler := http2.NewV1(book.New(&store, service_middleware.NewServiceMiddleware()), 0)

	t.Run("Create Book", func(t *testing.T) {
		createReq := reqresp.SaveBookRequest{
			Name: "Jumanji",
		}
		body, _ := json.Marshal(createReq)

		req := httptest.NewRequest(http.MethodPost, "/v1/book", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := bookHandler.SaveBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
