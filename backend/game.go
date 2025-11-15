package backend

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	log.Println("打开目录对话框")

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

	log.Printf("选择游戏路径: %s", path)

	a.UpdateConfig("game_path", path)

	err = os.Mkdir(filepath.Join(path, "JUNIMO_BACKUP"), fs.FileMode(0755))
	if err != nil {
		log.Printf("创建 backup 目录失败：%v \n", err)
	} else {
		log.Printf("创建 backup 目录成功: %s", filepath.Join(path, "JUNIMO_BACKUP"))

		a.UpdateConfig("backup_path", filepath.Join(path, "JUNIMO_BACKUP"))
	}

	err = os.Mkdir(filepath.Join(path, "JUNIMO_DISABLED"), fs.FileMode(0755))
	if err != nil {
		log.Printf("创建 disabled 目录失败：%v \n", err)
	} else {
		log.Printf("创建 disabled 目录成功: %s", filepath.Join(path, "JUNIMO_DISABLED"))

		a.UpdateConfig("disabled_path", filepath.Join(path, "JUNIMO_DISABLED"))
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

func (a *App) LoadEnabledMods(showSnackbar bool) {
	gamePath := a.ReadConfig("game_path")
	modsPath := filepath.Join(gamePath.(string), "Mods")
	SMAPIVersion := a.ReadConfig("smapi_version").(string)

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

	runtime.EventsEmit(a.ctx, "LoadEnabledMods", LoadModsOptions{
		Mods:  modsConfig,
		Total: len(modsConfig),
	})

	if SMAPIVersion != "" {
		runtime.WindowSetTitle(a.ctx, fmt.Sprintf("Junimo - SMAPI 版本：v%s", SMAPIVersion))
	}

	if showSnackbar {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("成功加载 %d 个 Mod", len(modsConfig)),
			ShowIcon: true,
			Color:    SnackbarColorSuccess,
			Variant:  SnackbarVariantSoft,
		})
	}
}

