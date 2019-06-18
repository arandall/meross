/*
Package meross-cloud contains a CLI tool to get device credentials from Meross to enable controlling of existing
devices.

Creating an account

	meross-cloud -create -email <email> -password <password> -region US

Viewing existing credentials and devices

	meross-cloud -email <email> -password <password>
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	args := os.Args
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	baseURL := fs.String("url", "https://iot.meross.com", "base URL of meross cloud service")
	createFlag := fs.Bool("create", false, "create account instead of logging in.")
	emailFlag := fs.String("email", "", "email address to sign in as")
	passwordFlag := fs.String("password", "", "password")
	regionFlag := fs.String("region", "", "2 letter country code for account (eg. US)")
	timeoutFlag := fs.Duration("timeout", 10*time.Second, "http timeout in seconds for requests (default: 10s)")

	fs.Parse(args[1:])

	url, err := url.Parse(*baseURL)
	if err != nil {
		fs.Usage()
		exit(fmt.Errorf("url: %v", err))
	}

	c := NewClient(&http.Client{
		Timeout: *timeoutFlag,
	}, url)

	var acc *Account
	if *createFlag {
		u := &User{
			Email:    *emailFlag,
			Password: *passwordFlag,
			Region:   *regionFlag,
		}
		acc, err = c.Signup(u)
		if err != nil {
			exit(fmt.Errorf("create: %v", err))
		}
	} else {
		u := &User{
			Email:    *emailFlag,
			Password: *passwordFlag,
		}
		acc, err = c.Login(u)
		if err != nil {
			exit(fmt.Errorf("login: %v", err))
		}
	}

	printAccountDetails(acc)

	c.ApplyCredentials(acc)
	devices, err := c.DeviceList()
	if err != nil {
		exit(fmt.Errorf("list: %v", err))
	}

	fmt.Println("Devices:")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(devices); err != nil {
		exit(fmt.Errorf("encode: %v", err))
	}
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func printAccountDetails(a *Account) {
	fmt.Println("Device Credentials:")
	fmt.Printf("\t User ID: %s\n", a.UserID)
	fmt.Printf("\t Key:     %s\n", a.Key)
}
