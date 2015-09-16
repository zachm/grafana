package sqlstore

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/log"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

func init() {
	bus.AddHandler("sql", GetQuotaInfoQuery)
	bus.AddHandler("sql", GetQuotasQuery)
	bus.AddHandler("sql", AddQuota)
	bus.AddHandler("sql", RemoveQuota)
}

type quotaCount struct {
	Count int64
}

func GetQuotaInfoQuery(query *m.GetQuotaInfoQuery) error {
	var quota m.Quota
	sess := x.Table("quota").Where("target=? AND org_id=? AND user_id=?", query.Target, query.OrgId, query.UserId)
	if exists, err := sess.Get(&quota); err != nil {
		return err
	} else if !exists {
		quota.Limit = setting.GetDefaultQuotaFor(query.Target, query.OrgId, query.UserId)
	}

	// if limit -1 skip checking used
	if quota.Limit == -1 {
		query.ResultLimit = -1
		query.ResultUsed = -1
		return nil
	}

	params := make([]interface{}, 0)
	rawSql := fmt.Sprintf("SELECT COUNT(*) as count from %s", dialect.Quote(query.Target))

	if query.OrgId != 0 {
		rawSql += " WHERE org_id=?"
		params = append(params, query.OrgId)
	} else if query.UserId != 0 {
		rawSql += "WHERE user_id=?"
		params = append(params, query.UserId)
	}

	resp := make([]*quotaCount, 0)
	if err := x.Sql(rawSql, params...).Find(&resp); err != nil {
		return err
	}

	log.Debug("sqlstore: IsQuotaReachedQuery (OrgId, %d, UserId: %d) Limit: %d, Count: %d", query.OrgId, query.UserId, quota.Limit, resp[0].Count)

	query.ResultLimit = quota.Limit
	query.ResultUsed = resp[0].Count
	return nil
}

func GetQuotasQuery(query *m.GetQuotasQuery) error {
	query.Result = make([]*m.Quota, 0)
	sess := x.Limit(100).Table("quota")

	if query.OrgId != 0 {
		sess.Where("org_id=?", query.OrgId)
	} else if query.UserId != 0 {
		sess.Where("user_id=?", query.UserId)
	}

	if err := sess.Find(&query.Result); err != nil {
		return err
	}

	return nil
}

func AddQuota(cmd *m.AddQuotaCommand) error {
	return inTransaction(func(sess *xorm.Session) error {
		t := m.Quota{
			Target:  cmd.Target,
			UserId:  cmd.UserId,
			OrgId:   cmd.OrgId,
			Limit:   cmd.Limit,
			Created: time.Now(),
			Updated: time.Now(),
		}

		if _, err := sess.Insert(&t); err != nil {
			return err
		}
		cmd.Result = &t
		return nil
	})
}

func RemoveQuota(cmd *m.RemoveQuotaCommand) error {
	return inTransaction(func(sess *xorm.Session) error {
		var rawSql = "DELETE FROM quota WHERE id=?"
		_, err := sess.Exec(rawSql, cmd.Id)
		return err
	})
}
