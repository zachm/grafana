package setting

import (
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/log"
)

type QuotaLimit struct {
	GlobalLimit int64
	OrgLimit    int64
	UserLimit   int64
}

type QuotaSettings struct {
	Enabled bool
	Limits  map[string]*QuotaLimit
}

func GetDefaultQuotaFor(target string, orgId int64, userId int64) int64 {
	limits, exists := Quota.Limits[target]
	if !exists {
		return -1
	}

	if orgId != 0 {
		return limits.OrgLimit
	} else if userId != 0 {
		return limits.UserLimit
	} else {
		return limits.GlobalLimit
	}
}

func parseQuotaLimit(str string) int64 {
	if str == "NA" {
		return -1
	}

	val, _ := strconv.ParseInt(str, 10, 0)
	return val
}

func readQuotaSettings() {
	// set global defaults.
	cfgSection := Cfg.Section("quota")
	Quota.Enabled = cfgSection.Key("enabled").MustBool(false)

	Quota.Limits = make(map[string]*QuotaLimit)
	for _, key := range cfgSection.Keys() {
		keyName := key.Name()

		if strings.HasPrefix(keyName, "limit_") {
			keyName = strings.TrimPrefix(keyName, "limit_")
			keyName = strings.TrimSuffix(keyName, "s")
			vals := key.Strings(",")
			Quota.Limits[keyName] = &QuotaLimit{
				GlobalLimit: parseQuotaLimit(vals[0]),
				OrgLimit:    parseQuotaLimit(vals[1]),
				UserLimit:   parseQuotaLimit(vals[2]),
			}
		}
	}

	log.Info("Values %v", Quota.Limits)
}
