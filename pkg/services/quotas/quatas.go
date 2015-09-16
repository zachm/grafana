package quotas

import "github.com/grafana/grafana/pkg/setting"

func Init() {
	if !setting.Quota.Enabled {
		return
	}
}
