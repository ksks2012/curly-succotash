// example main function
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/sys/unix"
	"gopkg.in/natefinch/lumberjack.v2"

	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/model"
	"curly-succotash/backend/pkg/logger"
	"curly-succotash/backend/pkg/setting"
)

var (
	runMode string
	cfg     string
)

func init() {
	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}
	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	// err = updateDB()
	// if err != nil {
	// 	log.Fatalf("init.updateDB err: %v", err)
	// }
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	global.Logger.Infof(ctx, "db_lock")
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, unix.SIGTERM)

	cancel()
}

func setupFlag() error {
	flag.StringVar(&runMode, "mode", "", "running level (info, debug)")
	flag.StringVar(&cfg, "config", "etc/", "assgin the path of config file")
	flag.Parse()

	return nil
}

func setupSetting() error {
	s, err := setting.NewSetting(strings.Split(cfg, ",")...)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	// TODO: run mode

	return nil
}

// TODO: DB mirgration
// func updateDB() error {
// 	var err error
// 	updateDBSetup := &config.StorageSetup{}
// 	err = updateDBSetup.NewDBEngine(global.DatabaseSetting)
// 	if err != nil {
// 		return err
// 	}
// 	if err = updateDBSetup.Instance.Open(); nil != err {
// 		log.Fatalf("open storage connection failed: %v", err)
// 		return err
// 	}

// 	return nil
// }

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}
