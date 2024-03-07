package dribbble

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cli/browser"
	"github.com/fatih/color"
	"golang.org/x/oauth2"
)

// Dribbble scopes
const (
	Public = "public"
	Upload = "upload"
)

// Default token file
const tokenFile = "token.json"

// Endpoint is Dribbbles's OAuth 2.0 default endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://dribbble.com/oauth/authorize",
	TokenURL: "https://dribbble.com/oauth/token",
}

// newOauthConf returns a new oauth2.Config with values from environment variables.
func newOauthConf() *oauth2.Config {
	scope := []string{Public}
	if os.Getenv("WRITE_SCOPE") == "true" {
		scope = append(scope, Upload)
	}

	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       scope,
		Endpoint:     Endpoint,
	}
}

// OauthStart returns a token from a file if it exists, otherwise it starts a server to authenticate the user.
func OauthStart(config *Config, filePath ...string) error {
	// validate number of arguments
	if len(filePath) > 1 {
		return fmt.Errorf("expected at most one argument; got %d", len(filePath))
	}

	// determine file path
	var path string
	if len(filePath) == 0 {
		path = tokenFile // use default token file if no file path provided
	} else {
		path = filePath[0]
	}

	var token *oauth2.Token
	var err error
	// get token from file
	token, err = OauthFile(path)
	if err != nil && strings.Contains(err.Error(), "token file does not exist") {
		// no token file so start server to authenticate user
		token, err = OauthServer()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// add token to config
	config.Token = token

	// write token to file
	if config.Flags.TokenFile {
		if err := writeTokenToFile(token); err != nil {
			return fmt.Errorf("failed to write token to file: %v", err)
		}
	}

	return nil
}

// OauthFile returns a token from a file if it exists. If no file path is
// provided or if the file does not exist, it returns nil.
func OauthFile(filePath string) (*oauth2.Token, error) {
	// checks if file path is empty
	if filePath == "" {
		return nil, fmt.Errorf("expected a file path; got %s", filePath)
	}

	// checks if token file exists
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("token file does not exist: %v", err)
	}

	// read token from file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %v", err)
	}
	defer file.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(token); err != nil {
		return nil, fmt.Errorf("failed to decode token from file: %v", err)
	}

	if !token.Valid() {
		return nil, fmt.Errorf("invalid token; can't use API without accessToken: %v", err)
	}

	return token, nil
}

// OauthServer starts a server to authenticate the user and returns a token.
// It writes the token to a file if writeTokenFile is true.
func OauthServer() (*oauth2.Token, error) {
	conf := newOauthConf()

	// prepare server
	ctx := context.Background()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	server := &http.Server{Addr: ":9000"}

	// create a channel to receive the authorization code
	codeChan := make(chan string)

	http.HandleFunc("/oauth/callback", handleOauthCallback(codeChan))

	// start server
	errChan := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start server: %v", err)
			return
		}
		errChan <- nil
	}()
	if err := <-errChan; err != nil {
		return nil, err
	}

	// get the OAuth authorization URL
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// redirect user to consent page to ask for permission
	// for the scopes specified above
	fmt.Printf("Your browser has been opened to visit:\n%s\n", url)

	// open user's browser to login page
	if err := browser.OpenURL(url); err != nil {
		return nil, fmt.Errorf("failed to open browser for authentication %s", err.Error())
	}

	// wait for the authorization code to be received
	code := <-codeChan

	// exchange the authorization code for an access token
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code for token: %v", err)
	}

	if !token.Valid() {
		return nil, fmt.Errorf("invalid token; can't use API without accessToken: %v", err)
	}

	// shut down the HTTP server
	if err := server.Shutdown(ctx); err != nil {
		return nil, fmt.Errorf("failed to shut down server: %v", err)
	}

	log.Println(color.CyanString("Authentication successful"))

	return token, nil
}

func handleOauthCallback(codeChan chan string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParts, _ := url.ParseQuery(r.URL.RawQuery)

		// use the authorization code that is pushed to the redirect URL.
		code := queryParts["code"][0]
		log.Printf("code: %s\n", code)

		// write the authorization code to the channel
		codeChan <- code

		msg := "<p><strong>Authentication successful</strong>. You may now close this tab.</p>"
		// send a success message to the browser
		fmt.Fprint(w, msg)
	}
}

func writeTokenToFile(token *oauth2.Token) error {
	// create file with 0600 permissions
	file, err := os.OpenFile("token.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to create token file: %v", err)
	}
	defer file.Close()

	// encode token as JSON and write to file
	if err := json.NewEncoder(file).Encode(token); err != nil {
		return fmt.Errorf("unable to write token to file: %v", err)
	}

	return nil
}
