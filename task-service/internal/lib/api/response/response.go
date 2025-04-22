package response

import "net/http"

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func StatusCreated() Response {
	return Response{
		Status: http.StatusCreated,
	}
}

func StatusOK() Response {
	return Response{
		Status: http.StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: http.StatusInternalServerError,
		Error:  msg,
	}
}

func ErrorClient(msg string) Response {
	return Response{
		Status: http.StatusBadRequest,
		Error:  msg,
	}
}
