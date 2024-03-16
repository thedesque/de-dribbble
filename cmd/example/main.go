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
	// returns a wrapped slice of shots and an error
	// 1 is the first page of paginated results, true to traverse all pages
	shots, err := c.Shots.GetShots(1, true)
	if err != nil {
		log.Fatal(err)
	}
	// access underlying slice of shots
	// checks if we have any shots
	shotCount := len(shots.Slice)
	if shotCount == 0 {
		log.Fatal("no shots found")
	}
	// range the slice of single ShotOut struct to do something
	for i, shot := range shots.Slice {
		fmt.Printf("num %d of %d\n%s", i+1, shotCount, shot.String())
	}

	// ------------------------------------------------------------------------

	// get specific shot via API
	// could use GetShots and search slice result to reduce API calls
	lastShot := shots.Slice[shotCount-1] // in terms of shot date, probably oldest from paginated results
	shot, err := c.Shots.GetShot(lastShot.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--SPECIFIC SHOT--\n%s", shot.String())

	// ------------------------------------------------------------------------

	// example shot to toml stdout
	tomlShot, err := dribbble.ToToml(shot)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--TOML--\n%s", tomlShot)

	// example shots to toml to os file
	tomlShots, err := dribbble.ToToml(shots)
	if err != nil {
		log.Fatal(err)
	}

	err = dribbble.ToFile("shots.toml", tomlShots)
	if err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------------------

	// example shot to yaml stdout
	yamlShot, err := dribbble.ToYaml(shot)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n--YAML--\n%s", yamlShot)

	// example shots to yaml to os file
	yamlShots, err := dribbble.ToYaml(shots)
	if err != nil {
		log.Fatal(err)
	}

	err = dribbble.ToFile("shots.yaml", yamlShots)
	if err != nil {
		log.Fatal(err)
	}
}
