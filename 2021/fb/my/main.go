package main

import (
	"os"
	"fmt"
	"net/http"
	"golang.org/x/oauth2"
	oauth2fb "golang.org/x/oauth2/facebook"
	fb "github.com/huandu/facebook"
)

var (
    // You must register the app at https://github.com/settings/applications
    // Set callback to http://127.0.0.1:7000/github_oauth_cb
    // Set ClientId and ClientSecret to
    oauthConf = &oauth2.Config{
        ClientID:     "316542315986013",
        ClientSecret: "6dbee45cf73f688053d95f24527720c9",
        // select level of access you want https://developer.github.com/v3/oauth/#scopes
	Scopes:   []string{"email"},
	Endpoint: oauth2fb.Endpoint,
	RedirectURL: "https://smartpatty1.appspot.com/fb_oauth_cb",
    }
    // random string for oauth2 API calls to protect against CSRF
    oauthStateString = "thisshouldberandom"
)

func indexhandle(w http.ResponseWriter, r *http.Request) {
	indexhtml := `<html><body><a href="/login">Facebook Alone</a></body></html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexhtml))
}

func loginhandle(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func cbhandle(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Create a client to manage access token life cycle.
	client := oauthConf.Client(oauth2.NoContext, token)

	// Use OAuth2 client with session.
	session := &fb.Session{
	    Version:    "v2.4",
	    HttpClient: client,
	}

	// Use session.
	session.Get("/me", nil)
	fmt.Fprint(w, "Smart Patty 1")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", indexhandle)
	http.HandleFunc("/login", loginhandle)
	http.HandleFunc("/fb_oauth_cb", cbhandle)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Failed to ListenAndServe")
	}
}
	
