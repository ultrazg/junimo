package backend

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SaveGamePathResultFlag struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

const gameExecFile = "horizon.exe" //

func (a *App) SaveGamePath() SaveGamePathResultFlag {
	options := runtime.OpenDialogOptions{
		Title: "选择游戏路径",
	}

	path, err := runtime.OpenDirectoryDialog(a.ctx, options)
	if err != nil {
		log.Printf("选择游戏路径失败: %v", err)
		return SaveGamePathResultFlag{
			Success: false,
			Message: fmt.Sprintf("选择游戏路径失败: %v", err),
			Path:    "",
		}
	}

	if path == "" {
		return SaveGamePathResultFlag{
			Success: false,
			Message: "未选择游戏路径",
			Path:    "",
		}
	}

	if !verifyGamePath(path) {
		return SaveGamePathResultFlag{
			Success: false,
			Message: "请选择正确的游戏路径",
			Path:    "",
		}
	}

	a.UpdateConfig("game_path", path)

	return SaveGamePathResultFlag{
		Success: true,
		Message: "",
		Path:    filepath.ToSlash(path),
	}
}

func verifyGamePath(path string) bool {
	gameExecFilePath := filepath.Join(path, gameExecFile)

	if _, err := os.Stat(gameExecFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func (a *App) LoadMods() {
	gamePath := a.ReadConfig("game_path")

	fmt.Println("game_path: ", gamePath)
}
