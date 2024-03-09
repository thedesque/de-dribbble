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

	// ------------------------------------------------------------------------

	// get currently logged in user
	user, err := c.User.GetUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", user.String())

	// ------------------------------------------------------------------------

	// get all shots from logged in user
	shots, err := c.Shots.GetShots()
	if err != nil {
		log.Fatal(err)
	}

	maxShotsNum := len(*shots)
	for i, shot := range *shots {
		fmt.Printf("num %d of %d\n%s", i+1, maxShotsNum, shot.String())
	}

	// ------------------------------------------------------------------------

	// get specific shot
	shot, err := c.Shots.GetShot(23275914)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", shot.String())

	// ------------------------------------------------------------------------

	// example toml output
	shotTomlString, err := shot.ToToml()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s", shotTomlString)

	// ------------------------------------------------------------------------

	// example yaml output
	shotYamlString, err := shot.ToYaml()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s", shotYamlString)
}
