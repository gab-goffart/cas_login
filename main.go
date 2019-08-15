package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	fmt.Println("===========================================")
	fmt.Println("Sending request to login route ...")
	res, err := http.Get("https://cas.usherbrooke.ca/login")
	fmt.Println("Response arrived !")
	fmt.Println("===========================================")

	if err != nil {
		fmt.Println("There was an error while trying to get to authentication service from UdeS")
		return
	}

	document, err := html.Parse(res.Body)
	defer res.Body.Close()

	if err != nil {
		fmt.Println("There was an error parsing the request's body")
		return
	}

	cookies := res.Cookies()
	form := url.Values{}
	form.Set("username", "gofg2301")
	form.Set("password", "Goffart2006")

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

	fmt.Println(form)
	fmt.Println("===========================================")
	fmt.Println("Sending request to authenticate")
	fmt.Println("===========================================")

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

	fmt.Println("response headers : ", res.Header)
	fmt.Println("response status : ", res.Status)
	fmt.Println("response set-cookie : ", res.Header.Get("Set-Cookie"))
	fmt.Println("response cookies : ", res.Cookies())
}
