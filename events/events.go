package events

import (
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Title       string
	Date        time.Time
	Description string
}

func parseEvent(n *html.Node) (*Event, error) {
	var event Event
	title, ok := scrape.Find(n, scrape.ByClass("event-title-many"))
	if ok {
		event.Title = scrape.Text(title)
	}
	date, ok := scrape.Find(n, scrape.ByClass("event-date"))
	if ok {
		var err error
		datestr := scrape.Text(date)
		datestr = strings.Replace(datestr, "th ", " ", 1)
		datestr = strings.Replace(datestr, "nd ", " ", 1)
		datestr = strings.Replace(datestr, "rd ", " ", 1)
		datestr = strings.Replace(datestr, "st ", " ", 1)
		event.Date, err = time.Parse("Monday 2 January", datestr)
		if err != nil {
			return nil, err
		}
		event.Date = event.Date.AddDate(time.Now().Year(), 0, 0)
	}
	description, ok := scrape.Find(n, scrape.ByClass("unstyled"))
	if ok {
		event.Description = scrape.Text(description)
	}
	return &event, nil
}

func pagination(n *html.Node) (int, error) {
	articles := scrape.FindAll(n, scrape.ByTag(atom.Span))
	var text []string
	for _, article := range articles {
		text = append(text, scrape.Text(article))
	}
	return strconv.Atoi(strings.Split(strings.Join(text, " "), " ")[2])
}

func getPages() ([]string, error) {
	resp, err := http.Get("http://thelanesbristol.skiddletickets.com/events.php")
	if err != nil {
		return nil, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	pages, ok := scrape.Find(root, scrape.ByClass("pagination"))
	var numEvents int
	if ok {
		numEvents, err = pagination(pages)
		if err != nil {
			return nil, err
		}
	}
	var urls []string
	for page := 0; page < numEvents; page += 10 {
		urls = append(urls, "http://thelanesbristol.skiddletickets.com/events.php?&page="+strconv.Itoa(page))
	}
	return urls, nil
}

func GetLanesEvents() ([]*Event, error) {
	var events []*Event
	urls, err := getPages()
	if err != nil {
		return nil, err
	}
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		root, err := html.Parse(resp.Body)
		if err != nil {
			return nil, err
		}
		articles := scrape.FindAllNested(root, scrape.ByClass("eventli"))
		for _, article := range articles {
			event, err := parseEvent(article)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	}
	return events, nil
}
