package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

var twitterAPIKey, twitterAPISecret, twitterAccessToken, twitterAccessTokenSecret string

func loadEnvironmentVariables() {
	twitterAPIKey = os.Getenv("TWITTER_API_KEY")
	if twitterAPIKey == "" {
		log.Fatal("TWITTER_API_KEY not found")
	}
	twitterAPISecret = os.Getenv("TWITTER_API_SECRET")
	if twitterAPIKey == "" {
		log.Fatal("TWITTER_API_SECRET not found")
	}
	twitterAccessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
	if twitterAPIKey == "" {
		log.Fatal("TWITTER_ACCESS_TOKEN not found")
	}
	twitterAccessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	if twitterAPIKey == "" {
		log.Fatal("TWITTER_ACCESS_TOKEN_SECRET not found")
	}
}

func main() {

	// stop forgetting to load environment variables !!
	loadEnvironmentVariables()
	anaconda.SetConsumerKey(twitterAPIKey)
	anaconda.SetConsumerSecret(twitterAPISecret)
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterAccessTokenSecret)
	log.Println("Connected to Twitter API ...")

	user, err := api.GetSelf(nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Retrieve logged user information ...")

	stream := api.UserStream(nil)
	log.Println("User stream initialized ...")

	for t := range stream.C {
		switch v := t.(type) {
		case anaconda.Tweet:
			fmt.Printf("Tweet by: %s \n", v.User.Name)
			if isMentionned(user.ScreenName, v) && isOnBlackList(v.User) {
				log.Println("User mentionnned in tweet by blacklisted guy: Exterminate !")
			}
		}
	}
}

func isMentionned(name string, tweet anaconda.Tweet) bool {
	for _, mention := range tweet.Entities.User_mentions {

		if mention.Screen_name == name {
			return true
		}
	}
	return false
}

func isOnBlackList(tweet anaconda.User) bool {
	return true
}
