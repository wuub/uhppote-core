package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Task struct {
	Task     TaskType `json:"task"`
	Door     uint8    `json:"door,omitempty"`
	From     *Date    `json:"start-date,omitempty"`
	To       *Date    `json:"end-date,omitempty"`
	Weekdays Weekdays `json:"weekdays,omitempty"`
	Start    HHmm     `json:"start,omitempty"`
	Cards    uint8    `json:"cards,omitempty"`
}

type TaskType int

const (
	DoorControlled TaskType = iota
	DoorOpen
	DoorClosed
	DisableTimeProfile
	EnableTimeProfile
	CardNoPassword
	CardInPassword
	CardInOutPassword
	EnableMoreCards
	DisableMoreCards
	TriggerOnce
	DisablePushButton
	EnablePushButton
)

func (tt TaskType) String() string {
	return [...]string{
		"CONTROL DOOR",
		"UNLOCK DOOR",
		"LOCK DOOR",
		"DISABLE TIME PROFILE",
		"ENABLE TIME PROFILE",
		"ENABLE CARD, NO PASSWORD",
		"ENABLE CARD+IN PASSWORD",
		"ENABLE CARD+PASSWORD",
		"ENABLE MORE CARDS",
		"DISABLE MORE CARDS",
		"TRIGGER ONCE",
		"DISABLE PUSH BUTTON",
		"ENABLE PUSH BUTTON",
	}[tt]
}

func (t Task) String() string {
	from := fmt.Sprintf("%v", t.From)
	to := fmt.Sprintf("%v", t.To)
	dates := ""

	if from != "" && to != "" {
		dates = from + ":" + to
	} else if from != "" {
		dates = from + ":-"
	} else {
		dates = "-:" + to
	}

	weekdays := fmt.Sprintf("%v", t.Weekdays)
	start := fmt.Sprintf("%v", t.Start)
	door := fmt.Sprintf("%v", t.Door)
	task := fmt.Sprintf("%v", t.Task)
	cards := ""

	if t.Task == EnableMoreCards {
		fmt.Sprintf("%v", t.Cards)
	}

	list := []string{}
	for _, s := range []string{task, door, dates, weekdays, start, cards} {
		if s != "" {
			list = append(list, s)
		}
	}

	return strings.Join(list, " ")
}

func (t *Task) UnmarshalJSON(bytes []byte) error {
	task := struct {
		From     *Date    `json:"start-date,omitempty"`
		To       *Date    `json:"end-date,omitempty"`
		Weekdays Weekdays `json:"weekdays,omitempty"`
		Start    HHmm     `json:"start,omitempty"`
		Door     uint8    `json:"door,omitempty"`
		Task     TaskType `json:"task"`
		Cards    uint8    `json:"cards,omitempty"`
	}{
		Weekdays: Weekdays{},
	}

	if err := json.Unmarshal(bytes, &task); err != nil {
		return err
	}

	t.From = task.From
	t.To = task.To
	t.Weekdays = task.Weekdays
	t.Start = task.Start
	t.Door = task.Door
	t.Task = task.Task
	t.Cards = task.Cards

	return nil
}
