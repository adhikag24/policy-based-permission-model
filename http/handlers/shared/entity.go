package shared

type Response[T any] struct {
	Code    int      `json:"code"`
	Message string   `json:"message,omitempty"`
	Errors  []Errors `json:"errors,omitempty"`
	Data    T        `json:"data,omitempty"`
}

type Errors struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CommonRequest[T any] struct {
	Data T `json:"data"`
}
