package response

type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []Error     `json:"errors,omitempty"`
}

type Error struct {
	Params  string
	Message string
}
