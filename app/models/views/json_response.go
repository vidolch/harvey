package views

type JsonResponse struct {
	Error string `json:"error""`
	Data interface{} `json:"data""`
}
