package main

import (
	"fmt"
	"log"

	dribbble "github.com/thedesque.com/desqfolio"
)

func main() {
	cfg := dribbble.NewConfig()
	if err := dribbble.OauthStart(cfg); err != nil {
		log.Fatal(err)
	}

	c := dribbble.NewClient(cfg)

	// get currently logged in user
	user, err := c.User.GetUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", user)

	// get all shots from logged in user
	// shots, _ := client.Shots.GetShots()
	// fmt.Printf("%v", shots)
}
