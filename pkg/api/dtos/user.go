package dtos

import m "github.com/grafana/grafana/pkg/models"

type SignUpForm struct {
	Email    string `json:"email" binding:"Required"`
	Password string `json:"password" binding:"Required"`
}

type AdminCreateUserForm struct {
	Email    string     `json:"email"`
	Login    string     `json:"login"`
	Name     string     `json:"name"`
	OrgName  string     `json:"orgName"`
	OrgRole  m.RoleType `json:"orgRole"`
	NewOrg   bool       `json:"newOrg"`
	Password string     `json:"password" binding:"Required"`
}

type AdminUpdateUserForm struct {
	Email string `json:"email"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type AdminUpdateUserPasswordForm struct {
	Password string `json:"password" binding:"Required"`
}

type AdminUpdateUserPermissionsForm struct {
	IsGrafanaAdmin bool `json:"IsGrafanaAdmin"`
}

type AdminUserListItem struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	Login          string `json:"login"`
	IsGrafanaAdmin bool   `json:"isGrafanaAdmin"`
}

type SendResetPasswordEmailForm struct {
	UserOrEmail string `json:"userOrEmail" binding:"Required"`
}

type ResetUserPasswordForm struct {
	Code            string `json:"code"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}
