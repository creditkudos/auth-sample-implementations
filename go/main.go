package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/template/pug"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

/* #nosec G1010*/
const (
	accessTokenURL   = "https://api.creditkudos-staging.com/v2/oauth/token"
	authorizationURL = "https://auth.creditkudos-staging.com/oauth/authorize"
)

// Update the env in docker with your redirect URL and client id/secret
var (
	redirectURL  = ""
	clientID     = ""
	clientSecret = ""
)

func main() {
	redirectURL = os.Getenv("REDIRECT_URL")
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	if redirectURL == "" || clientID == "" || clientSecret == "" {
		log.Fatal("Set env variables before running")
	}

	engine := pug.New("./views", ".pug")
	app := fiber.New(&fiber.Settings{Views: engine})
	oauthClient := createOauthClient()

	createRoutes(app, oauthClient)

	if err := app.Listen(3000); err != nil {
		log.Fatal(fmt.Errorf("App errored: %w", err))
	}
}

func createOauthClient() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizationURL,
			TokenURL: accessTokenURL,
		},
		RedirectURL: redirectURL,
	}
}

func createRoutes(app *fiber.App, oauthConfig *oauth2.Config) {
	// Setup landing page
	app.Get("/", func(c *fiber.Ctx) {
		if err := c.Render(
			"index",
			fiber.Map{
				"title":   "Example Go implementation",
				"message": "Click below to start a journey",
			},
		); err != nil {
			log.Println("Failed to load main page:", err)
			return
		}
	})

	// Setup redirection
	app.Get("/redirect", func(c *fiber.Ctx) {
		// Create and sign JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss":              clientID,
			"sub":              "customer",
			"iat":              1000,
			"jti":              "234342",
			"email":            "sam.pull@example.com",
			"first_name":       "Samuel",
			"last_name":        "Pull",
			"custom_reference": "SP-123",
			"date_of_birth":    "1985-10-25",
			"postcode":         "XY12AB",
		})
		tokenString, err := token.SignedString([]byte(clientSecret))
		if err != nil {
			handleError(err, c)
			return
		}

		// Generate the redirection URI and redirect
		customerToken := oauth2.SetAuthURLParam("customer_token", tokenString)
		// The first parameter here is the state, which is passed back at the end of the oauth journey
		url := oauthConfig.AuthCodeURL("sample state", oauth2.AccessTypeOnline, customerToken)
		log.Println("Redirect URL: ", url)
		c.Redirect(url)
	})

	// Setup return callback
	app.Get("/callback", func(c *fiber.Ctx) {
		code := c.Query("code")
		log.Println("Response code:", code)
		ctx := c.Context()
		token, err := oauthConfig.Exchange(ctx, code)
		if err != nil {
			handleError(fmt.Errorf("unable to exchange for token: %w", err), c)
			return
		}

		log.Println("Access Token:", token.AccessToken)
		// oauthConfig.Client returns an HTTP client that automatically uses and
		// renews the token
		oauthConfig.Client(ctx, token)

		c.Redirect("/complete")
	})

	// Setup on complete
	app.Get("/complete", func(c *fiber.Ctx) {
		if err := c.Render(
			"index",
			fiber.Map{
				"title":   "Complete!",
				"message": "Journey complete! Click below to start another",
			},
		); err != nil {
			handleError(fmt.Errorf("unable to render /complete: %w", err), c)
		}
	})
}

// Handle errors and render error response
func handleError(e error, c *fiber.Ctx) {
	log.Println(e.Error())
	if err := c.Render("error", fiber.Map{"message": e.Error()}); err != nil {
		c.SendStatus(http.StatusInternalServerError)
		log.Fatal(fmt.Errorf("Failed to render error page: %w", err))
	}
}
