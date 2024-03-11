package dribbble

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
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
	public = "public"
	upload = "upload"
)

// Default token file
const tokenFile = "token.json"

// endpoint is Dribbbles's OAuth 2.0 default endpoint.
var endpoint = oauth2.Endpoint{
	AuthURL:  "https://dribbble.com/oauth/authorize",
	TokenURL: "https://dribbble.com/oauth/token",
}

// newOauthConf returns a new oauth2.Config with values from environment variables.
func newOauthConf() *oauth2.Config {
	scope := []string{public}
	if os.Getenv("WRITE_SCOPE") == "true" {
		scope = append(scope, upload)
	}

	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       scope,
		Endpoint:     endpoint,
	}
}

// ------------------------------------------------------------------------

// OauthStart checks if a token file exists and starts a server to authenticate
// the user if it does not. It returns an error if there is a problem with the
// token file or the server. On success it sets the token in the config and
// writes the token to a file if writeTokenFile is true.
func OauthStart(config *Config, filePath ...string) error {
	// validate number of arguments
	if len(filePath) > 1 {
		return fmt.Errorf("expected at most one argument; got %d", len(filePath))
	}

	f := tokenFile          // defaults
	if len(filePath) == 1 { // use provided file path if available
		f = filePath[0]
	}

	// get token from file
	token, err := OauthFile(f)
	if err != nil {
		if strings.Contains(err.Error(), "token file does not exist") {
			// no token file so start server to authenticate user
			token, err = OauthServer()
			if err != nil {
				return err // error with auth server
			}
		} else {
			return err // some other error
		}
	}

	// add token to config
	config.Token = token

	// write token to file
	if config.Flags.TokenFile {
		if err := writeTokenToFile(token, f); err != nil {
			return fmt.Errorf("failed to write token to file: %v", err)
		}
	}

	return nil
}

// OauthFile returns a token from a file if it exists. If no file path is
// provided or if the file does not exist, it returns nil.
func OauthFile(filePath string) (*oauth2.Token, error) {
	if filePath == "" {
		return nil, errors.New("expected a file path; got empty string")
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
// It writes the token to a file if TokenFile flag is true.
func OauthServer() (*oauth2.Token, error) {
	conf := newOauthConf()

	// prepare server
	ctx := context.Background()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	urlComps, err := extractURLComponents(conf.RedirectURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract URL components: %v", err)
	}

	// set server
	addr := fmt.Sprintf("%s:%s", urlComps.Domain, urlComps.Port)
	server := &http.Server{Addr: addr}

	// create a channel to receive the authorization code
	codeChan := make(chan string)

	http.HandleFunc(urlComps.Path, handleOauthCallback(codeChan))

	// start server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

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

// ------------------------------------------------------------------------

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

func writeTokenToFile(token *oauth2.Token, filePath string) error {
	if filePath == "" {
		return errors.New("expected a file path; got empty string")
	}

	// create file with 0600 permissions
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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

// ------------------------------------------------------------------------

// urlComponents holds the extracted components of the redirect URL.
type urlComponents struct {
	Proto  string
	Domain string
	Port   string
	Path   string
}

func extractURLComponents(urlStr string) (*urlComponents, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	components := &urlComponents{
		Proto: parsedURL.Scheme,
		Path:  parsedURL.Path,
	}

	// split host into domain and port
	hostParts := strings.Split(parsedURL.Host, ":")
	components.Domain = hostParts[0]
	if len(hostParts) > 1 {
		components.Port = hostParts[1]
	}

	return components, nil
}
