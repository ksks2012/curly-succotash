// example main function
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
	"gopkg.in/natefinch/lumberjack.v2"

	"curly-succotash/backend/global"
	"curly-succotash/backend/internal/dao/config"
	"curly-succotash/backend/internal/model"
	"curly-succotash/backend/pkg/logger"
	"curly-succotash/backend/pkg/setting"
	"curly-succotash/backend/routers"
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
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
	err = updateDB()
	if err != nil {
		log.Fatalf("init.updateDB err: %v", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set Gin mode
	gin.SetMode(global.AppSetting.RunMode)

	// Initialize router
	router := routers.NewRouter()

	// Create HTTP server
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", global.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server in a goroutine
	go func() {
		global.Logger.Infof(ctx, "Starting server on port %d", global.ServerSetting.HttpPort)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Fatalf(ctx, "Failed to start server: %v", err)
		}
	}()

	// Setup signal handling
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, os.Interrupt, unix.SIGTERM)

	// Wait for shutdown signal
	<-stopChannel
	global.Logger.Infof(ctx, "Shutting down server")

	// Create a context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Shutdown server
	if err := s.Shutdown(shutdownCtx); err != nil {
		global.Logger.Errorf(ctx, "Server shutdown failed: %v", err)
	}

	global.Logger.Infof(ctx, "Server stopped")
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
	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("StoragePath", &global.StoragePathSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("AI", &global.AISetting)
	if err != nil {
		return err
	}

	// TODO: run mode

	return nil
}

// DB mirgration
func updateDB() error {
	var err error
	updateDBSetup := &config.StorageSetup{}
	err = updateDBSetup.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	if err = updateDBSetup.Instance.Open(); nil != err {
		log.Fatalf("open storage connection failed: %v", err)
		return err
	}

	return nil
}

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
