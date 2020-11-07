package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	token        string
	userID       string
	messageCount int = 0
	failureCount int = 0
	client           = &http.Client{}
	scan             = bufio.NewScanner(os.Stdin)
)

type sendbody1 struct {
	Recipients []interface{} `json:"recipients,omitempty"`
}

type sendbody2 struct {
	Content string `json:"content"`
}

func main() {
	fmt.Print("Enter your token:")
	fmt.Scan(&token)
	fmt.Print("Enter the UserID of the person you want to message:")
	fmt.Scan(&userID)

	var user sendbody1
	user.Recipients = []interface{}{userID}
	jsonres, _ := json.Marshal(user)
	code, response := sendrequest("POST", "https://discord.com/api/v8/users/@me/channels", bytes.NewReader(jsonres))

	switch code {
	case 200:
		var decoded map[string]string
		json.Unmarshal([]byte(response), &decoded)
		fmt.Print("Enter your message: ")
		for {
			scan.Scan()
			text := scan.Text()
			if len(text) > 0 {
				sendMessage(decoded["id"], text)
				fmt.Print(fmt.Sprintf("Enter your next message[Msg(s):%d,fail(s):%d]: ", messageCount, failureCount))
			}
		}
	case 400:
		exit("wrong UserID")
	case 401:
		exit("wrong token")
	default:
		exit("wtf happened????")
	}
}

func sendMessage(channelID string, message string) {
	sendmessage := &sendbody2{
		Content: message,
	}
	jsonres, _ := json.Marshal(sendmessage)
	code, _ := sendrequest("POST", fmt.Sprintf("https://discord.com/api/v8/channels/%s/messages", channelID), bytes.NewReader(jsonres))
	switch code {
	case 200:
		messageCount++
	case 403:
		failureCount++
		fmt.Println("Couldn't DM that user (no mutual servers?)")
	case 429:
		failureCount++
		fmt.Println("Hit a rate limit try resending that last message in a few seconds.")
	default:
		failureCount++
		fmt.Println("something odd happened")
	}
}

func sendrequest(method string, url string, body *bytes.Reader) (int, []byte) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.130 Safari/537.36")
	res, _ := client.Do(req)
	resb, _ := ioutil.ReadAll(res.Body)
	return res.StatusCode, resb
}

func exit(message string) {
	fmt.Println(message)
	time.Sleep(time.Duration(5) * time.Second)
	os.Exit(0)
}
