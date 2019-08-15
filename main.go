package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	username := flag.String("u", "", "-u <username>")
	password := flag.String("p", "", "-p <password>")

	flag.Parse()

	if *username == "" || *password == "" {
		fmt.Println("-u and -p are required")
		return
	}

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("Sending request to login route ...")
	fmt.Println("===========================================")

	res, err := http.Get("https://cas.usherbrooke.ca/login")
	if err != nil {
		fmt.Println("There was an error while trying to get to authentication service from UdeS")
		return
	}
	defer res.Body.Close()

	document, err := html.Parse(res.Body)

	if err != nil {
		fmt.Println("There was an error parsing the request's body")
		return
	}

	cookies := res.Cookies()
	form := url.Values{}
	form.Set("username", *username)
	form.Set("password", *password)
	form.Set("_eventId", "submit")
	form.Set("submit", "")

	var f func(*html.Node)

	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "input" {
			for _, a := range node.Attr {
				if a.Key == "name" {
					if a.Val == "lt" {
						form.Set("lt", node.Attr[2].Val)
					}
					if a.Val == "execution" {
						form.Set("execution", node.Attr[2].Val)
					}
				}
			}
		}
		for next := node.FirstChild; next != nil; next = next.NextSibling {
			f(next)
		}
	}

	f(document)

	req, err := http.NewRequest("POST", "https://cas.usherbrooke.ca/login", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println("There was an error building the request")
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("Sending request to authenticate")
	fmt.Println("===========================================")
	fmt.Println()

	client := http.Client{}

	res, err = client.Do(req)

	if err != nil {
		fmt.Println("Error trying to authenticate")
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if err != nil {
		fmt.Println("Could not read body of response")
		return
	}

	fmt.Println("response status : ", res.Status)
	for i, cookie := range res.Cookies() {
		fmt.Println("\t Cookie ", i, " : ", cookie)
	}
}
