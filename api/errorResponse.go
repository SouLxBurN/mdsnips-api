package api

// ErrorResponse
type ErrorResponse struct {
	Code    int32  `json:"code,omitempty" format:"int" example:"400"`
	Message string `json:"message,omitempty" format:"string"`
}
