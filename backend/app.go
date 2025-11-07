package backend

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	logFile *os.File
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	logFile, err := initLog()
	if err != nil {
		log.Printf("初始化日志失败: %v", err)
	}

	a.logFile = logFile

	log.Println("Startup")
}

func (a *App) DomReady(ctx context.Context) {
	gamePath := a.ReadConfig("game_path").(string)

	if gamePath != "" {
		a.LoadEnabledMods(true)
		a.LoadDisabledMods()
	} else {
		SnackbarShow(ctx, &SnackbarShowOptions{
			Message:          "请先在设置中选择游戏路径",
			ShowIcon:         true,
			AutoHideDuration: 6000,
			Color:            SnackbarColorWarning,
			Variant:          SnackbarVariantSoft,
		})
	}
}

func (a *App) Shutdown(ctx context.Context) {
	log.Println("Shutdown")

	if a.logFile != nil {
		a.logFile.Close()
	}
}

func SnackbarShow(ctx context.Context, options *SnackbarShowOptions) {
	runtime.EventsEmit(ctx, "snackbarShow", options)
}

func (a *App) LogDirSize() int64 {
	appPath, err := os.Getwd()
	if err != nil {
		log.Printf("获取应用目录失败: %v", err)
		return 0
	}

	logDir := filepath.Join(appPath, "logs")

	totalSize, err := CalcDirSize(logDir)
	if err != nil {
		log.Printf("计算日志目录 %s 大小失败: %v", logDir, err)
		return 0
	}

	return totalSize
}
