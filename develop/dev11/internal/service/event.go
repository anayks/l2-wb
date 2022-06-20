package event

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Event struct {
	ID          int
	Description string
	Date        time.Time
	CreatedAt   time.Time
	CreatedBy   string
}

var ErrorEmptyResult = errors.New("Day data doesn't exists")

var incr int

var DayEvents map[time.Time][]*Event = map[time.Time][]*Event{}
var WeekEvents map[string][]*Event = map[string][]*Event{}
var MonthEvents map[string][]*Event = map[string][]*Event{}
var EventsMap map[int]*Event = map[int]*Event{}

var rw sync.RWMutex

func Create(description string, date, createdAt time.Time, createdBy string) Event {
	rw.Lock()
	defer rw.Unlock()
	newEvent := Event{
		ID:          incr,
		Description: description,
		Date:        date,
		CreatedAt:   createdAt,
		CreatedBy:   createdBy,
	}
	EventsMap[incr] = &newEvent

	incr++

	addToDay(date, &newEvent)
	addToWeek(date, &newEvent)
	addToYear(date, &newEvent)
	return newEvent
}

func addToDay(date time.Time, e *Event) {
	eDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	_, ok := DayEvents[eDate]
	if !ok {
		DayEvents[eDate] = make([]*Event, 0)
	}
	DayEvents[eDate] = append(DayEvents[eDate], e)
}

func addToWeek(date time.Time, e *Event) {
	eDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	year, week := eDate.ISOWeek()
	key := fmt.Sprintf("%v_%v", year, week)

	_, ok := WeekEvents[key]
	if !ok {
		WeekEvents[key] = make([]*Event, 0)
	}
	WeekEvents[key] = append(WeekEvents[key], e)
}

func addToYear(date time.Time, e *Event) {
	eDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	month := eDate.Month()
	year := eDate.Year()
	key := fmt.Sprintf("%v_%v", year, month)

	_, ok := MonthEvents[key]
	if !ok {
		MonthEvents[key] = make([]*Event, 0)
	}
	MonthEvents[key] = append(MonthEvents[key], e)
}

func UpdateDescription(id int, desc string) error {
	_, ok := EventsMap[id]
	if !ok {
		return fmt.Errorf("Error on update desc by id #%v.", id)
	}
	EventsMap[id].Description = desc
	return nil
}

func UpdateDate(id int, newDate time.Time) error {
	event, ok := EventsMap[id]
	if !ok {
		return fmt.Errorf("Event by id #%v doesn't exists", id)
	}
	date := event.Date
	err := deleteFromMonth(date, id)
	if err != nil {
		return fmt.Errorf("Error on update date: %v", err)
	}

	err = deleteFromWeek(date, id)
	if err != nil {
		return fmt.Errorf("Error on update date: %v", err)
	}

	err = deleteFromDay(date, id)
	if err != nil {
		return fmt.Errorf("Error on update date: %v", err)
	}
	EventsMap[id].Date = newDate
	addToDay(date, EventsMap[id])
	addToWeek(date, EventsMap[id])
	addToYear(date, EventsMap[id])
	return nil
}

func Update(id int, date time.Time, desc string) (Event, error) {
	rw.Lock()
	defer rw.Unlock()

	ev, ok := EventsMap[id]
	if !ok {
		return Event{}, fmt.Errorf("on update desc by id #%v.", id)
	}

	var err error

	if desc != "" {
		err = UpdateDescription(id, desc)
	}

	if err != nil {
		return Event{}, err
	}

	if (date == time.Time{}) {
		err = UpdateDate(id, date)
	}

	return *ev, err
}

func Delete(id int) error {
	rw.Lock()
	defer rw.Unlock()
	event, ok := EventsMap[id]
	if !ok {
		return fmt.Errorf("Event by id #%v doesn't exists", id)
	}

	date := event.Date
	err := deleteFromMonth(date, id)
	if err != nil {
		return err
	}

	err = deleteFromWeek(date, id)
	if err != nil {
		return err
	}

	err = deleteFromDay(date, id)
	if err != nil {
		return err
	}
	delete(EventsMap, id)
	return nil
}

func deleteFromMonth(date time.Time, id int) error {
	month := date.Month()
	year := date.Year()
	key := fmt.Sprintf("%v_%v", year, month)
	mth := MonthEvents[key]
	for vid, v := range mth {
		if v.ID != id {
			continue
		}
		if id == 0 {
			MonthEvents[key] = MonthEvents[key][1:]
		} else if id == len(MonthEvents[key])-1 {
			MonthEvents[key] = MonthEvents[key][:len(mth)-1]
		} else {
			MonthEvents[key] = append(MonthEvents[key][0:vid], MonthEvents[key][vid+1:]...)
		}
		return nil
	}
	return fmt.Errorf("element with id #%v not found in MonthEvents", id)
}

func deleteFromWeek(date time.Time, id int) error {
	_, week := date.ISOWeek()
	key := fmt.Sprintf("%v_%v", date.Year(), week)
	mth := WeekEvents[key]
	for vid, v := range mth {
		if v.ID != id {
			continue
		}
		if id == 0 {
			WeekEvents[key] = WeekEvents[key][1:]
		} else if id == len(MonthEvents[key])-1 {
			WeekEvents[key] = WeekEvents[key][:len(mth)-1]
		} else {
			WeekEvents[key] = append(WeekEvents[key][0:vid], WeekEvents[key][vid+1:]...)
		}
		return nil
	}
	return fmt.Errorf("element with id #%v not found in WeekEvents", id)
}

func deleteFromDay(date time.Time, id int) error {
	grp := DayEvents[date]
	key := date
	for vid, v := range grp {
		if v.ID != id {
			continue
		}
		if id == 0 {
			DayEvents[key] = DayEvents[key][1:]
		} else if id == len(DayEvents[key])-1 {
			DayEvents[key] = DayEvents[key][:len(grp)-1]
		} else {
			DayEvents[key] = append(DayEvents[key][0:vid], DayEvents[key][vid+1:]...)
		}
		return nil
	}
	return fmt.Errorf("element with id #%v not found in DayEvents", id)
}

func GetForDay(d time.Time) ([]Event, error) {
	rw.RLock()
	defer rw.RUnlock()

	var result []Event = []Event{}
	v, ok := DayEvents[d]
	if !ok {
		return result, ErrorEmptyResult
	}
	for _, value := range v {
		result = append(result, *value)
	}
	return result, nil
}

func GetForWeek(d time.Time) ([]Event, error) {
	rw.RLock()
	defer rw.RUnlock()

	var result []Event = []Event{}
	year, week := d.ISOWeek()
	key := fmt.Sprintf("%v_%v", year, week)
	v, ok := WeekEvents[key]
	if !ok {
		return result, ErrorEmptyResult
	}
	for _, value := range v {
		result = append(result, *value)
	}
	return result, nil
}

func GetForMonth(d time.Time) ([]Event, error) {
	rw.RLock()
	defer rw.RUnlock()

	var result []Event = []Event{}
	year, month := d.Year(), d.Month()
	key := fmt.Sprintf("%v_%v", year, month)
	v, ok := MonthEvents[key]
	if !ok {
		return result, ErrorEmptyResult
	}
	for _, value := range v {
		result = append(result, *value)
	}
	return result, nil
}
