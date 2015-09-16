package middleware

import (
	"fmt"

	"github.com/Unknwon/macaron"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

type QuotaScope struct {
	Target  string
	Handler QuotaCheckHandler
}

type QuotaDef struct {
	Name   string
	Scopes []QuotaScope
}

type QuotaCheckHandler func(c *Context, scope *QuotaScope) (bool, error)

var QuotaDefOrgs *QuotaDef = &QuotaDef{
	Name: "Organization count",
	Scopes: []QuotaScope{
		{Target: "org", Handler: globalQuotaCheck},
		{Target: "org_user", Handler: userQuotaCheck},
	},
}

var QuotaDefDashboards *QuotaDef = &QuotaDef{
	Name: "Dashboard count",
	Scopes: []QuotaScope{
		{Target: "dashboard", Handler: globalQuotaCheck},
		{Target: "dashboard", Handler: orgQuotaCheck},
	},
}

var QuotaDefDataSources *QuotaDef = &QuotaDef{
	Name: "Data sources count",
	Scopes: []QuotaScope{
		{Target: "data_source", Handler: globalQuotaCheck},
		{Target: "data_source", Handler: orgQuotaCheck},
	},
}

var QuotaDefUsers *QuotaDef = &QuotaDef{
	Name: "User count",
	Scopes: []QuotaScope{
		{Target: "user", Handler: globalQuotaCheck},
		{Target: "org_user", Handler: orgQuotaCheck},
	},
}

var QuotaDefApiKeys *QuotaDef = &QuotaDef{
	Name: "Api keys count",
	Scopes: []QuotaScope{
		{Target: "api_key", Handler: globalQuotaCheck},
		{Target: "api_key", Handler: orgQuotaCheck},
	},
}

var QuotaDefSessions *QuotaDef = &QuotaDef{
	Name:   "Concurrent users",
	Scopes: []QuotaScope{{Handler: sessionQuotaCheck}},
}

func Quota(def *QuotaDef) macaron.Handler {
	return func(c *Context) {
		limitReached, err := QuotaReached(c, def)
		if err != nil {
			c.JsonApiErr(500, "failed to get quota", err)
			return
		}
		if limitReached {
			c.JsonApiErr(403, fmt.Sprintf("%s Quota reached", def.Name), nil)
			return
		}
	}
}

func QuotaReached(c *Context, quota *QuotaDef) (bool, error) {
	if !setting.Quota.Enabled {
		return false, nil
	}
	for _, scope := range quota.Scopes {
		if overLimit, err := scope.Handler(c, &scope); err != nil || overLimit {
			return overLimit, err
		}
	}
	return false, nil
}

func globalQuotaCheck(c *Context, scope *QuotaScope) (bool, error) {
	checkQuery := m.IsQuotaReachedQuery{Target: scope.Target}
	if err := bus.Dispatch(&checkQuery); err != nil || checkQuery.Result {
		return true, err
	}
	return false, nil
}

func orgQuotaCheck(c *Context, scope *QuotaScope) (bool, error) {
	checkQuery := m.IsQuotaReachedQuery{Target: scope.Target, OrgId: c.OrgId}
	if err := bus.Dispatch(&checkQuery); err != nil || checkQuery.Result {
		return true, err
	}
	return false, nil
}

func userQuotaCheck(c *Context, scope *QuotaScope) (bool, error) {
	checkQuery := m.IsQuotaReachedQuery{Target: scope.Target, UserId: c.UserId}
	if err := bus.Dispatch(&checkQuery); err != nil || checkQuery.Result {
		return true, err
	}
	return false, nil
}

func sessionQuotaCheck(c *Context, scope *QuotaScope) (bool, error) {
	usedSessions := int64(sessionManager.Count())
	limit := setting.GetDefaultQuotaFor("sessions", 0, 0)
	if limit != -1 && usedSessions > limit {
		log.Info(fmt.Sprintf("%d sessions active, limit is %d", usedSessions, limit))
		return true, nil
	}

	return false, nil
}
