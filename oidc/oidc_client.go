// Package oidc evaluates https://github.com/zitadel/oidc
package oidc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/oidc/pkg/client/rp"
	httpHelper "github.com/zitadel/oidc/pkg/http"
	"github.com/zitadel/oidc/pkg/oidc"
)

var (
	callbackPath = "/auth/callback"
	key          = []byte("test1234test1234")
)

// RunClient starts a http server do illustrate OIDC Flow
// Before you have to start oidc op server
// oidc discovery http://localhost:9998/.well-known/openid-configuration
// go run github.com/zitadel/oidc/example/server
// CLIENT_ID=web CLIENT_SECRET=secret ISSUER=http://localhost:9998/ SCOPES="openid profile" PORT=9999 go run github.com/zitadel/oidc/example/client/app
func RunClient(args []string) {
	log.Info().Msgf("Start client with args %v", args)
	clientID := "web"
	clientSecret := "secret"
	issuer := "http://localhost:9998/"
	port := 9999
	scopes := []string{"openid", "profile"}
	// keyPath := ""

	redirectURI := fmt.Sprintf("http://localhost:%v%v", port, callbackPath)
	cookieHandler := httpHelper.NewCookieHandler(key, key, httpHelper.WithUnsecure())

	options := []rp.Option{
		rp.WithCookieHandler(cookieHandler),
		rp.WithVerifierOpts(rp.WithIssuedAtOffset(5 * time.Second)),
	}
	// if clientSecret == "" {
	options = append(options, rp.WithPKCE(cookieHandler))
	// }
	// if keyPath != "" {options = append(options, rp.WithJWTProfile(rp.SignerFromKeyPath(keyPath)))}

	provider, err := rp.NewRelyingPartyOIDC(issuer, clientID, clientSecret, redirectURI, scopes, options...)
	if err != nil {
		logrus.Fatalf("error creating provider %s", err.Error())
	}

	//generate some state (representing the state of the user in your application,
	//e.g. the page where he was before sending him to login
	state := func() string {
		return uuid.New().String()
	}

	//register the AuthURLHandler at your preferred path
	//the AuthURLHandler creates the auth request and redirects the user to the auth server
	//including state handling with secure cookie and the possibility to use PKCE
	http.Handle("/login", rp.AuthURLHandler(state, provider))

	//for demonstration purposes the returned userinfo response is written as JSON object onto response
	marshalUserinfo := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens, state string, rp rp.RelyingParty, info oidc.UserInfo) {
		data, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
	}

	//you could also just take the access_token and id_token without calling the userinfo endpoint:
	//
	//marshalToken := func(w http.ResponseWriter, r *http.Request, tokens *oidc.Tokens, state string, rp rp.RelyingParty) {
	//	data, err := json.Marshal(tokens)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	w.Write(data)
	//}

	//register the CodeExchangeHandler at the callbackPath
	//the CodeExchangeHandler handles the auth response, creates the token request and calls the callback function
	//with the returned tokens from the token endpoint
	//in this example the callback function itself is wrapped by the UserinfoCallback which
	//will call the Userinfo endpoint, check the sub and pass the info into the callback function
	http.Handle(callbackPath, rp.CodeExchangeHandler(rp.UserinfoCallback(marshalUserinfo), provider))

	//if you would use the callback without calling the userinfo endpoint, simply switch the callback handler for:
	//
	//http.Handle(callbackPath, rp.CodeExchangeHandler(marshalToken, provider))
	log.Info().Msgf("listening on %s://localhost:%d/", "http", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal().Err(err).Msg("Serve failed")
	}
}
