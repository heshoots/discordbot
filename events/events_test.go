package events

import (
	"fmt"
	"testing"
	"time"
)

func Test_getEvents(context *testing.T) {
	events, err := GetLanesEvents()
	if err != nil {
		context.Fatal(err)
		return
	}
	for _, event := range events {
		if event.Date.Weekday() == time.Sunday {
			fmt.Println(event)
		}
	}
}
