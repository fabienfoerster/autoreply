package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

type Configuration struct {
	Users []struct {
		Name     string
		Response struct {
			Text     string
			MediaURL string
		}
	}
}

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

func loadConfiguration() Configuration {
	log.Println("Loading configuration ...")
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	return configuration
}

func main() {

	configuration := loadConfiguration()

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
		switch tw := t.(type) {
		case anaconda.Tweet:
			fmt.Printf("Tweet by: %s \n", tw.User.Name)
			if isMentionned(user.ScreenName, tw) && isOnBlackList(tw.User, configuration) {
				log.Println("User mentionnned in tweet by blacklisted guy: Exterminate !")
				v := url.Values{}
				v.Set("in_reply_to_status_id", tw.IdStr)
				var reponse = fmt.Sprintf("@%s Réponse automatique et boum (j'espère) cc @haitaar", tw.User.ScreenName)
				api.PostTweet(reponse, v)
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

func isOnBlackList(u anaconda.User, config Configuration) bool {
	for _, user := range config.Users {
		if user.Name == u.ScreenName {
			return true
		}
	}
	return false
}
