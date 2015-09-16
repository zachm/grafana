package api

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/setting"
)

func AdminGetQuotas(c *middleware.Context) Response {
	if !setting.Quota.Enabled {
		return ApiError(404, "Quotas not enabled", nil)
	}

	query := m.GetQuotasQuery{
		UserId: c.QueryInt64("userId"),
		OrgId:  c.QueryInt64("orgId"),
	}

	if err := bus.Dispatch(&query); err != nil {
		return ApiError(500, "Failed to get org quotas", err)
	}

	result := make([]*dtos.QuotaInfo, 0)
	for _, quota := range query.Result {
		result = append(result, &dtos.QuotaInfo{
			Id:     quota.Id,
			Target: quota.Target,
			UserId: quota.UserId,
			OrgId:  quota.OrgId,
			Limit:  quota.Limit,
		})
	}

	return Json(200, result)
}

// POST /api/admin/quotas/
func AdminAddQuota(c *middleware.Context, cmd m.AddQuotaCommand) Response {
	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "Failed to updated org quota", err)
	}

	return ApiSuccess("Quota Added")
}

// DELETE /api/admin/quotas/
func AdminRemoveQuota(c *middleware.Context) Response {
	cmd := m.RemoveQuotaCommand{Id: c.ParamsInt64(":id")}

	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "Failed to updated org quota", err)
	}

	return ApiSuccess("Quota Removed")
}
