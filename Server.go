package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// These are the server's settings
var port string = ":8080"
var timezone string = "Europe/Zurich"
var blacklist = [...]string{
	"substring1",
	"substring2",
	"substring3",
	// ...
}

// These are other global variables
var users = [128]string{}
var lastPinged = [128]time.Time{}
var channels = map[string]string{}

func Register(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if len(name) < 1 {
		w.Write([]byte("No name provided"))
		return
	}
	for _, v := range blacklist {
		if strings.Contains(name, v) {
			w.Write([]byte("Name contains offensive language"))
			return
		}
	}
	for _, v := range users {
		if name == v {
			w.Write([]byte("Name already taken"))
			return
		}
	}
	for i, v := range users {
		if v == "" {
			users[i] = name
			w.Write([]byte("Registered successfully"))
			return
		}
	}
	w.Write([]byte("Server full"))
}
func Send(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	channel := query.Get("channel")
	message := query.Get("message")

	loc, _ := time.LoadLocation(timezone)
	currTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	channels[channel] += currTime + " | " + name + " > " + message + "\n"
}
func GetMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	channel := query.Get("channel")
	fmt.Fprint(w, channels[channel])
}
func UsersList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, users)
}
func Ping(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	for i, v := range users {
		if v == name {
			lastPinged[i] = time.Now()
		}
	}
}

func remUsers() {
	for true {
		for i, v := range lastPinged {
			if time.Since(v) > time.Minute*10 {
				users[i] = ""
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

func main() {
	go remUsers()

	http.HandleFunc("/register", Register)
	http.HandleFunc("/send", Send)
	http.HandleFunc("/get_messages", GetMessages)
	http.HandleFunc("/users", UsersList)
	http.HandleFunc("/ping", Ping)
	http.ListenAndServe(port, nil)
}
