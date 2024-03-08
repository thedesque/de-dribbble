package main

import (
	"fmt"
	"log"

	dribbble "github.com/thedesque.com/desqfolio"
)

func main() {
	// holds our default flags and a nil pointer to token
	cfg := dribbble.NewConfig()
	// starts the automated auth flow and sets token to config
	if err := dribbble.OauthStart(cfg); err != nil {
		log.Fatal(err)
	}
	// create new API client with config
	c := dribbble.NewClient(cfg)

	// get currently logged in user
	user, err := c.User.GetUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", user.String())

	// get all shots from logged in user
	shots, err := c.Shots.GetShots()
	if err != nil {
		log.Fatal(err)
	}

	maxNum := len(*shots)
	for i, shot := range *shots {
		fmt.Printf("num %d of %d\n%s", i+1, maxNum, shot.String())
	}
}
