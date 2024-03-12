# De-Dribbble
A Go library for interacting with the Dribbble API v2 (WIP).
Please refer to official [Dribbble API v2 Docs](http://developer.dribbble.com/v2/) for more information about API itself.
 
## Install
```
go get github.com/thedesque/de-dribbble
```

## Usage
```go
import "github.com/thedesque/de-dribbble"

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

fmt.Printf("%s", user.String()) // or .ToToml or .ToYaml
```

Inspired by [go-dribbble](https://github.com/bedakb/go-dribbble).
