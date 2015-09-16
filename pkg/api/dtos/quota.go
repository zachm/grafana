package dtos

type QuotaInfo struct {
	Id     int64  `json:"id"`
	Target string `json:"target"`
	UserId int64  `json:"userId"`
	OrgId  int64  `json:"orgId"`
	Limit  int64  `json:"limit"`
}
