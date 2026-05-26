package backend

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	nexus   *Nexus
	logFile *os.File
}

func NewApp() *App {
	return &App{
		nexus: NewNexus(),
	}
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
	runtime.EventsEmit(ctx, "snackbar:show", options)
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

func (a *App) ValidateUser(apiKey string) *NexusUserValidateResult {
	r, err := a.nexus.ValidateUser(apiKey)
	if err != nil {
		log.Printf("验证 Nexus Mods 用户失败：%v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          "验证 Nexus Mods 用户失败：" + err.Error(),
			ShowIcon:         true,
			AutoHideDuration: 6000,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
		})

		return &NexusUserValidateResult{
			Flag: false,
		}
	}

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message: fmt.Sprintf("验证 Nexus Mods 用户成功！\r\n已登录用户：%s\r\nNexus Mods API 速率限制：\r\n当前小时限制：%s，当前小时剩余：%s，当前小时重置时间：%s\r\n当前天限制：%s，当前天剩余：%s，当前天重置时间：%s",
			r.Name, r.RateLimit.HourlyLimit, r.RateLimit.HourlyRemaining, r.RateLimit.HourlyReset, r.RateLimit.DailyLimit, r.RateLimit.DailyRemaining, r.RateLimit.DailyReset),
		ShowIcon:         true,
		AutoHideDuration: 3000,
		Color:            SnackbarColorSuccess,
		Variant:          SnackbarVariantSoft,
	})

	log.Printf("验证 Nexus Mods 用户成功：%v", r)

	return &NexusUserValidateResult{
		Flag:       true,
		UserID:     r.UserID,
		Key:        r.Key,
		Name:       r.Name,
		Email:      r.Email,
		ProfileUrl: r.ProfileUrl,
	}
}
