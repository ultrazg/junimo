package backend

import (
	"context"

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

func SnackbarShow(ctx context.Context, options *SnackbarShowOptions) {
	runtime.EventsEmit(ctx, "snackbarShow", options)
}
