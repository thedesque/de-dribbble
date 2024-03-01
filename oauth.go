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

	"github.com/cli/browser"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		color.Red("error: could not load .env file")
		os.Exit(1)
	}

	requiredEnvVars := []string{"CLIENT_ID", "CLIENT_SECRET", "REDIRECT_URL"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			color.Red("error: %s is not set", envVar)
			os.Exit(1)
		}
	}

	writeScope := os.Getenv("WRITE_SCOPE")
	if writeScope != "true" && writeScope != "false" {
		color.Red("error: WRITE_SCOPE is not set; should either be true or false")
		os.Exit(1)
	}
}

// Dribbble scopes
const (
	Public = "public"
	Upload = "upload"
)

// Endpoint is Dribbbles's OAuth 2.0 default endpoint.
var Endpoint = oauth2.Endpoint{
	AuthURL:  "https://dribbble.com/oauth/authorize",
	TokenURL: "https://dribbble.com/oauth/token",
}

func Auth() *oauth2.Token {
	scope := []string{Public}
	if os.Getenv("WRITE_SCOPE") == "true" {
		scope = append(scope, Upload)
	}

	// setup dribbble OAuth2 conf
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes:       scope,
		Endpoint:     Endpoint,
	}

	// start server
	ctx := context.Background()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, sslcli)

	server := &http.Server{Addr: ":9000"}

	// create a channel to receive the authorization code
	codeChan := make(chan string)

	http.HandleFunc("/oauth/callback", handleOauthCallback(ctx, conf, codeChan))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// get the OAuth authorization URL
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// redirect user to consent page to ask for permission
	// for the scopes specified above
	fmt.Printf("Your browser has been opened to visit:\n%s\n", url)

	// open user's browser to login page
	if err := browser.OpenURL(url); err != nil {
		panic(fmt.Errorf("failed to open browser for authentication %s", err.Error()))
	}

	// wait for the authorization code to be received
	code := <-codeChan

	// exchange the authorization code for an access token
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Failed to exchange authorization code for token: %v", err)
	}

	if !token.Valid() {
		log.Fatalf("Can't get source information without accessToken: %v", err)
		return nil
	}

	// write the access token to a file
	// if err := writeTokenToFile(token); err != nil {
	// 	log.Fatalf("Failed to write token to file: %v", err)
	// }

	// shut down the HTTP server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to shut down server: %v", err)
	}

	log.Println(color.CyanString("Authentication successful"))

	return token
}

// todo: not using ctx or conf yet
func handleOauthCallback(ctx context.Context, conf *oauth2.Config, codeChan chan string) func(w http.ResponseWriter, r *http.Request) {
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

// todo: not in use, check if we need this
func writeTokenToFile(token *oauth2.Token) error {
	// create file with 0600 permissions
	file, err := os.OpenFile("token.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to create token file: %v", err)
	}
	defer file.Close()

	// encode token as JSON and write to file
	if err := json.NewEncoder(file).Encode(token); err != nil {
		return fmt.Errorf("Unable to write token to file: %v", err)
	}

	return nil
}
