package global

import (
	"curly-succotash/backend/pkg/logger"
	"curly-succotash/backend/pkg/setting"
)

var (
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
	ServerSetting   *setting.ServerSettingS
)
