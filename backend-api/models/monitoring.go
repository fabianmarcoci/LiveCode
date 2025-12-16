package models

type ClientErrorLog struct {
	Timestamp    string `json:"timestamp" binding:"required"`
	ErrorType    string `json:"error_type" binding:"required,max=100"`
	ErrorMessage string `json:"error_message" binding:"required,max=1000"`
	AppVersion   string `json:"app_version" binding:"required,max=50"`
	OS           string `json:"os" binding:"required,max=20"`
}
