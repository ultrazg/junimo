package backend

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func SnackbarShow(ctx context.Context, options *SnackbarShowOptions) {
	runtime.EventsEmit(ctx, "snackbarShow", options)
}

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
