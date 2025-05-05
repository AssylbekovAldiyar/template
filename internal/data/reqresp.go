package data

import "template/pkg/reqresp"

type SaveBookRequest struct {
	// meta if needed, such as
	RequestID string
	Language  string
	UserID    string

	// main req
	Req reqresp.SaveBookRequest
}
