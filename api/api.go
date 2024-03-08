package api

type (
	GenerateZipRequest struct {
		FileName string `validate:"required,min=5,max=20,fileName"`
		Password string `validate:"required,min=5,max=20"`
		Content  string `validate:"required,max=100"`
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)
