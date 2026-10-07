package dto

type Pagination struct {
  Page      int   `json:"page"`
  Limit     int   `json:"limit"`
  TotalData int64 `json:"totalData"`
  TotalPage int   `json:"totalPage"`
}

type ErrorResponse struct {
  Code    string            `json:"code"`    // VALIDATION_ERROR, CONFLICT, UNAUTHORIZED, dst.
  Message string            `json:"message"`
  Errors  map[string]string `json:"errors,omitempty"`
}
