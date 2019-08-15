package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	cookies := res.Cookies()
	form := url.Values{}
	form.Set("username", "gofg2301")
	form.Set("password", "Goffart2006")

	req, err := http.NewRequest("POST", "https://cas.usherbrooke.ca/login", strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println("There was an error building the request")
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	for _, cookie := range cookies {
		req.AddCookie(cookie)
		fmt.Println("Cookie added : ", cookie)
	}

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

	fmt.Println("response headers: ", res.Header)
	fmt.Println("response status", res.Status)
	fmt.Println("response cookies : ", res.Cookies())
}
