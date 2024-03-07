package dribbble

import (
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("%s: %v",
			color.HiRedString("error"),
			color.WhiteString("could not load .env file"),
		)
	}

	requiredEnvVars := []string{"CLIENT_ID", "CLIENT_SECRET", "REDIRECT_URL"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("%s: %v",
				color.HiRedString("error"),
				color.WhiteString("%s is not set", envVar),
			)
		}
	}

	writeScope := os.Getenv("WRITE_SCOPE")
	if writeScope != "true" && writeScope != "false" {
		log.Printf("%s: %v\n",
			color.HiYellowString("warning"),
			color.WhiteString("WRITE_SCOPE is not set; should be a bool; setting to false"),
		)
		os.Setenv("WRITE_SCOPE", "false")
	}
}

// Config for Dribbble Auth
type Config struct {
	Flags      Flags
	Token      *oauth2.Token
	HTTPClient *http.Client
}

// Flags for Dribbble Client
type Flags struct {
	Verbose    bool
	TokenFile  bool
	WriteScope bool
}

// NewConfig for Dribbble Client
func NewConfig() *Config {
	return &Config{
		Token:      nil,
		HTTPClient: http.DefaultClient,
		Flags: Flags{
			Verbose:    true,
			TokenFile:  true,
			WriteScope: os.Getenv("WRITE_SCOPE") == "true", // bool check; conversion
		},
	}
}