func (a *App) LoadDisabledMods() {
	gamePath := a.ReadConfig("game_path")
	disabledModsPath := filepath.Join(gamePath.(string), "JUNIMO_DISABLED")

	disabledMods, err := os.ReadDir(disabledModsPath)
	if err != nil {
		log.Printf("读取 JUNIMO_DISABLED 目录失败: %v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("无法读取已禁用 Mod 目录: %v", err),
			ShowIcon:         true,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
			AutoHideDuration: 6000,
		})

		return
	}

	var disabledModsConfig []ModManifestJson

	for _, mod := range disabledMods {
		if mod.IsDir() {
			modManifestPath := findModManifestFile(filepath.Join(disabledModsPath, mod.Name()))
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

				disabledModsConfig = append(disabledModsConfig, modManifest)
			}
		}
	}

	runtime.EventsEmit(a.ctx, "LoadDisabledMods", LoadModsOptions{
		Mods:  disabledModsConfig,
		Total: len(disabledModsConfig),
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

	log.Printf("开始解析文件 %s ", manifestFile)

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

func (a *App) BackupModDir() {
	gamePath := a.ReadConfig("game_path").(string)
	backupPath := a.ReadConfig("backup_path").(string)

	log.Printf("gamePath: %s", gamePath)
	log.Printf("backupPath: %s", backupPath)

	if gamePath == "" || backupPath == "" {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  "无法备份，请在设置中检查或配置游戏目录",
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	modDir := filepath.Join(gamePath, "Mods")

	if _, err := os.Stat(modDir); os.IsNotExist(err) {
		log.Printf("Mods 目录不存在: %v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  "Mods 目录不存在，无法备份",
			ShowIcon: true,
			Color:    SnackbarColorWarning,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	timeSuffix := time.Now().Format("20060102_150405")
	targetDir := filepath.Join(backupPath, "Mods_backup_"+timeSuffix)

	err := copyDir(modDir, targetDir)
	if err != nil {
		log.Printf("备份 Mods 目录失败: %v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("备份失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	log.Printf("Mods 目录已成功备份至 %s", targetDir)

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:          fmt.Sprintf("Mods 目录已成功备份至 %s", targetDir),
		ShowIcon:         true,
		Color:            SnackbarColorSuccess,
		Variant:          SnackbarVariantSoft,
		AutoHideDuration: 6000,
	})
}

func (a *App) ListBackupDirs() []ListBackupDirsResult {
	backupPath := a.ReadConfig("backup_path").(string)
	if backupPath == "" {
		log.Printf("备份路径 %s 为空", backupPath)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  "备份路径为空，请在设置中检查或配置游戏目录",
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

			log.Printf("读取的备份：%s", info.Name())

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

func (a *App) RestoreBackup(name string) {
	backPath := a.ReadConfig("backup_path").(string)
	gamePath := a.ReadConfig("game_path").(string)
	modsPath := filepath.Join(gamePath, "Mods")
	backupDir := filepath.Join(backPath, name)

	err := clearDir(modsPath)
	if err != nil {
		log.Printf("清空 Mods 目录失败: %v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("恢复备份失败: %v", err),
			ShowIcon:         true,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
			AutoHideDuration: 3000,
		})

		return
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		log.Printf("读取备份目录 %s 失败: %v", backupDir, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("恢复备份失败: %v", err),
			ShowIcon:         true,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
			AutoHideDuration: 3000,
		})

		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(backupDir, entry.Name())
		dstPath := filepath.Join(modsPath, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				log.Printf("恢复备份 %s 失败: %v", entry.Name(), err)

				SnackbarShow(a.ctx, &SnackbarShowOptions{
					Message:          fmt.Sprintf("恢复备份 %s 失败: %v", entry.Name(), err),
					ShowIcon:         true,
					Color:            SnackbarColorDanger,
					Variant:          SnackbarVariantSoft,
					AutoHideDuration: 3000,
				})

				continue
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				log.Printf("恢复备份 %s 失败: %v", entry.Name(), err)

				SnackbarShow(a.ctx, &SnackbarShowOptions{
					Message:          fmt.Sprintf("恢复备份 %s 失败: %v", entry.Name(), err),
					ShowIcon:         true,
					Color:            SnackbarColorDanger,
					Variant:          SnackbarVariantSoft,
					AutoHideDuration: 3000,
				})

				continue
			}
		}
	}

	log.Printf("成功恢复备份 %s", name)

	a.LoadEnabledMods(false)

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:          fmt.Sprintf("成功恢复备份 %s", name),
		ShowIcon:         true,
		Color:            SnackbarColorSuccess,
		Variant:          SnackbarVariantSoft,
		AutoHideDuration: 3000,
	})
}

func (a *App) DisableMod(path string) {
	disabledPath := a.ReadConfig("disabled_path").(string)

	log.Printf("disabledPath: %s", disabledPath)

	info, err := os.Stat(path)
	if err != nil {
		log.Printf("读取目录 %s 失败: %v", path, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("无法读取目录 %s，禁用失败: %v", path, err),
			ShowIcon:         true,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
			AutoHideDuration: 3000,
		})

		return
	}

	dstPath := filepath.Join(disabledPath, info.Name())

	err = MoveDir(path, dstPath)
	if err != nil {
		log.Printf("禁用 Mod %s 失败: %v", path, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("禁用 Mod 失败: %v", err),
			ShowIcon:         true,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
			AutoHideDuration: 3000,
		})

		return
	}

	a.LoadEnabledMods(false)
	a.LoadDisabledMods()

	log.Printf("成功禁用 Mod %s", path)

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:          "成功禁用 Mod",
		ShowIcon:         true,
		Color:            SnackbarColorSuccess,
		Variant:          SnackbarVariantSoft,
		AutoHideDuration: 3000,
	})
}

func (a *App) EnableMod(path string) {
	baseDisabledPath := a.ReadConfig("disabled_path").(string)
	gamePath := a.ReadConfig("game_path").(string)
	modsPath := filepath.Join(gamePath, "Mods")
	path = filepath.Clean(path)
	parts := strings.Split(path, string(os.PathSeparator))

	log.Printf("baseDisabledPath: %s", baseDisabledPath)
	log.Printf("gamePath: %s", gamePath)
	log.Printf("modsPath: %s", modsPath)
	log.Printf("path: %s", path)
	log.Printf("parts: %s", parts)

	var modBaseDir string
	for i, part := range parts {
		if part == "JUNIMO_DISABLED" {
			if i+1 < len(parts) {
				modBaseDir = parts[i+1]
			}
			break
		}
	}

	srcPath := filepath.Join(baseDisabledPath, modBaseDir)
	dstPath := filepath.Join(modsPath, modBaseDir)

	err := MoveDir(srcPath, dstPath)
	if err != nil {
		log.Printf("启用 Mod %s 失败: %v", srcPath, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("启用 Mod 失败: %v", err),
			ShowIcon:         true,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
			AutoHideDuration: 3000,
		})

		return
	}

	log.Printf("成功启用 Mod %s", srcPath)

	a.LoadEnabledMods(false)
	a.LoadDisabledMods()

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:          "成功启用 Mod",
		ShowIcon:         true,
		Color:            SnackbarColorSuccess,
		Variant:          SnackbarVariantSoft,
		AutoHideDuration: 3000,
	})
}
