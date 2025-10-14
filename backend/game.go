package backend

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const SMAPIExecFileName = "StardewModdingAPI.exe"
const ModManifestFileName = "manifest.json"

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
	gameExecFilePath := filepath.Join(path, SMAPIExecFileName)

	if _, err := os.Stat(gameExecFilePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func (a *App) LoadMods() {
	gamePath := a.ReadConfig("game_path")
	modsPath := filepath.Join(gamePath.(string), "Mods")

	mods, err := os.ReadDir(modsPath)
	if err != nil {
		log.Printf("读取Mods目录失败: %v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("读取Mods目录失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	for _, mod := range mods {
		if mod.IsDir() {
			findModManifestFile(filepath.Join(modsPath, mod.Name()))
		}
	}
}

func findModManifestFile(path string) {
	manifestFilePath := filepath.Join(path, ModManifestFileName)
	if _, err := os.Stat(manifestFilePath); os.IsNotExist(err) {

		dirs, err := os.ReadDir(path)
		if err != nil {
			log.Printf("读取目录 %s 失败: %v", path, err)

			return
		}

		for _, dir := range dirs {

			if dir.IsDir() {
				findModManifestFile(filepath.Join(path, dir.Name()))
			}

		}
	} else {
		// ...
		fmt.Printf("√ 目录 %s 中存在 %s 文件\n", path, ModManifestFileName)
	}
}
