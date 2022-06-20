package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RequestLog struct {
	Time       time.Time
	RemoteAddr string
	Method     string
	Header     string
	UserAgent  string
}

func Error(r *http.Request, log error) {
	rlog := r.Context().Value("request")
	v, _ := json.Marshal(rlog.(RequestLog))
	fmt.Printf("\nRequest <%v> error: %v", string(v), log)
}

func Log(r *http.Request, log string) {
	rlog := r.Context().Value("request")
	v, _ := json.Marshal(rlog.(RequestLog))
	fmt.Printf("\nRequest <%v>: %v", string(v), log)
}
