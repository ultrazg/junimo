package backend

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

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

func (a *App) CheckForUpdatesByNexus() CheckForUpdatesResult {
	apiKey, _ := a.ReadConfig("nexus_api_key").(string)
	if apiKey == "" {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          "请先在设置中配置 Nexus Mods API Key",
			ShowIcon:         true,
			AutoHideDuration: 6000,
			Color:            SnackbarColorWarning,
			Variant:          SnackbarVariantSoft,
		})
		return CheckForUpdatesResult{Success: false, Message: "未配置 Nexus Mods API Key"}
	}

	gamePath, _ := a.ReadConfig("game_path").(string)
	if gamePath == "" {
		return CheckForUpdatesResult{Success: false, Message: "未配置游戏目录"}
	}

	modsPath := filepath.Join(gamePath, "Mods")
	entries, err := os.ReadDir(modsPath)
	if err != nil {
		log.Printf("读取 Mods 目录失败: %v", err)
		return CheckForUpdatesResult{Success: false, Message: fmt.Sprintf("读取 Mods 目录失败: %v", err)}
	}

	var targets []ModManifestJson
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		manifestDir := findModManifestFile(filepath.Join(modsPath, entry.Name()))
		if manifestDir == "" {
			continue
		}
		m, err := parseManifestFile(a, manifestDir)
		if err != nil {
			continue
		}
		if m.NexusKey > 0 {
			targets = append(targets, m)
		}
	}

	if len(targets) == 0 {
		return CheckForUpdatesResult{Success: true, Total: 0, Items: []ModUpdateInfo{}}
	}

	items := a.nexus.CheckForUpdate(targets)
	log.Printf("通过 Nexus 检查 %d 个 mod 更新完成", len(items))

	return CheckForUpdatesResult{
		Success: true,
		Total:   len(items),
		Items:   items,
	}
}

func (a *App) ViewModChangelog(nexusKey int) ViewModChangelogResult {
	entries, err := a.nexus.ViewModChangelog(strconv.Itoa(nexusKey))
	if err != nil {
		log.Printf("获取 Mod %d 更新日志失败：%v", nexusKey, err)
		return ViewModChangelogResult{Success: false, Message: err.Error()}
	}

	return ViewModChangelogResult{
		Success: true,
		Entries: entries,
	}
}
