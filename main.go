package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
	"os"
	"strings"
)

func main() {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")
	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessSecret, consumerKey, consumerSecret)

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		log.Fatalln("VerifyCredentials", err)
	}
	//spew.Dump(user)
	fmt.Println("currently have", user.StatusesCount)
	//spew.Dump(user)

	//client.Statuses.Retweets()
	maxID := int64(1290974423576920065)
	tru := true

	var removedCount int
	for i := 0; i < 10; i++ {
		tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			UserID:          user.ID,
			Count:           30,
			MaxID:           maxID,
			IncludeRetweets: &tru,
		})
		if err != nil {
			log.Fatalln("UserTimeline", err)
		}
		for _, tweet := range tweets {
			//fmt.Println(i, tweet.ID, tweet.CreatedAt)
			fmt.Printf("-%q %d %q\n", tweet.CreatedAt, tweet.ID, tweet.Text)
		}

		fmt.Println(len(tweets), maxID)
		for _, tweet := range tweets {
			maxID = tweet.ID // Maybe ID - 1?
			if !strings.Contains(tweet.Source, "kyeett-twitterbot") {
				continue
			}
			removedCount++
			//fmt.Println(i, tweet.ID, tweet.CreatedAt)
			fmt.Printf("-%q %d %q\n", tweet.CreatedAt, tweet.ID, tweet.Text)
			//fmt.Printf("%q %q\n", tweet.CreatedAt, tweet.Text)
			//twt, resp, err := client.Statuses.Unretweet(tweet.RetweetedStatus.ID, &twitter.StatusUnretweetParams{
			//	ID:        tweet.RetweetedStatus.ID,
			//	TrimUser:  &tru,
			//	TweetMode: "compact",
			//})
			//if err != nil || resp.StatusCode != 200 {
			//	log.Fatalln("Unretweet", err, resp.StatusCode, resp.Status)
			//}
			twt, err := api.UnRetweet(tweet.RetweetedStatus.ID, true)
			if err != nil {
				log.Fatalln("unretweet", err)
			}
			fmt.Printf("x%q %d %q\n", twt.CreatedAt, twt.Id, twt.Text)
			//fmt.Println(res.Text)
		}
	}
	fmt.Println("Removed", removedCount)
	fmt.Println("Last ID was", maxID)

	//
	//res, err := api.UnRetweet(1067750435631652864, true)
	//if err != nil {
	//	log.Fatalln("unretweet", err)
	//}
	//spew.Dump(res)
}
