package model

// ApiResponse
type ApiResponse struct {
	Code    int32  `json:"code,omitempty" format:"int" example:"400"`
	Message string `json:"message,omitempty" format:"string"`
}
