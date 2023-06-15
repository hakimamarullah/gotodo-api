package models

type ResponseBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count"`
	Code    int         `json:"code"`
}
