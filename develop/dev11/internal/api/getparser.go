package http

import (
	"fmt"
	"net/http"
)

func parseGetData(r *http.Request, key string) (string, error) {
	keys, ok := r.URL.Query()["date"]

	if !ok || len(keys[0]) < 1 {
		return "", fmt.Errorf("key %v doesn't exists", key)
	}

	result := keys[0]
	return result, nil
}
