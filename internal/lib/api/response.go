package api

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	statusOK    = "OK"
	statusERROR = "Error"
)

func OK() Response {
	return Response{
		Status: statusOK,
	}
}
func Error(msg string) Response {
	return Response{
		Status: statusERROR,
		Error:  msg,
	}
}
