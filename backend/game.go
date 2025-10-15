package backend

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
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

	var modsConfig []ModManifestJson

	for _, mod := range mods {
		if mod.IsDir() {
			modManifestPath := findModManifestFile(filepath.Join(modsPath, mod.Name()))
			if modManifestPath != "" {
				modManifest, err := parseManifestFile(a, modManifestPath)
				if err != nil {
					log.Printf("解析ModManifest文件 %s 失败: %v", modManifestPath, err)

					SnackbarShow(a.ctx, &SnackbarShowOptions{
						Message:  fmt.Sprintf("解析ModManifest文件 %s 失败: %v", modManifestPath, err),
						ShowIcon: true,
						Color:    SnackbarColorDanger,
						Variant:  SnackbarVariantSoft,
					})

					continue
				}

				modsConfig = append(modsConfig, modManifest)
			}
		}
	}

	a.UpdateConfig("mods", modsConfig)

	runtime.EventsEmit(a.ctx, "loadMods", LoadModsOptions{
		Mods:  modsConfig,
		Total: len(modsConfig),
	})

	runtime.WindowSetTitle(a.ctx, fmt.Sprintf("Junimo - 已加载 %d 个 Mod - SMAPI 版本：v%s", len(modsConfig), a.ReadConfig("smapi_version")))

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:  fmt.Sprintf("成功加载 %d 个 Mod", len(modsConfig)),
		ShowIcon: true,
		Color:    SnackbarColorSuccess,
		Variant:  SnackbarVariantSoft,
	})
}

func findModManifestFile(path string) string {
	manifestFilePath := filepath.Join(path, ModManifestFileName)

	if _, err := os.Stat(manifestFilePath); err == nil {
		return path
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Printf("读取目录 %s 失败: %v", path, err)
		return ""
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			subDirPath := filepath.Join(path, dir.Name())

			if result := findModManifestFile(subDirPath); result != "" {
				return result
			}
		}
	}

	return ""
}

func parseManifestFile(a *App, path string) (ModManifestJson, error) {
	manifestFile := filepath.Join(path, ModManifestFileName)

	jsonData, err := os.ReadFile(manifestFile)
	if err != nil {
		log.Printf("读取文件 %s 失败: %v", manifestFile, err)
		return ModManifestJson{}, err
	}

	jsonStr := string(jsonData)

	modManifest := ModManifestJson{
		Name:              gjson.Get(jsonStr, "Name").String(),
		Author:            gjson.Get(jsonStr, "Author").String(),
		Version:           gjson.Get(jsonStr, "Version").String(),
		MinimumApiVersion: gjson.Get(jsonStr, "MinimumApiVersion").String(),
		Description:       gjson.Get(jsonStr, "Description").String(),
		UniqueID:          gjson.Get(jsonStr, "UniqueID").String(),
		EntryDll:          gjson.Get(jsonStr, "EntryDll").String(),
	}

	updateKeys := gjson.Get(jsonStr, "UpdateKeys")
	if updateKeys.Exists() && updateKeys.IsArray() {
		for _, v := range updateKeys.Array() {
			modManifest.UpdateKeys = append(modManifest.UpdateKeys, v.String())
		}
	}

	if gjson.Get(jsonStr, "Name").String() == "Console Commands" {
		a.UpdateConfig("smapi_version", gjson.Get(jsonStr, "Version").String())
	}

	return modManifest, nil
}
