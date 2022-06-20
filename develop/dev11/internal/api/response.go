package http

import (
	"encoding/json"
	"httpserv/internal/logger"
	"net/http"
)

type ErrorResponse struct {
	Err string `json:"error"`
}

func ResponseErr(r *http.Request, w http.ResponseWriter, v error, code int) {
	w.WriteHeader(code)
	var err string
	if code == 503 {
		err = "server internal error"
	} else {
		err = "bad request"
	}
	response := ErrorResponse{
		Err: err,
	}
	logger.Error(r, v)
	res, _ := json.Marshal(response)
	w.Write(res)
}

type ResultResponse struct {
	Result interface{} `json:"result"`
}

func ResponseResult(r *http.Request, w http.ResponseWriter, data interface{}) {
	w.WriteHeader(200)
	response := ResultResponse{
		Result: data,
	}
	logger.Log(r, "responsed")
	res, _ := json.Marshal(response)
	w.Write(res)
}
