package middleware

import (
	"fmt"

	"github.com/Unknwon/macaron"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

type targetLimits struct {
	GlobalLimit int64
	OrgLimit    int64
	UserLimit   int64
}

func Quota(target string) macaron.Handler {
	return func(c *Context) {
		limitReached, err := QuotaReached(c, target)
		if err != nil {
			log.Error(3, "Error: %v", err)
			c.JsonApiErr(500, "failed to get quota", err)
			return
		}
		if limitReached {
			c.JsonApiErr(403, fmt.Sprintf("%s Quota reached", target), nil)
			return
		}
	}
}

func getTargetLimits(quotas []*m.Quota, target string) targetLimits {
	limits := targetLimits{GlobalLimit: -1, UserLimit: -1, OrgLimit: -1}

	for _, q := range quotas {
		if q.OrgId != 0 {
			limits.OrgLimit = q.Limit
		} else if q.UserId != 0 {
			limits.UserLimit = q.Limit
		}
	}

	targetConf, exists := setting.Quota.Limits[target]
	if exists {
		limits.GlobalLimit = targetConf.GlobalLimit
		if limits.OrgLimit == -1 {
			limits.OrgLimit = targetConf.OrgLimit
		}
		if limits.UserLimit == -1 {
			limits.UserLimit = targetConf.UserLimit
		}
	}

	return limits
}

func QuotaReached(c *Context, target string) (bool, error) {
	if !setting.Quota.Enabled {
		return false, nil
	}

	var err error
	query := m.GetQuotasByTargetQuery{Target: target, UserId: c.UserId, OrgId: c.OrgId}
	if err = bus.Dispatch(&query); err != nil {
		return true, err
	}

	limits := getTargetLimits(query.Result, target)
	log.Info("Checking quota: %s, limits: %v", target, limits)

	if target == "session" {
		usedSessions := sessionManager.Count()
		if int64(usedSessions) > limits.GlobalLimit {
			log.Info(fmt.Sprintf("%d sessions active, limit is %d", usedSessions, limits.GlobalLimit))
			return true, nil
		}
	}

	if limits.UserLimit != -1 {
		userQuery := m.IsQuotaReachedQuery{Target: target, UserId: c.UserId, Limit: limits.UserLimit}
		err = bus.Dispatch(&userQuery)
		if err != nil || userQuery.Result {
			return true, err
		}
	}

	if limits.OrgLimit != -1 {
		orgQuery := m.IsQuotaReachedQuery{Target: target, OrgId: c.OrgId, Limit: limits.OrgLimit}
		err = bus.Dispatch(&orgQuery)
		if err != nil || orgQuery.Result {
			return true, err
		}
	}

	if limits.GlobalLimit != -1 {
		globalQuery := m.IsQuotaReachedQuery{Target: target, Limit: limits.GlobalLimit}
		err = bus.Dispatch(&globalQuery)
		if err != nil || globalQuery.Result {
			return true, err
		}
	}

	return false, nil
}
