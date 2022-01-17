package main

import (
        "github.com/gin-gonic/gin"
        "net/http"
        "fmt"
	"os"
	"github.com/zricethezav/go-tdameritrade"
	"golang.org/x/oauth2"
	"context"
	"encoding/json"
)

var (
        router *gin.Engine
)

func init() {
        router = gin.Default()
}

type HTTPHeaderStore struct{}
type TDHandlers struct {
	authenticator *tdameritrade.Authenticator
}

func (s *HTTPHeaderStore) StoreToken(token *oauth2.Token, w http.ResponseWriter, req *http.Request) error {
	// DO NOT DO THIS IN A PRODUCTION ENVIRONMENT!
	// This is just an example.
	// Used signed cookies like those provided by https://github.com/gorilla/securecookie
	http.SetCookie(
		w,
		&http.Cookie{
			Name:    "refreshToken",
			Value:   token.RefreshToken,
			Expires: token.Expiry,
		},
	)
	http.SetCookie(
		w,
		&http.Cookie{
			Name:    "accessToken",
			Value:   token.AccessToken,
			Expires: token.Expiry,
		},
	)
	return nil
}

func (s HTTPHeaderStore) GetToken(req *http.Request) (*oauth2.Token, error) {
	// DO NOT DO THIS IN A PRODUCTION ENVIRONMENT!
	// This is just an example.
	// Used signed cookies like those provided by https://github.com/gorilla/securecookie
	refreshToken, err := req.Cookie("refreshToken")
	if err != nil {
		return nil, err
	}

	accessToken, err := req.Cookie("accessToken")
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
		Expiry:       refreshToken.Expires,
	}, nil
}

func (s HTTPHeaderStore) StoreState(state string, w http.ResponseWriter, req *http.Request) error {
	// DO NOT DO THIS IN A PRODUCTION ENVIRONMENT!
	// This is just an example.
	// Used signed cookies like those provided by https://github.com/gorilla/securecookie
	http.SetCookie(
		w,
		&http.Cookie{
			Name:  "state",
			Value: state,
		},
	)
	return nil
}

func (s HTTPHeaderStore) GetState(req *http.Request) (string, error) {
	// DO NOT DO THIS IN A PRODUCTION ENVIRONMENT!
	// This is just an example.
	// Used signed cookies like those provided by https://github.com/gorilla/securecookie
	cookie, err := req.Cookie("state")
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

//func (h *TDHandlers) Authenticate(w http.ResponseWriter, req *http.Request) {
func (h *TDHandlers) Authenticate(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var req *http.Request = c.Request

	redirectURL, err := h.authenticator.StartOAuth2Flow(w, req)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, redirectURL, http.StatusTemporaryRedirect)
}

//func (h *TDHandlers) Callback(w http.ResponseWriter, req *http.Request) {
func (h *TDHandlers) Callback(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var req *http.Request = c.Request
	ctx := context.Background()
	_, err := h.authenticator.FinishOAuth2Flow(ctx, w, req)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/apis/v1/quote?ticker=SPY", http.StatusTemporaryRedirect)
}

//func (h *TDHandlers) Quote(w http.ResponseWriter, req *http.Request) {
func (h *TDHandlers) Quote(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var req *http.Request = c.Request
	ctx := context.Background()

	client, err := h.authenticator.AuthenticatedClient(ctx, req)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	ticker, ok := req.URL.Query()["ticker"]
	if !ok || len(ticker) == 0 {
		w.Write([]byte("ticker is required"))
		return
	}

	quote, _, err := client.Quotes.GetQuotes(ctx, ticker[0])
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	body, err := json.Marshal(quote)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(body)
}

func (h *TDHandlers) Health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
        c.Writer.Write([]byte(fmt.Sprintf("Jui\n")))
}

func main() {
	authenticator := tdameritrade.NewAuthenticator(
		&HTTPHeaderStore{},
		oauth2.Config{
			ClientID: os.Getenv("clientID"),
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://api.tdameritrade.com/v1/oauth2/token",
				AuthURL:  "https://auth.tdameritrade.com/auth",
			},
			RedirectURL: "https://www.utnetscout.com/apis/v1/callback",
		},
	)
	handlers := &TDHandlers{authenticator: authenticator}

        v := router.Group("/apis/v1")
        {
                //v.GET("/get_user", get_user_secure_new_func)
		v.GET("/authenticate", handlers.Authenticate)
		v.GET("/callback", handlers.Callback)
		v.GET("/quote", handlers.Quote)
		v.GET("/health", handlers.Health)
        }
        router.Run(":80")
}

