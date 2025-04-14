package book

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"

	"libs/apperror"
	"template/pkg/domain"
)

const (
	bookServiceAddress = "order.svc.cluster.local"

	getBookById = bookServiceAddress + "/v1/book/{id}"
)

type HTTPClient interface {
	GetBookByID(ctx context.Context, checkoutID string) (domain.Book, error)
}

type httpClient struct {
	r *resty.Client
}

func NewHTTPClient(r *http.Client) HTTPClient {
	return &httpClient{r: resty.NewWithClient(r)}
}

func (c *httpClient) GetBookByID(ctx context.Context, checkoutID string) (domain.Book, error) {
	resp, err := c.r.R().SetPathParam("id", checkoutID).
		SetHeader("Session-Id", ctx.Value("session_id").(string)).
		Get(getBookById)
	if err != nil {
		return domain.Book{}, apperror.New(apperror.CommonErrInternal)
	}

	var book domain.Book

	err = json.Unmarshal(resp.Body(), &book)
	if err != nil {
		return domain.Book{}, apperror.New(apperror.CommonErrInternal)
	}

	return book, nil
}
