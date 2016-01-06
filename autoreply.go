package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
		log.Fatalln("TWITTER_API_KEY not found")
	}
	twitterAPISecret = os.Getenv("TWITTER_API_SECRET")
	if twitterAPIKey == "" {
		log.Fatalln("TWITTER_API_SECRET not found")
	}
	twitterAccessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
	if twitterAPIKey == "" {
		log.Fatalln("TWITTER_ACCESS_TOKEN not found")
	}
	twitterAccessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	if twitterAPIKey == "" {
		log.Fatalln("TWITTER_ACCESS_TOKEN_SECRET not found")
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

	//load configuration present in the config.json file
	configuration := loadConfiguration()

	// stop forgetting to load environment variables !!
	loadEnvironmentVariables()
	//setup twitter api
	anaconda.SetConsumerKey(twitterAPIKey)
	anaconda.SetConsumerSecret(twitterAPISecret)
	api := anaconda.NewTwitterApi(twitterAccessToken, twitterAccessTokenSecret)
	log.Println("Connected to Twitter API ...")

	// retrive information about the logged user
	user, err := api.GetSelf(nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Retrieve logged user information ...")
	//retrive the stream of the current user
	stream := api.UserStream(nil)
	log.Println("User stream initialized ...")

	//listen to the stream
	for t := range stream.C {
		switch tw := t.(type) {
		//when we detect a new tweet in the stream
		case anaconda.Tweet:
			fmt.Printf("Tweet by: %s \n", tw.User.Name)
			// we verify is we are mentionned and if the bot should respond to that
			if isMentionned(user.ScreenName, tw) && isOnBlackList(tw.User, configuration) {
				log.Println("Mentionnned in tweet by blacklisted guy: Exterminate !")

				values := url.Values{}
				var response string
				// retrive the text and media to send to the person who dare speak to you
				for _, user := range configuration.Users {
					if user.Name == tw.User.ScreenName {
						mediaID := getTwitterMediaID(user.Response.MediaURL, api)
						values.Set("media_ids", mediaID)
						response = fmt.Sprintf("@%s %s", tw.User.ScreenName, user.Response.Text)
					}
				}

				values.Set("in_reply_to_status_id", tw.IdStr)
				log.Printf("My response to that is : %s", response)
				_, err := api.PostTweet(response, values)
				if err != nil {
					log.Fatalf("Sending tweet failed miserably : %s", err)
				}
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

func getTwitterMediaID(mediaURL string, api *anaconda.TwitterApi) string {
	log.Println("Uploading media ....")
	res, _ := http.Get(mediaURL)
	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("CgetTwitterMediaID(mediaURL string, api *anaconda.TwitterApi)annot retrive media : %s \n", err)
	}
	mediaBase64Str := base64.StdEncoding.EncodeToString(buf)
	media, err := api.UploadMedia(mediaBase64Str)
	if err != nil {
		log.Fatalf("Media upload failed : %s \n", err)
	}
	return media.MediaIDString
}
