package user

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"clevergo.tech/clevergo"
	"github.com/coreos/go-oidc"
	"github.com/headwindfly/go-auth0-web-app/internal/core"
)

func (h *Handler) callback(c *clevergo.Context) error {
	ctx := c.Context()
	state := h.sessionManager.GetString(ctx, "auth_state")
	if c.QueryParam("state") != state {
		return errors.New("invalid state parameter")
	}

	authenticator, err := core.NewAuthenticator(ctx)
	if err != nil {
		return err
	}

	token, err := authenticator.Config.Exchange(ctx, c.QueryParam("code"))
	if err != nil {
		return err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return errors.New("No id_token field in oauth2 token.")
	}

	oidcConfig := &oidc.Config{
		ClientID: os.Getenv("AUTH0_CLIENT_ID"),
	}
	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(ctx, rawIDToken)
	if err != nil {
		return fmt.Errorf("failed to verify ID Token: %s", err.Error())
	}

	// Getting now the userInfo
	var user core.User
	if err := idToken.Claims(&user); err != nil {
		return err
	}
	user.IDToken = rawIDToken
	user.AccessToken = token.AccessToken
	h.sessionManager.Put(ctx, "auth_user", user)

	return c.Redirect(http.StatusSeeOther, "/dashboard")
}
