package models

import (
	"errors"
	"time"
)

var ErrInvalidQuotaTarget = errors.New("Invalid quota target")

type QuotaScope string

const (
	QuotaScopeGlobal QuotaScope = "global"
	QuotaScopeOrg    QuotaScope = "org"
	QuotaScopeUser   QuotaScope = "user"
)

type Quota struct {
	Id      int64
	OrgId   int64
	UserId  int64
	Target  string
	Limit   int64
	Created time.Time
	Updated time.Time
}

// type QuotaScope struct {
// 	Name         string
// 	Target       string
// 	DefaultLimit int64
// }
//
// type OrgQuotaDTO struct {
// 	OrgId  int64  `json:"org_id"`
// 	Target string `json:"target"`
// 	Limit  int64  `json:"limit"`
// 	Used   int64  `json:"used"`
// }

type QuotaDTO struct {
	OrgId  int64  `json:"org_id"`
	UserId int64  `json:"user_id"`
	Target string `json:"target"`
	Limit  int64  `json:"limit"`
	Used   int64  `json:"used"`
}

type GetQuotasByTargetQuery struct {
	Target string
	OrgId  int64
	UserId int64
	Result []*Quota
}

type IsQuotaReachedQuery struct {
	Target string
	OrgId  int64
	UserId int64
	Limit  int64

	Result bool
}

type SetQuotaCommand struct {
	Target string
	Limit  int64
	UserId int64
	OrgId  int64
}

// type UserQuotaDTO struct {
// 	UserId int64  `json:"user_id"`
// 	Target string `json:"target"`
// 	Limit  int64  `json:"limit"`
// 	Used   int64  `json:"used"`
// }
//
// type GlobalQuotaDTO struct {
// 	Target string `json:"target"`
// 	Limit  int64  `json:"limit"`
// 	Used   int64  `json:"used"`
// }
//
// type GetOrgQuotaByTargetQuery struct {
// 	Target  string
// 	OrgId   int64
// 	Default int64
// 	Result  *OrgQuotaDTO
// }
//
// type GetOrgQuotasQuery struct {
// 	OrgId  int64
// 	Result []*OrgQuotaDTO
// }
//
// type GetUserQuotaByTargetQuery struct {
// 	Target  string
// 	UserId  int64
// 	Default int64
// 	Result  *UserQuotaDTO
// }
//
// type GetQuotasByTargetQuery struct {
// 	Target string
// 	OrgId  int64
// 	UserId int64
// 	Result []*QuotaDTO
// }
//
// type GetUserQuotasQuery struct {
// 	UserId int64
// 	Result []*UserQuotaDTO
// }
//
// type GetGlobalQuotaByTargetQuery struct {
// 	Target  string
// 	Default int64
// 	Result  *GlobalQuotaDTO
// }
//
// type UpdateOrgQuotaCmd struct {
// 	Target string `json:"target"`
// 	Limit  int64  `json:"limit"`
// 	OrgId  int64  `json:"-"`
// }
//
// type UpdateUserQuotaCmd struct {
// 	Target string `json:"target"`
// 	Limit  int64  `json:"limit"`
// 	UserId int64  `json:"-"`
// }
//
// func GetQuotaScopes(target string) ([]QuotaScope, error) {
// 	scopes := make([]QuotaScope, 0)
// 	switch target {
// 	case "user":
// 		scopes = append(scopes,
// 			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.User},
// 			QuotaScope{Name: "org", Target: "org_user", DefaultLimit: setting.Quota.Org.User},
// 		)
// 		return scopes, nil
// 	case "org":
// 		scopes = append(scopes,
// 			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.Org},
// 			QuotaScope{Name: "user", Target: "org_user", DefaultLimit: setting.Quota.User.Org},
// 		)
// 		return scopes, nil
// 	case "dashboard":
// 		scopes = append(scopes,
// 			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.Dashboard},
// 			QuotaScope{Name: "org", Target: target, DefaultLimit: setting.Quota.Org.Dashboard},
// 		)
// 		return scopes, nil
// 	case "data_source":
// 		scopes = append(scopes,
// 			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.DataSource},
// 			QuotaScope{Name: "org", Target: target, DefaultLimit: setting.Quota.Org.DataSource},
// 		)
// 		return scopes, nil
// 	case "api_key":
// 		scopes = append(scopes,
// 			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.ApiKey},
// 			QuotaScope{Name: "org", Target: target, DefaultLimit: setting.Quota.Org.ApiKey},
// 		)
// 		return scopes, nil
// 	case "session":
// 		scopes = append(scopes,
// 			QuotaScope{Name: "global", Target: target, DefaultLimit: setting.Quota.Global.Session},
// 		)
// 		return scopes, nil
// 	default:
// 		return scopes, ErrInvalidQuotaTarget
// 	}
// }
