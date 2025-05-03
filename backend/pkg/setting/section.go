package setting

import (
	"time"
)

type AppSettingS struct {
	RunMode      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LogSavePath string
	LogFileName string
	LogFileExt string
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
