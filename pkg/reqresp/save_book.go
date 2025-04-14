package reqresp

type SaveBookRequest struct {
	Name string `json:"name"`
}

type SaveBookResponse struct {
	Name string `json:"name"`
}
