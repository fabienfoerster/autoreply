# twitter-autoreply

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

Auto reply to some tweet with predetermined response

##Description
It super easy, if you want to use it, just fork it, and change some values in the config.json file.
```json
{
  "users" : [
    {
      "name" : "clement0210",
      "response" : {
        "text" : "",
        "mediaURL" : "http://media2.giphy.com/media/ZsmCoNbVPy0QE/giphy.gif"
      }
    }
  ]
}
```

You specify in this the screenname of the person you want to respond to automatically (the @name), the text and an optional media ( photo or gif).

You click on the deploy to heroku button.

And finally you specify some environment variable for the Twitter API.

##Environment variables
The app use the twitter API so you'll need to set 5 environment variables
```
heroku config:set TWITTER_API_KEY=XXX
heroku config:set TWITTER_API_SECRET=XXX
heroku config:set TWITTER_ACCESS_TOKEN=XXX
heroku config:set TWITTER_ACCESS_TOKEN_SECRET=XXX
```

And you are good to go <3 ( just turn on the worker in the heroku dashboard ;)
