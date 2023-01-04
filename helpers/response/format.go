package response

type Base struct {
	Status  int         `json:"-"`
	Message string      `json:"message"`
	Data    interface{} `json:"-"`
}
