package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var SECRET string = "Hello"
var TELEGRAM_KEY string = "mock"
var TELEGRAM_CHAT string = ""
var WEB_HOST string = "http://localhost:8000"
var ASK_EVERY time.Duration = time.Duration(3) * time.Second
var REVEAL_MESSAGE_SECRET_AFTER uint = 4

var send_in uint = 0
var mutex = &sync.Mutex{}

func resetTimer() {
	mutex.Lock()
	send_in = 0
	mutex.Unlock()
}

func sendMessage(msg string) {
	if TELEGRAM_KEY == "mock" {
		fmt.Println("========")
		fmt.Println(msg)
		fmt.Println("========")
		return
	}

	payload := map[string]string{
		"chat_id": TELEGRAM_CHAT,
		"text":    msg,
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_KEY), "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Error sending message:", string(body))
	}
}

func messageThread() {
	resetTimer()
	for {
		time.Sleep(ASK_EVERY)
		mutex.Lock()
		send_in = send_in + 1
		if send_in >= REVEAL_MESSAGE_SECRET_AFTER {
			sendMessage(SECRET)
		} else {
			sendMessage(fmt.Sprintf("Is valhalla still waiting?\n Click here\n%s/valhalla_awaits\n to reset timer\nIteration: %v", WEB_HOST, send_in))
		}
		mutex.Unlock()
	}
}

func main() {
	sendMessage(fmt.Sprintf("Started Valhalla Bot..\n * Asks every %s\n * If no response within %v iterations I sends secret", ASK_EVERY, REVEAL_MESSAGE_SECRET_AFTER))
	go messageThread()

	http.HandleFunc("/general/ok", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	http.HandleFunc("/valhalla_awaits", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, `
                <form method="post" action="/valhalla_awaits">
                    <button type="submit">Is Valhalla still waiting?</button>
                </form>
            `)
		case http.MethodPost:
			resetTimer()
			sendMessage("Timer reset. Valhalla Awaits..")
			fmt.Fprint(w, "ok")
		}
	})

	http.ListenAndServe(":8000", nil)
	sendMessage("Shutting Down Server")
}
