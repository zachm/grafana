package setting

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/log"
)

type OrgQuota struct {
	User       int64 `target:"org_user"`
	DataSource int64 `target:"data_source"`
	Dashboard  int64 `target:"dashboard"`
	ApiKey     int64 `target:"api_key"`
}

type UserQuota struct {
	Org int64 `target:"org_user"`
}

type GlobalQuota struct {
	Org        int64 `target:"org"`
	User       int64 `target:"user"`
	DataSource int64 `target:"data_source"`
	Dashboard  int64 `target:"dashboard"`
	ApiKey     int64 `target:"api_key"`
	Session    int64 `target:"-"`
}

func (q *OrgQuota) ToMap() map[string]int64 {
	return quotaToMap(*q)
}

func (q *UserQuota) ToMap() map[string]int64 {
	return quotaToMap(*q)
}

func (q *GlobalQuota) ToMap() map[string]int64 {
	return quotaToMap(*q)
}

func quotaToMap(q interface{}) map[string]int64 {
	qMap := make(map[string]int64)
	typ := reflect.TypeOf(q)
	val := reflect.ValueOf(q)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := field.Tag.Get("target")
		if name == "" {
			name = field.Name
		}
		if name == "-" {
			continue
		}
		value := val.Field(i)
		qMap[name] = value.Int()
	}
	return qMap
}

type QuotaLimit struct {
	GlobalLimit int64
	OrgLimit    int64
	UserLimit   int64
}

type QuotaSettings struct {
	Enabled bool
	Limits  map[string]*QuotaLimit
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

	// // per ORG Limits
	// Quota.Org = &OrgQuota{
	// 	User:       quota.Key("org_user").MustInt64(10),
	// 	DataSource: quota.Key("org_data_source").MustInt64(10),
	// 	Dashboard:  quota.Key("org_dashboard").MustInt64(10),
	// 	ApiKey:     quota.Key("org_api_key").MustInt64(10),
	// }
	//
	// // per User limits
	// Quota.User = &UserQuota{
	// 	Org: quota.Key("user_org").MustInt64(10),
	// }
	//
	// // Global Limits
	// Quota.Global = &GlobalQuota{
	// 	User:       quota.Key("global_user").MustInt64(-1),
	// 	Org:        quota.Key("global_org").MustInt64(-1),
	// 	DataSource: quota.Key("global_data_source").MustInt64(-1),
	// 	Dashboard:  quota.Key("global_dashboard").MustInt64(-1),
	// 	ApiKey:     quota.Key("global_api_key").MustInt64(-1),
	// 	Session:    quota.Key("global_session").MustInt64(-1),
	// }

}
