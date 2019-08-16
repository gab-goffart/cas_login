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

	var cookies []*http.Cookie

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("Sending request to login route")
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

	cookies = append(cookies, res.Cookies()...)
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
	res.Body.Close()

	if err != nil {
		fmt.Println("Could not read body of response")
		return
	}

	for _, cookie := range res.Cookies() {
		isPresent := false
		for _, oldCookie := range cookies {
			if oldCookie.String() == cookie.String() {
				isPresent = true
			}
		}

		if isPresent == false {
			cookies = append(cookies, cookie)
		}
	}

	fmt.Println("===========================================")
	fmt.Println("Trying to get results for session H19")
	fmt.Println("===========================================")
	fmt.Println()

	req, err = http.NewRequest("GET", "https://cas.usherbrooke.ca/login?service=https%3A%2F%2Fwww.gel.usherbrooke.ca%2Fgrille-notes%2Fapi%2Fgrid%2Fresults%3Ftrimester%3DH19", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	fmt.Println(req.Header)

	res, err = client.Do(req)
	if err != nil || res.StatusCode != 302 {
		fmt.Println("Error trying to GET for grades authentication", err)
		return
	}
	defer res.Body.Close()

	location, err := res.Location()
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err = http.NewRequest("GET", location.String(), nil)
	if err != nil {
		fmt.Println("Error trying to create a request")
		return
	}

	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Error fetching for the link")
		return
	}
	res.Body.Close()

	cookies = res.Cookies()

	req, err = http.NewRequest("GET", "https://gel.usherbrooke.ca/grille-notes/api/grid/results?trimester=H19", nil)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Error fetching for the grades")
		return
	}

	var body []byte
	nbBytes, err := res.Body.Read(body)
	if err != nil {
		fmt.Println("Error parsing the body")
		return
	}

	if nbBytes == 0 {
		fmt.Println("body is of size 0, authentication probably failed ...")
		return
	}

	fmt.Println("Grades : ", string(body))

}
