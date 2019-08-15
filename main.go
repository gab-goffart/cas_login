package main

import (
	"bytes"
	"fmt"
	"net/http"
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

	var cookieValues map[string]string
	for _, cookie := range cookies {
		cookieValues[cookie.Name] = cookie.Value
	}

	var form map[string]string
	form["username"] = "gab.goffart@gmail.com"
	form["password"] = "Goffart2006"

	formBytes := []byte(form)

	req, err := http.NewRequest("POST", "cas.usherbrooke.ca/login", bytes.NewBuffer(form))
	if err != nil {
		fmt.Println("There was an error while trying to get to authentication service from UdeS")
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookies", string(cookieValues))

	fmt.Println("===========================================")
	fmt.Println("Here is the response the server sent: ")
	fmt.Println("Headers : ", res.Header)
	fmt.Println("Cookies : ", res.Cookies())
	fmt.Println("===========================================")
}
