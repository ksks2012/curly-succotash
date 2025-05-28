package setting

import (
	"time"
)

type AppSettingS struct {
	RunMode     string
	LogSavePath string
	LogFileName string
	LogFileExt  string
}

type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         []string
	SocketPath   string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
	Path         string
}

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type StoragePathSettingS struct {
	PDFFoldar string
}

type AISettingS struct {
	APIKey          string
	Model           string
	Temperature     float32
	TopP            float32
	TopK            float32
	MaxOutputTokens int32
	Stream          bool
}

var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}

	return nil
}

func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
