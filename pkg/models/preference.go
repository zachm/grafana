package models

import "time"

const (
	PREF_KEY_HOME_DASHBOARD = "home.dashboard.id"
	PREF_KEY_THEME          = "home.theme"
)

type PreferenceLevel int

var (
	PREF_LEVEL_ORG         PreferenceLevel = 1
	PREF_LEVEL_ORG_OR_USER PreferenceLevel = 2
	PREF_LEVEL_USER        PreferenceLevel = 3
)

type PreferenceDef struct {
	Key   string
	Level PreferenceLevel
}

var PreferenceDefinitions = make(map[string]*PreferenceDef)

func init() {
	PreferenceDefinitions[PREF_KEY_HOME_DASHBOARD] = &PreferenceDef{
		Level: PREF_LEVEL_ORG_OR_USER,
	}
}

type Preference struct {
	Id     int64
	OrgId  int64
	UserId int64

	Key   string
	Value string

	Created time.Time
	Updated time.Time
}

// type Preferences struct {
// 	Id     int64
// 	OrgId  int64
// 	UserId int64
//
// 	HomeDashboardId int64
// 	Theme           string
// 	Timezone        string
//
// 	Created time.Time
// 	Updated time.Time
// }
//
// type SavePreferencesCommand struct {
// 	HomeDashboardId int64
// 	Theme           string
// 	Timezone        string
//
// 	OrgId  int64
// 	UserId int64
// }
//
// type GetPreferencesQuery struct {
// 	OrgId  int64
// 	UserId int64
// 	Merge  bool
//
// 	Result *Preferences
// }
