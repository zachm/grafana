package api

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/metrics"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/util"
)

func AdminCreateUser(c *middleware.Context, form dtos.AdminCreateUserForm) Response {
	cmd := m.CreateUserCommand{
		Login:    form.Login,
		Email:    form.Email,
		Password: form.Password,
		OrgName:  form.OrgName,
		OrgRole:  form.OrgRole,
		NewOrg:   form.NewOrg,
		Name:     form.Name,
	}

	if len(cmd.Login) == 0 {
		cmd.Login = cmd.Email
		if len(cmd.Login) == 0 {
			return ApiError(400, "Validation error, need specify either username or email", nil)
		}
	}

	if len(cmd.Password) < 4 {
		return ApiError(400, "Password is missing or too short", nil)
	}

	if err := bus.Dispatch(&cmd); err != nil {
		return ApiError(500, "failed to create user", err)
	}

	metrics.M_Api_Admin_User_Create.Inc(1)

	return ApiSuccess("User created")
}

func AdminUpdateUserPassword(c *middleware.Context, form dtos.AdminUpdateUserPasswordForm) {
	userId := c.ParamsInt64(":id")

	if len(form.Password) < 4 {
		c.JsonApiErr(400, "New password too short", nil)
		return
	}

	userQuery := m.GetUserByIdQuery{Id: userId}

	if err := bus.Dispatch(&userQuery); err != nil {
		c.JsonApiErr(500, "Could not read user from database", err)
		return
	}

	passwordHashed := util.EncodePassword(form.Password, userQuery.Result.Salt)

	cmd := m.ChangeUserPasswordCommand{
		UserId:      userId,
		NewPassword: passwordHashed,
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to update user password", err)
		return
	}

	c.JsonOK("User password updated")
}

func AdminUpdateUserPermissions(c *middleware.Context, form dtos.AdminUpdateUserPermissionsForm) {
	userId := c.ParamsInt64(":id")

	cmd := m.UpdateUserPermissionsCommand{
		UserId:         userId,
		IsGrafanaAdmin: form.IsGrafanaAdmin,
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to update user permissions", err)
		return
	}

	c.JsonOK("User permissions updated")
}

func AdminDeleteUser(c *middleware.Context) {
	userId := c.ParamsInt64(":id")

	cmd := m.DeleteUserCommand{UserId: userId}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to delete user", err)
		return
	}

	c.JsonOK("User deleted")
}
