package models

import "time"

type Quota struct {
	Id      int64
	OrgId   int64
	UserId  int64
	Target  string
	Limit   int64
	Created time.Time
	Updated time.Time
}

type GetQuotasQuery struct {
	OrgId  int64
	UserId int64

	Result []*Quota
}

type GetQuotaInfoQuery struct {
	Target string
	OrgId  int64
	UserId int64

	ResultLimit int64
	ResultUsed  int64
}

type AddQuotaCommand struct {
	Target string `json:"target"`
	Limit  int64  `json:"limit"`
	UserId int64  `json:"userId"`
	OrgId  int64  `json:"orgId"`

	Result *Quota
}

type RemoveQuotaCommand struct {
	Id int64
}
