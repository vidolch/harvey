package views

type JsonResponse struct {
	Error string `json:"error""`
	Name string `json:"name""`
	Id string `json:"id"`
	Data interface{}
}
