package main

import (
	"os"
	"net/http"
	"fmt"
	"golang.org/x/oauth2"
	oauth2fb "golang.org/x/oauth2/facebook"
	fb "github.com/huandu/facebook"
)

var (
	authconf = &oauth2.Config {
		ClientID:	ClientID,
		ClientSecret:	ClientSecret,
		RedirectURL:	RedirectURL,
		Scopes:		[]string{"email"},
		Endpoint:	oauth2fb.Endpoint,
	}
	statestring = "Jeevan"	
)

func indexhandle(w http.ResponseWriter, r *http.Request) {
	indexhtml := `<html><body><a href="/login">Login with Facebook </a></body></html>`
	
	w.Header().Set("Content-Type", "charset=utf-8, text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexhtml))
}

func loginhandle(w http.ResponseWriter, r *http.Request) {
	url := authconf.AuthCodeURL(statestring, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func fb_oauth_cb_handle(w http.ResponseWriter, r *http.Request) {
	authcode := r.FormValue("code")
	accesstoken, err := authconf.Exchange(oauth2.NoContext, authcode)
	if err != nil {
		fmt.Fprint(w, "Failed to get accesstoken")
	}

	fmt.Fprint(w, "Access Token", accesstoken.AccessToken)

	res, err1 := fb.Get("/me", fb.Params{
				"fields": "first_name, email",
				"access_token": accesstoken.AccessToken,
			})
	if err1 != nil {
		fmt.Fprint(w, "Failed fb.Get")
	}
	fmt.Fprint(w, "Name: ", res["first_name"])
	fmt.Fprint(w, "email: ", res["email"])

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", indexhandle)
	http.HandleFunc("/login", loginhandle)
	http.HandleFunc("/fb_oauth_cb", fb_oauth_cb_handle)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Failed Listen and Serve")
	}	
}
