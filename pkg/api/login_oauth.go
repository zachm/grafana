package api

import (
	"fmt"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/metrics"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/social"
)

func OAuthLogin(ctx *middleware.Context) Response {
	if setting.OAuthService == nil {
		return HtmlErrorView(404, "login.OAuthLogin(oauth service not enabled)", nil)
	}

	name := ctx.Params(":name")
	connect, ok := social.SocialMap[name]
	if !ok {
		return HtmlErrorView(404, "login.OAuthLogin(social login not enabled)", nil)
	}

	code := ctx.Query("code")
	if code == "" {
		return Redirect(connect.AuthCodeURL("", oauth2.AccessTypeOnline))
	}

	// handle call back
	token, err := connect.Exchange(oauth2.NoContext, code)
	if err != nil {
		return HtmlErrorView(500, "login.OAuthLogin(NewTransportWithCode)", err)
	}

	log.Trace("login.OAuthLogin(Got token)")

	userInfo, err := connect.UserInfo(token)
	if err != nil {
		if err == social.ErrMissingTeamMembership {
			return Redirect("/login?failedMsg=" + url.QueryEscape("Required Github team membership not fulfilled"))
		} else if err == social.ErrMissingOrganizationMembership {
			return Redirect("/login?failedMsg=" + url.QueryEscape("Required Github organization membership not fulfilled"))
		}
		return HtmlErrorView(500, fmt.Sprintf("login.OAuthLogin(get info from %s)", name), err)
	}

	log.Trace("login.OAuthLogin(social login): %s", userInfo)

	// validate that the email is allowed to login to grafana
	if !connect.IsEmailAllowed(userInfo.Email) {
		log.Info("OAuth login attempt with unallowed email, %s", userInfo.Email)
		return Redirect("/login?failedMsg=" + url.QueryEscape("Required email domain not fulfilled"))
	}

	userQuery := m.GetUserByLoginQuery{LoginOrEmail: userInfo.Email}
	err = bus.Dispatch(&userQuery)

	// create account if missing
	if err == m.ErrUserNotFound {
		if !connect.IsSignupAllowed() {
			return Redirect("/login?failedMsg=" + url.QueryEscape("OAuth signup is not allowed"))
		}

		if resp := CheckQuota(ctx, middleware.QuotaDefUsers); resp != nil {
			return HtmlErrorView(403, "User Quota Reached", nil)
		}
		cmd := m.CreateUserCommand{
			Login:   userInfo.Email,
			Email:   userInfo.Email,
			Name:    userInfo.Name,
			Company: userInfo.Company,
		}
		if err = bus.Dispatch(&cmd); err != nil {
			return HtmlErrorView(500, "Failed to create account", err)
		}
		userQuery.Result = &cmd.Result
	} else if err != nil {
		return HtmlErrorView(500, "Unexpected error", err)
	}

	// login
	loginUserWithUser(userQuery.Result, ctx)

	metrics.M_Api_Login_OAuth.Inc(1)

	return Redirect("/")
}
