/**
* Ask user for search keywords
* Ask for duration
* Ask user for tweet include a template to insert usernames
* Ask user if they want to include multiple users in single tweet
* TODO:
* 	1. [x] Tweet length
* 	2. [ ] Tweet velocity
* 					[x] There is a limit on GET requests
* 							So we can run a subroutine will fill the temoUserList on regular interval
* 					[ ] Limit can not be determined for POST request in realtime :(
								Expriment to get the error value when limit is crossed and when it is reset
								Once the limit is crossed wait till it is reset
* 	3. [x] Exit condition
* 					[x] Duration
* 					[x] Number of tweets
* 	4.  [x] Solve concurrent access issue
* 	5.  [x] Subroutine for fetching latest users
* 	6.  [x] Subroutine to send tweets
* 	7.  [x] Note down all the output into a log file so that it can be analyzed later
* 	8.  [+] Add some random content while replyinh to tweet
* 	9.  [-] Add suport for Tor Proxy
* 	10. [x] Handle Non 200 error while replying
* 	11. [-] Update log filename as {name}-{date}.log
* 	12. [-] move log files inside log folder
* 	13. [ ] Add option to filter negative keywords
* 	14. [-] Multiple sender workers
* 	15. [ ] Reply to new tweets first
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	//"io"
	"bufio"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

var credPath = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")

func readCredentials() error {
	b, err = ioutil.ReadFile(*credPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &oauthClient.Credentials)
}

var userList XUserList
var xReplyStatuses XReplyStatuses
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	userList.Init()
	xReplyStatuses.Init()
	// Inital Logging
	InitLogging(true, true, true, true, false)
	// Trace.Println("Tracing 123")
	// Info.Println("Info 123")
	// Warning.Println("Warning 123")
	// Error.Println("Error 123")
	// Debug.Println("Debug 123")
	// os.Exit(0)

	var endCriteriaValue int = 0
	var tweetText, keywordSearch, endCriteria, direction string = "", "", "", ""

	r := bufio.NewReader(os.Stdin)
	// var resp http.Response
	if err := readCredentials(); err != nil {
		log.Fatal(err)
	}

	// tokenCred := &oauth.Credentials{Token: "2846849851-UNwMEPigXogDrdMAPfvsxxDsC8nY0wdzOHB8xVi", Secret: "YSR6OUbYqBkAPCwVq5TOH30YByd6TSniqERuUv8Ftp2sT"}
	tempCred, err := oauthClient.RequestTemporaryCredentials(http.DefaultClient, "oob", nil)
	if err != nil {
		log.Fatal("RequestTemporaryCredentials:", err)
	}
	u := oauthClient.AuthorizationURL(tempCred, nil)
	fmt.Printf("1. Go to %s\n2. Authorize the application\n3. Enter verification code:\n", u)
	var code string
	fmt.Scanln(&code)
	tokenCred, _, err := oauthClient.RequestToken(http.DefaultClient, tempCred, code)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(tokenCred)

	//contents, _ := ioutil.ReadAll()

	// formTweet := url.Values{"status": {"You are simply amazing buddy"}}
	// resp, err := oauthClient.Post(http.DefaultClient, tokenCred,
	// 	"https://api.twitter.com/1.1/statuses/update.json", formTweet)
	// defer resp.Body.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	contents, _ := ioutil.ReadAll(resp.Body)
	// 	//fmt.Printf("%s\n", contents)
	// 	Debug.Printf("%s", contents)
	// 	Error.Printf("%s", contents)
	// }
	// os.Exit(0)

	fmt.Println(">> Enter search keywords: ")
	keywordSearch, _ = r.ReadString('\n')

	for !(endCriteria == "D" || endCriteria == "d" || endCriteria == "T" || endCriteria == "t") {
		fmt.Println(">> End criteria ? (D: Duration / T: Number of tweets) ")
		fmt.Scanln(&endCriteria)
		endCriteria = strings.ToLower(endCriteria)
	}

	if endCriteria == "d" {
		fmt.Println(">> Duration value in minutes: ")
	} else if endCriteria == "t" {
		fmt.Println(">> Number of tweets to reply: ")
	}
	fmt.Scanln(&endCriteriaValue)

	fmt.Println(">> Enter tweet: ")
	fmt.Println("Ex: Hey [user], check this awesome sketch http://bitly/xyz")
	fmt.Println("1. [user] will be replaced by @username 2. Dont add important stuff like a link in the end, if username is long it will be truncated. 3. Keep some sapce for adding random #hashtag at the end to prevent getting blocked due to similar content")
	tweetText, _ = r.ReadString('\n')
	for len(tweetText) < 10 || len(tweetText) > 140 || (strings.Contains(tweetText, "[user]") && len(tweetText) >= 130) || !strings.Contains(tweetText, "[user]") {
		if !strings.Contains(tweetText, "[user]") {
			fmt.Println("[user] must be a part of the tweet, Please try again")
		}
		if len(tweetText) < 10 {
			fmt.Println("Tweet too small, Please try again")
		}
		if len(tweetText) > 140 {
			fmt.Println("Tweet too large, You entered", len(tweetText), "/140 characters. Please try again")
		}
		if strings.Contains(tweetText, "[user]") && len(tweetText) >= 130 {
			fmt.Println("You must leave some character for including username, your current tweet length is ", len(tweetText), "/140")
		}
		fmt.Println(">> Enter tweet: ")
		tweetText, _ = r.ReadString('\n')
	}

	for !(direction == "o" || direction == "O" || direction == "n" || direction == "N" || direction == "b" || direction == "B") {
		fmt.Println(">> Enter direction: (O: Old Tweets / N: New Tweets / B: Both alternatively)")
		fmt.Scanln(&direction)
		direction = strings.ToLower(direction)
	}

	if endCriteria == "d" {
		go func() {
			time.Sleep(time.Duration(endCriteriaValue) * time.Minute)
			os.Exit(0)
		}()
	}

	// Run Goroutines
	go searchTweets(keywordSearch, direction, tokenCred)
	go statusUpdate(tweetText, endCriteria, endCriteriaValue, tokenCred)

	// User can terminate by inputting "end" in console
	endNow := "n"
	for endNow != "end" {
		fmt.Println(">> End Now? (end)")
		fmt.Scanln(&endNow)
		if endNow == "end" {
			// TODO: Maybe dump some files
		}
	}
}

func searchTweets(keywordSearch string, direction string, tokenCred *oauth.Credentials) {
	k := 0
	var ramainingRequests, resetTimeStamp int64 = 0, 0
	var maxId, minId int64 = 0, 0

	for {
		form := url.Values{}
		// Find tweets. It returns only 100 whatever be the count. So sad :( Fuck you twitter
		if k == 0 {
			form = url.Values{"q": {keywordSearch}, "count": {"2"}, "result_type": {"recent"}}
			//Debug.Println("No min No max")
		}
		if direction == "o" {
			form = url.Values{"q": {keywordSearch}, "count": {"2"}, "result_type": {"recent"}, "max_id": {strconv.FormatInt(minId-1, 10)}}
			//Debug.Println("OLD: MinId = ", minId)
		}
		if direction == "n" {
			form = url.Values{"q": {keywordSearch}, "count": {"2"}, "result_type": {"recent"}, "since_id": {strconv.FormatInt(maxId+1, 10)}}
			//Debug.Println("NEW: MaxId = ", maxId)
		}

		if direction == "b" {
			if k%2 == 0 {
				form = url.Values{"q": {keywordSearch}, "count": {"2"}, "result_type": {"recent"}, "max_id": {strconv.FormatInt(minId-1, 10)}}
				//Debug.Println("BOTH: MinId = ", minId)
			} else {
				form = url.Values{"q": {keywordSearch}, "count": {"2"}, "result_type": {"recent"}, "since_id": {strconv.FormatInt(maxId+1, 10)}}
				//Debug.Println("BOTH: MaxId = ", maxId)
			}
		}
		//fmt.Println(form)
		resp, err := oauthClient.Get(http.DefaultClient, tokenCred,
			"https://api.twitter.com/1.1/search/tweets.json", form)
		defer resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		ramainingRequests, _ = strconv.ParseInt(resp.Header["X-Rate-Limit-Remaining"][0], 10, 64)
		//allowedRequests, _ = strconv.ParseInt(resp.Header["X-Rate-Limit-Limit"][0], 10, 64)
		resetTimeStamp, _ = strconv.ParseInt(resp.Header["X-Rate-Limit-Reset"][0], 10, 64) // converted to miliseconds
		resetTimeStamp *= 1000

		var srobj searchResponse
		searchResponseBody, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(searchResponseBody, &srobj)

		for i := range srobj.Statuses {
			//Debug.Println(srobj.Statuses[i].Id)

			if (strings.Contains(strings.ToLower(srobj.Statuses[i].Text), "trend")) != true {
				userList.Set(srobj.Statuses[i].User.Id, srobj.Statuses[i].User)
				// if _, ok := replyStatuses[srobj.Statuses[i].User.Id]; !ok {
				// 	replyStatuses[srobj.Statuses[i].User.Id] = ReplyStatus{Replied: false}
				// }
				if !xReplyStatuses.IsSet(srobj.Statuses[i].User.Id) {
					xReplyStatuses.Initiate(srobj.Statuses[i].User.Id)
				}
			}
			if minId == 0 {
				minId = srobj.Statuses[i].Id
			} else if minId > srobj.Statuses[i].Id {
				minId = srobj.Statuses[i].Id
			}

			if maxId == 0 {
				maxId = srobj.Statuses[i].Id
			} else if maxId < srobj.Statuses[i].Id {
				maxId = srobj.Statuses[i].Id
			}
		}

		if resp.StatusCode != 200 {
			Debug.Printf("%s %s %s", searchResponseBody, resp.Header, resp.Status)
			Error.Printf("%s %s %s", searchResponseBody, resp.Header, resp.Status)
		} else {
			Info.Printf("%s %s %s", searchResponseBody, resp.Header, resp.Status)
			if len(srobj.Statuses) > 0 {
				if direction == "b" {
					if k%2 == 0 {
						Debug.Printf("%d old tweets found", len(srobj.Statuses))
					} else {
						Debug.Printf("%d new tweets found", len(srobj.Statuses))
					}
				}
				if direction == "o" {
					Debug.Printf("%d old tweets found", len(srobj.Statuses))
				}
				if direction == "n" {
					Debug.Printf("%d new tweets found", len(srobj.Statuses))
				}
			}
		}

		// Insert calculated delay (Useful when there will be multiple senders)
		var delay int64
		if ramainingRequests != 0 {
			delay = (resetTimeStamp - time.Now().UnixNano()/int64(time.Millisecond)) / ramainingRequests
		} else {
			delay = (resetTimeStamp - time.Now().UnixNano()/int64(time.Millisecond))
		}
		//fmt.Printf("Pass ends sleeping for %d seconds", delay/1000)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		// Because there is a limit to how much tweet a user can send
		// time.Sleep(5 * time.Millisecond)
		k++
	}
}

func statusUpdate(tweetText, endCriteria string, endCriteriaValue int, tokenCred *oauth.Credentials) {
	time.Sleep(5 * time.Second) // Wait for some time till userList is populated
	k := 0
	totalPeopleReplied := 0
	for {
		// Fill tempUserList outside
		// While adding stuff to tempUserList check if replyStatuses has a value
		useridList := xReplyStatuses.ListUseridUnsent()
		peopleReplied := 0

		for len(useridList) == 0 {
			time.Sleep(1 * time.Second)
		}
		for j := range useridList {
			totalPeopleReplied++

			currentUserId := useridList[j]
			currentUser := userList.Get(currentUserId)
			processedTweetText := strings.Replace(tweetText, "[user]", "@"+currentUser.ScreenName, -1) + " #" + randSeq(3)

			formTweet := url.Values{"status": {processedTweetText}}
			resp, err := oauthClient.Post(http.DefaultClient, tokenCred,
				"https://api.twitter.com/1.1/statuses/update.json", formTweet)
			searchResponseBody, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				log.Fatal(err)
				// } else {
				// 	contents, _ := ioutil.ReadAll(resp.Body)
				// 	fmt.Printf("%s\n", contents)
			}
			xReplyStatuses.Sent(currentUserId)

			if resp.StatusCode != 200 {
				Debug.Printf("%s %s %s", searchResponseBody, resp.Header, resp.Status)
				Error.Printf("%s %s %s", searchResponseBody, resp.Header, resp.Status)
			} else {
				Info.Printf("%s %s %s", searchResponseBody, resp.Header, resp.Status)
				//Debug.Printf("%s", processedTweetText)
			}

			if resp.StatusCode == 403 {
				var errobj TwitterResponseError
				_ = json.Unmarshal(searchResponseBody, &errobj)
				if errobj.errors[0].Code == 185 {
					// TODO: Daily user limit reached
					// Stop for half an hour
					Debug.Printf("Tweet limit reached, Sleeping for half an hour")
					time.Sleep(30 * time.Minute)
				}
				if errobj.errors[0].Code == 226 {
					// Inform user to create another campaign and exit
					Debug.Printf("Seems like twitter has detected your campaign as spam try another campign, may be another user too")
					os.Exit(0)
				}
			}
			if endCriteria == "t" && totalPeopleReplied >= endCriteriaValue {
				fmt.Println("Exiting because end criteria reached. Number of tweets replied : ", totalPeopleReplied)
				os.Exit(0)
			}
			peopleReplied++

			// User can only send 50 request per half-an-hour (Post every 36 seconds)
			time.Sleep(time.Duration(rand.Intn(5)+3) * time.Second) // Wait for 3 to 7 seconds
		}
		k++

		Debug.Printf("%d people replied", peopleReplied)
	}
}
