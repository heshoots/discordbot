package twitter

import (
	"fmt"
	twt "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
)

type TwitterAuth struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func getClient(auth TwitterAuth) *twt.Client {
	twitterconfig := oauth1.NewConfig(auth.ConsumerKey, auth.ConsumerSecret)
	token := oauth1.NewToken(auth.AccessToken, auth.AccessSecret)
	httpClient := twitterconfig.Client(oauth1.NoContext, token)
	return twt.NewClient(httpClient)
}

func Tweet(auth TwitterAuth, message string) (string, error) {
	log.Println("tweeting")
	client := getClient(auth)
	tweet, _, err := client.Statuses.Update(message, nil)
	if err != nil {
		log.Println("could not tweet,", err)
		return "", err
	}
	return "http://twitter.com/sodiumshowdown/status/" + fmt.Sprint(tweet.ID), nil
}
