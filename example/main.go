package main

import (
	"fmt"

	dribbble "github.com/thedesque.com/desqfolio"
)

func main() {
	token := dribbble.Auth()

	cfg := dribbble.NewConfig(token)
	d := dribbble.New(cfg)

	// get currently logged in user
	user, _ := d.User.GetUser()
	fmt.Printf("%v", user)

	// get all shots from logged in user
	shots, _ := d.Shots.GetShots()
	fmt.Printf("%v", shots)
}
