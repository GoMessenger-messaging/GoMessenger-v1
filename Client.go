package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var urlStr, uname, channel, message string
var reader = bufio.NewReader(os.Stdin)

func checkUname(uname string, urlStr string) bool {
	urlStr = urlStr + "register?name=" + url.QueryEscape(uname)
	resp, err := http.Get(urlStr)
	if err != nil {
		ret := " - Error: " + err.Error()
		fmt.Print(ret, "\n\n")
		return false
	}
	defer resp.Body.Close()

	ret, _ := io.ReadAll(resp.Body)
	fmt.Print("\n - ", string(ret), "\n\n")

	if string(ret) == "Registered successfully" {
		return true
	} else {
		return false
	}
}
func getMessages() {
	resp, err := http.Get(urlStr + "get_messages?channel=" + url.QueryEscape(channel))
	if err != nil {
		fmt.Print("\n - Error: ", err, "\n\n")
		return
	}
	defer resp.Body.Close()

	ret, _ := io.ReadAll(resp.Body)
	fmt.Print("\n", string(ret), "\n\n")
}
func ping() {
	for true {
		http.Get(urlStr + "ping?name=" + url.QueryEscape(uname))
		time.Sleep(5 * time.Minute)
	}
}

func Server() {
	for true {
		fmt.Print("|Enter the server's url (:h for more info): ")
		urlStr, _ = reader.ReadString('\n')
		urlStr = strings.TrimRight(urlStr, "\r\n")
		if urlStr == ":h" {
			fmt.Print("\n - You may enter the url in the format 'http://{ip}:{port}/' (port only if needed)\n - As of 14/3/2023 (d/m/y) 'https://gomessenger.link/' is the official server\n - :d for the official server\n\n")
		} else {
			if urlStr == ":d" {
				urlStr = "https://gomessenger.link/"
			}
			_, err := http.Get(urlStr)
			if err != nil {
				fmt.Print("\n - Error: ", err, "\n\n")
			} else {
				fmt.Print("\n - Connected to ", urlStr, "\n\n")
				if urlStr[len(urlStr)-1:] != "/" {
					urlStr = urlStr + "/"
				}
				return
			}
		}
	}
}
func Username() {
	for true {
		fmt.Print("|Enter your username (:h for more info): ")
		uname, _ = reader.ReadString('\n')
		uname = strings.TrimRight(uname, "\r\n")
		if uname == ":h" {
			fmt.Print("\n - Your username is used to identify you to other users\n - You may change your username at any time\n - You may not use offensive language in your username\n - You may chose a username that is still available on the server\n\n")
		} else {
			if checkUname(uname, urlStr) {
				return
			}
		}
	}
}
func Channel() {
	for true {
		fmt.Print("|Enter a channel (:h for more info): ")
		channel, _ = reader.ReadString('\n')
		channel = strings.TrimRight(channel, "\r\n")
		if channel == ":h" {
			fmt.Print("\n - You may join any channel that exists on the server\n - You may create a new channel by entering a name that is not taken\n\n")
		} else {
			fmt.Print("\n - Joined channel ", channel, "\n - Press enter to get all messages\n\n")
			return
		}
	}
}

func main() {
	fmt.Print(" --- Welcome to Jeroen's Messenger - v1.1.0 --- \n\n")
	Server()
	Username()
	Channel()

	go ping()

	for true {
		fmt.Print("|Enter a message (:h for more info): ")
		message, _ = reader.ReadString('\n')
		message = strings.TrimRight(message, "\r\n")
		if len(message) > 0 && message[0] == ':' {
			if message == ":h" {
				fmt.Print("\n - You may enter any message you like\n - Press enter without writing anything to get all messages in the channel\n - ':h' to get help\n - ':cs' to get the current server\n - ':cu' to get the current username\n - ':cc' to get the current channel\n - ':s' to change the server\n - ':u' to change your username\n - ':c' to change the channel\n - ':q' to quit\n\n")
			} else if message == ":cs" {
				fmt.Print("\n - Current server: ", urlStr, "\n\n")
			} else if message == ":cu" {
				fmt.Print("\n - Current username: ", uname, "\n\n")
			} else if message == ":cc" {
				fmt.Print("\n - Current channel: ", channel, "\n\n")
			} else if message == ":s" {
				Server()
				Username()
				Channel()
			} else if message == ":u" {
				Username()
			} else if message == ":c" {
				Channel()
			} else if message == ":q" {
				return
			} else {
				fmt.Print("\n - Error: Unknown command\n\n")
			}
		} else {
			if message != "" {
				_, err := http.Get(urlStr + "send?name=" + url.QueryEscape(uname) + "&channel=" + url.QueryEscape(channel) + "&message=" + url.QueryEscape(message))
				if err != nil {
					fmt.Print("\n - Error: ", err, "\n\n")
				}
			}
			getMessages()
		}
	}
}
