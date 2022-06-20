package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	json_models "httpserv/internal/json_models"
	"httpserv/internal/logger"
	event "httpserv/internal/service"
)

func ParseKey(r *http.Request, skey string) (string, error) {
	keys, ok := r.URL.Query()[skey]
	if !ok || len(keys[0]) < 1 {
		return "", fmt.Errorf("key by %v name not found", skey)
	}
	key := keys[0]
	return key, nil
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(400)
		logger.Error(r, fmt.Errorf("Method is not POST: %v", r.Method))
		return
	}

	model := json_models.EventModel{}

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		w.WriteHeader(400)
		logger.Error(r, fmt.Errorf("error on json unmarshal: %v", err))
		return
	}

	createdBy := r.RemoteAddr

	newDate, err := json_models.ParseDate(model.Date)
	if err != nil {
		ResponseErr(r, w, err, 400)
		logger.Error(r, fmt.Errorf("error on parse date: %v", err))
		return
	}

	res := event.Create(model.Desc, newDate, time.Now(), createdBy)

	ResponseResult(r, w, res)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	model := json_models.EventModelUpdate{}

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		ResponseErr(r, w, fmt.Errorf("error on json unmarshal: %v", err), 400)
		return
	}

	date, err := json_models.ParseDate(model.Date)
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on json unmarshal date: %v", err), 400)
		return
	}

	result, err := event.Update(model.ID, date, model.Desc)
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on update: %v", err), 503)
		return
	}

	ResponseResult(r, w, result)
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	model := json_models.EventModelDelete{}

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		ResponseErr(r, w, fmt.Errorf("error on json unmarshal: %v", err), 400)
		return
	}

	intID, err := strconv.Atoi(model.ID)
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on delete: %v", err), 400)
	}

	if err := event.Delete(intID); err != nil {
		ResponseErr(r, w, fmt.Errorf("error on delete: %v", err), 503)
		return
	}

	ResponseResult(r, w, "success")
}

func GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	result, err := parseGetData(r, "date")
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on parse query: %v", err), 400)
		return
	}

	date, err := json_models.ParseDate(result)
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on parse data: %v", err), 400)
		return
	}

	response, err := event.GetForDay(date)
	if err != nil {
		if err == event.ErrorEmptyResult {
			ResponseResult(r, w, response)
			return
		}
		ResponseErr(r, w, fmt.Errorf("error on getting: %v", err), 503)
		return
	}

	ResponseResult(r, w, response)
}

func GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	result, err := parseGetData(r, "date")
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on parse query: %v", err), 400)
		return
	}

	date, err := json_models.ParseDate(result)
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on parse data: %v", err), 400)
		return
	}

	response, err := event.GetForWeek(date)
	if err != nil {
		if err == event.ErrorEmptyResult {
			ResponseResult(r, w, response)
			return
		}
		ResponseErr(r, w, fmt.Errorf("error on getting: %v", err), 503)
		return
	}

	ResponseResult(r, w, response)
}

func GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	result, err := parseGetData(r, "date")
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on parse query: %v", err), 400)
		return
	}

	date, err := json_models.ParseDate(result)
	if err != nil {
		ResponseErr(r, w, fmt.Errorf("error on parse data: %v", err), 400)
		return
	}

	response, err := event.GetForMonth(date)
	if err != nil {
		if err == event.ErrorEmptyResult {
			ResponseResult(r, w, response)
			return
		}
		ResponseErr(r, w, fmt.Errorf("error on getting: %v", err), 503)
		return
	}

	ResponseResult(r, w, response)
}
