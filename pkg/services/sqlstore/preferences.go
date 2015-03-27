package sqlstore

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

func init() {
	bus.AddHandler("sql", SavePreferences)
	bus.AddHandler("sql", GetPreferences)
}

func SavePreferences(cmd *m.SavePreferencesCommand) error {
	return inTransaction(func(sess *xorm.Session) error {
		// try get existing prefs
		prefs := m.Preferences{OrgId: cmd.OrgId, UserId: cmd.UserId}
		has, err := sess.Get(&prefs)
		if err != nil {
			return err
		}

		prefs.HomeDashboardId = cmd.HomeDashboardId

		if !has {
			prefs.Created = time.Now()
			if _, err := sess.Insert(&prefs); err != nil {
				return err
			}
		} else {
			if _, err := sess.Id(prefs.Id).Update(&prefs); err != nil {
				return err
			}
		}

		return nil
	})
}

func GetPreferences(query *m.GetPreferencesQuery) error {
	sess := x.NewSession()

	if query.UserId > 0 {
		sess.Where("org_id=? OR (org_id=? AND user_id=?)", query.OrgId, query.UserId)
	}

	query.Result = make([]*m.ApiKey, 0)
	return sess.Find(&query.Result)
}
