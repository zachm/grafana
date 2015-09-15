package middleware

import (
	"testing"

	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMiddlewareQuota(t *testing.T) {

	Convey("Given the grafana quota middleware", t, func() {
		setting.Quota = setting.QuotaSettings{
			Enabled: true,
			Limits: map[string][]int64{
				"orgs":        []int64{5, 5},
				"users":       []int64{5, 5},
				"dashboards":  []int64{5, 5},
				"datasources": []int64{5, 5},
				"apikeys":     []int64{5, 5},
				"sessions":    []int64{20},
			},
		}

		middlewareScenario("with user not logged in", func(sc *scenarioContext) {
			bus.AddHandler("globalQuota", func(query *m.GetGlobalQuotaByTargetQuery) error {
				query.Result = &m.GlobalQuotaDTO{
					Target: query.Target,
					Limit:  query.Default,
					Used:   4,
				}
				return nil
			})
			Convey("global user quota not reached", func() {
				sc.m.Get("/user", Quota("users"), sc.defaultHandler)
				sc.fakeReq("GET", "/user").exec()
				So(sc.resp.Code, ShouldEqual, 200)
			})
			Convey("global quota reached", func() {
				setting.Quota.Limits["users"][0] = 4
				sc.m.Get("/user", Quota("users"), sc.defaultHandler)
				sc.fakeReq("GET", "/user").exec()
				So(sc.resp.Code, ShouldEqual, 403)
			})
			Convey("global session quota not reached", func() {
				sc.m.Get("/user", Quota("sessions"), sc.defaultHandler)
				sc.fakeReq("GET", "/user").exec()
				So(sc.resp.Code, ShouldEqual, 200)
			})
			Convey("global session quota reached", func() {
				setting.Quota.Limits["sessions"][0] = 1
				sc.m.Get("/user", Quota("sessions"), sc.defaultHandler)
				sc.fakeReq("GET", "/user").exec()
				So(sc.resp.Code, ShouldEqual, 403)
			})
		})

		middlewareScenario("with user logged in", func(sc *scenarioContext) {
			// log us in, so we have a user_id and org_id in the context
			sc.fakeReq("GET", "/").handler(func(c *Context) {
				c.Session.Set(SESS_KEY_USERID, int64(12))
			}).exec()

			bus.AddHandler("test", func(query *m.GetSignedInUserQuery) error {
				query.Result = &m.SignedInUser{OrgId: 2, UserId: 12}
				return nil
			})
			bus.AddHandler("globalQuota", func(query *m.GetGlobalQuotaByTargetQuery) error {
				query.Result = &m.GlobalQuotaDTO{
					Target: query.Target,
					Limit:  query.Default,
					Used:   4,
				}
				return nil
			})
			bus.AddHandler("userQuota", func(query *m.GetUserQuotaByTargetQuery) error {
				query.Result = &m.UserQuotaDTO{
					Target: query.Target,
					Limit:  query.Default,
					Used:   4,
				}
				return nil
			})
			bus.AddHandler("orgQuota", func(query *m.GetOrgQuotaByTargetQuery) error {
				query.Result = &m.OrgQuotaDTO{
					Target: query.Target,
					Limit:  query.Default,
					Used:   4,
				}
				return nil
			})
			Convey("global datasource quota reached", func() {
				setting.Quota.Limits["datasources"][0] = 4
				sc.m.Get("/ds", Quota("datasources"), sc.defaultHandler)
				sc.fakeReq("GET", "/ds").exec()
				So(sc.resp.Code, ShouldEqual, 403)
			})
			Convey("user Org quota not reached", func() {
				setting.Quota.Limits["orgs"][1] = 5
				sc.m.Get("/org", Quota("orgs"), sc.defaultHandler)
				sc.fakeReq("GET", "/org").exec()
				So(sc.resp.Code, ShouldEqual, 200)
			})
			Convey("user Org quota reached", func() {
				setting.Quota.Limits["orgs"][1] = 4
				sc.m.Get("/org", Quota("orgs"), sc.defaultHandler)
				sc.fakeReq("GET", "/org").exec()
				So(sc.resp.Code, ShouldEqual, 403)
			})
			Convey("org dashboard quota not reached", func() {
				setting.Quota.Limits["dashboards"][1] = 10
				sc.m.Get("/dashboard", Quota("dashboard"), sc.defaultHandler)
				sc.fakeReq("GET", "/dashboard").exec()
				So(sc.resp.Code, ShouldEqual, 200)
			})
			Convey("org dashboard quota reached", func() {
				setting.Quota.Limits["dashboards"][1] = 4
				sc.m.Get("/dashboard", Quota("dashboards"), sc.defaultHandler)
				sc.fakeReq("GET", "/dashboard").exec()
				So(sc.resp.Code, ShouldEqual, 403)
			})
			Convey("org dashboard quota reached but quotas disabled", func() {
				setting.Quota.Limits["dashboards"][1] = 4
				setting.Quota.Enabled = false
				sc.m.Get("/dashboard", Quota("dashboard"), sc.defaultHandler)
				sc.fakeReq("GET", "/dashboard").exec()
				So(sc.resp.Code, ShouldEqual, 200)
			})

		})

	})
}
