package backend

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const SMAPIExecFileName = "StardewModdingAPI.exe"
const ModManifestFileName = "manifest.json"
const ModConfigFileName = "config.json"

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

	err = os.Mkdir(filepath.Join(path, "junimo_backup"), fs.FileMode(0755))
	if err != nil {
		fmt.Printf("创建目录失败：%v \n", err)
	} else {
		a.UpdateConfig("backup_path", filepath.Join(path, "junimo_backup"))
	}

	err = os.Mkdir(filepath.Join(path, "junimo_disabled"), fs.FileMode(0755))
	if err != nil {
		fmt.Printf("创建目录失败：%v \n", err)
	} else {
		a.UpdateConfig("disabled_path", filepath.Join(path, "junimo_disabled"))
	}

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

func (a *App) LoadMods(showSnackbar bool) {
	gamePath := a.ReadConfig("game_path")
	modsPath := filepath.Join(gamePath.(string), "Mods")

	mods, err := os.ReadDir(modsPath)
	if err != nil {
		log.Printf("读取 Mods 目录失败: %v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("无法读取 Mods 目录，请先在设置中配置或检查游戏目录: %v", err),
			ShowIcon:         true,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
			AutoHideDuration: 6000,
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
					log.Printf("解析 ModManifest 文件 %s 失败: %v", modManifestPath, err)

					SnackbarShow(a.ctx, &SnackbarShowOptions{
						Message:  fmt.Sprintf("解析 ModManifest 文件 %s 失败: %v", modManifestPath, err),
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

	if showSnackbar {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("成功加载 %d 个 Mod", len(modsConfig)),
			ShowIcon: true,
			Color:    SnackbarColorSuccess,
			Variant:  SnackbarVariantSoft,
		})
	}
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
	var configPath string

	if _, err := os.Stat(filepath.Join(path, ModConfigFileName)); err == nil {
		configPath = filepath.Join(path, ModConfigFileName)
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
		ManifestPath:      path,
		ModPath:           trimToFirstSubdirUnderMods(path),
		ConfigPath:        configPath,
		Disabled:          false,
	}
	updateKeys := gjson.Get(jsonStr, "UpdateKeys")
	if updateKeys.Exists() && updateKeys.IsArray() {
		for _, v := range updateKeys.Array() {
			modManifest.UpdateKeys = append(modManifest.UpdateKeys, v.String())
		}
	}

	//
	if gjson.Get(jsonStr, "Name").String() == "Console Commands" {
		a.UpdateConfig("smapi_version", gjson.Get(jsonStr, "Version").String())
	}

	return modManifest, nil
}

func trimToFirstSubdirUnderMods(path string) string {
	cleanPath := filepath.Clean(path)
	parts := strings.Split(cleanPath, string(filepath.Separator))

	for i, part := range parts {
		if strings.EqualFold(part, "Mods") {
			if i+1 < len(parts) {
				return FixDrivePath(filepath.Join(parts[:i+2]...))
			}

			return FixDrivePath(filepath.Join(parts[:i+1]...))
		}
	}

	return cleanPath
}

func (a *App) ListBackupDirs() []ListBackupDirsResult {
	backupPath := a.ReadConfig("backup_path").(string)
	if backupPath == "" {
		log.Printf("备份路径为空")

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  "备份路径为空",
			ShowIcon: true,
			Color:    SnackbarColorWarning,
			Variant:  SnackbarVariantSoft,
		})

		return []ListBackupDirsResult{}
	}

	dirs, err := os.ReadDir(backupPath)
	if err != nil {
		log.Printf("无法读取备份目录：%v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("无法读取备份目录：%v", err),
			ShowIcon:         true,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
			AutoHideDuration: 6000,
		})

		return []ListBackupDirsResult{}
	}

	if len(dirs) == 0 {
		return []ListBackupDirsResult{}
	}

	var backupDirs []ListBackupDirsResult
	for _, dir := range dirs {
		if dir.IsDir() {
			info, err := dir.Info()
			if err != nil {
				continue
			}

			size, err := CalcDirSize(filepath.Join(backupPath, info.Name()))
			if err != nil {
				log.Printf("计算备份目录 %s 大小失败: %v", info.Name(), err)
				size = 0

				continue
			}

			fmt.Printf("备份名称: %s, 大小: %d bytes, 修改时间: %v, 是否目录: %v\n",
				info.Name(), size, info.ModTime(), info.IsDir())

			backupDirs = append(backupDirs, ListBackupDirsResult{
				Name:       info.Name(),
				Size:       size,
				CreateTime: info.ModTime(),
			})
		}
	}

	return backupDirs
}
