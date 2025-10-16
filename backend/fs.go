package backend

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func openDir(path string) error {
	cmd := exec.Command("explorer", path)

	return cmd.Start()
}

func (a *App) OpenGameDir() {
	gamePath := a.ReadConfig("game_path")

	err := openDir(gamePath.(string))
	if err != nil {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("打开游戏目录失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})
	}
}

func (a *App) OpenAppDir() {
	appPath, _ := os.Getwd()

	err := openDir(appPath)
	if err != nil {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("打开应用目录失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})
	}
}

func (a *App) OpenModDir(path string) {
	err := openDir(path)
	if err != nil {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("打开 Mod 目录失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})
	}
}

func (a *App) ReadModConfigFile(path string) string {
	configStr, err := os.ReadFile(path)
	if err != nil {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("读取 Mod 配置文件失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})
	}

	return string(configStr)
}

func FixDrivePath(p string) string {
	if len(p) >= 2 && p[1] == ':' && (len(p) == 2 || (len(p) > 2 && p[2] != '\\' && p[2] != '/')) {
		return p[:2] + string(filepath.Separator) + p[2:]
	}

	return p
}
