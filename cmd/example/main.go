package main

import (
	"fmt"
	"log"

	dribbble "github.com/thedesque.com/de-dribbble"
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
	// 1 is the first page of paginated results, true to traverse all pages
	shots, err := c.Shots.GetShots(1, true)
	if err != nil {
		log.Fatal(err)
	}
	// checks if we have any shots
	shotCount := len(shots)
	if shotCount == 0 {
		log.Fatal("no shots found")
	}

	for i, shot := range shots {
		fmt.Printf("num %d of %d\n%s", i+1, shotCount, shot.String())
	}

	// ------------------------------------------------------------------------

	// get specific shot
	lastShot := shots[shotCount-1] // in terms of shot date, probably oldest from paginated results
	shot, err := c.Shots.GetShot(lastShot.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--SPECIFIC SHOT--\n%s", shot.String())

	// ------------------------------------------------------------------------

	// example toml output
	shotTomlString, err := shot.ToToml()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--TOML--\n%s", shotTomlString)

	// ------------------------------------------------------------------------

	// example yaml output
	shotYamlString, err := shot.ToYaml()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--YAML--\n%s", shotYamlString)
}
