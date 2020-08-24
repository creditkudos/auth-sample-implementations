package main

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/pug"
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
	redirectURL = os.Getenv("REDIRECTURL")
	clientID = os.Getenv("CLIENTID")
	clientSecret = os.Getenv("CLIENTSECRET")
	if redirectURL == "" || clientID == "" || clientSecret == "" {
		panic("Set env variables before running")
	}

	engine := pug.New("./views", ".pug")
	app := fiber.New(&fiber.Settings{Views: engine})
	oauthClient := createOauthClient()

	createRoutes(app, oauthClient)

	if err := app.Listen(3000); err != nil {
		fmt.Println("App errored: ", err.Error())
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
			fmt.Println("Failed to load main page:", err)
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
			fmt.Println("Failed to create jwt token:", err)
			return
		}

		// Generate the redirection URI and redirect
		customerToken := oauth2.SetAuthURLParam("customer_token", tokenString)
		url := oauthConfig.AuthCodeURL("abc", oauth2.AccessTypeOnline, customerToken)
		fmt.Println("Redirect URL:", url)
		c.Redirect(url)
	})

	// Setup return callback
	app.Get("/callback", func(c *fiber.Ctx) {
		code := c.Query("code")
		fmt.Println("Response code:", code)
		ctx := c.Context()
		token, err := oauthConfig.Exchange(ctx, code)
		if err != nil {
			handleError(fmt.Errorf("unable to exchange for token: %s", err), c)
			return
		}

		fmt.Println("Access Token:", token.AccessToken)
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
			handleError(fmt.Errorf("unable to render /complete: %s", err), c)
		}
	})
}

// Handle errors and render error response
func handleError(e error, c *fiber.Ctx) {
	fmt.Println(e.Error())
	if err := c.Render("error", fiber.Map{"message": e.Error()}); err != nil {
		panic(err)
	}
}
