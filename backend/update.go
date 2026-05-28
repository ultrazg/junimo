package backend

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const ModUpdateBackupSubdir = "mod_updates"
const ModUpdateBackupMetaFile = ".junimo_backup.json"

// UpdateMod 通过 Nexus 下载并安装最新版 mod，先备份再替换。keepConfig 为 true 时
// 把旧版 config.json 拷贝到新版目录里。完成后删除下载的 zip。
func (a *App) UpdateMod(modPath string, nexusKey int, fileID int, keepConfig bool) UpdateModResult {
	if modPath == "" || nexusKey <= 0 || fileID <= 0 {
		return UpdateModResult{Success: false, Message: "缺少必要参数"}
	}

	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		return UpdateModResult{Success: false, Message: fmt.Sprintf("Mod 目录 %s 不存在", modPath)}
	}

	manifestDir := findModManifestFile(modPath)
	if manifestDir == "" {
		return UpdateModResult{Success: false, Message: "未找到原 Mod 的 manifest.json"}
	}

	oldManifest, err := parseManifestFile(a, manifestDir)
	if err != nil {
		return UpdateModResult{Success: false, Message: fmt.Sprintf("解析原 manifest.json 失败: %v", err)}
	}

	uri, err := a.nexus.GetDownloadLink(strconv.Itoa(nexusKey), strconv.Itoa(fileID))
	if err != nil {
		log.Printf("获取下载链接失败：%v", err)
		return UpdateModResult{Success: false, Message: fmt.Sprintf("获取下载链接失败：%v", err)}
	}

	tempDir, err := os.MkdirTemp("", "junimo_update_*")
	if err != nil {
		return UpdateModResult{Success: false, Message: fmt.Sprintf("创建临时目录失败: %v", err)}
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, fmt.Sprintf("mod_%d_%d.zip", nexusKey, fileID))
	if err := DownloadFile(uri, zipPath); err != nil {
		log.Printf("下载 mod 失败：%v", err)
		return UpdateModResult{Success: false, Message: fmt.Sprintf("下载 mod 失败：%v", err)}
	}

	extractDir := filepath.Join(tempDir, "extract")
	if err := Unzip(zipPath, extractDir); err != nil {
		return UpdateModResult{Success: false, Message: fmt.Sprintf("解压失败: %v", err)}
	}

	newManifestDir := findModManifestFile(extractDir)
	if newManifestDir == "" {
		return UpdateModResult{Success: false, Message: "下载的 zip 中未找到 manifest.json"}
	}

	newManifest, err := parseManifestFile(a, newManifestDir)
	if err != nil {
		return UpdateModResult{Success: false, Message: fmt.Sprintf("解析新版 manifest.json 失败: %v", err)}
	}

	srcPath := newManifestDir
	if newManifestDir != extractDir {
		if rel, err := filepath.Rel(extractDir, newManifestDir); err == nil {
			top := strings.SplitN(rel, string(filepath.Separator), 2)[0]
			srcPath = filepath.Join(extractDir, top)
		}
	}

	backupRoot, err := ensureModUpdateBackupRoot(a)
	if err != nil {
		return UpdateModResult{Success: false, Message: fmt.Sprintf("创建备份目录失败: %v", err)}
	}

	timeSuffix := time.Now().Format("20060102_150405")
	originalDirName := filepath.Base(modPath)
	backupName := fmt.Sprintf("%s_v%s_%s", sanitizeDirName(originalDirName), sanitizeDirName(oldManifest.Version), timeSuffix)
	backupDir := filepath.Join(backupRoot, backupName)

	if err := copyDir(modPath, backupDir); err != nil {
		log.Printf("备份旧版 mod 失败：%v", err)
		return UpdateModResult{Success: false, Message: fmt.Sprintf("备份旧版 mod 失败: %v", err)}
	}

	meta := ModUpdateBackup{
		Name:        backupName,
		ModName:     oldManifest.Name,
		ModPath:     modPath,
		OldVersion:  oldManifest.Version,
		NewVersion:  newManifest.Version,
		UniqueID:    oldManifest.UniqueID,
		KeptConfig:  keepConfig,
		OriginalDir: originalDirName,
		CreateTime:  time.Now(),
	}
	if err := writeBackupMeta(backupDir, meta); err != nil {
		log.Printf("写入备份元数据失败：%v", err)
	}

	var savedConfig []byte
	if keepConfig && oldManifest.ConfigPath != "" {
		if data, err := os.ReadFile(oldManifest.ConfigPath); err == nil {
			savedConfig = data
		} else {
			log.Printf("读取旧版 config.json 失败：%v", err)
		}
	}

	if err := os.RemoveAll(modPath); err != nil {
		log.Printf("删除旧版 mod 目录失败：%v", err)
		return UpdateModResult{Success: false, Message: fmt.Sprintf("删除旧版 mod 目录失败: %v", err)}
	}

	if err := copyDir(srcPath, modPath); err != nil {
		log.Printf("写入新版 mod 失败：%v", err)
		// 尝试回滚
		_ = os.RemoveAll(modPath)
		_ = copyDir(backupDir, modPath)
		return UpdateModResult{Success: false, Message: fmt.Sprintf("写入新版 mod 失败: %v，已回滚为旧版", err)}
	}

	if savedConfig != nil {
		newManifestDirInPlace := findModManifestFile(modPath)
		if newManifestDirInPlace != "" {
			configTarget := filepath.Join(newManifestDirInPlace, ModConfigFileName)
			if err := os.WriteFile(configTarget, savedConfig, 0644); err != nil {
				log.Printf("写回旧版 config.json 失败：%v", err)
			}
		}
	}

	a.LoadEnabledMods(false)

	log.Printf("成功更新 Mod %s -> v%s（备份：%s）", oldManifest.Name, newManifest.Version, backupName)

	return UpdateModResult{
		Success:    true,
		Message:    fmt.Sprintf("已更新 %s 到 v%s", oldManifest.Name, newManifest.Version),
		BackupName: backupName,
		NewModPath: modPath,
	}
}

// ListModUpdateBackups 列出所有 mod 更新备份
func (a *App) ListModUpdateBackups() []ModUpdateBackup {
	backupRoot, err := modUpdateBackupRoot(a)
	if err != nil || backupRoot == "" {
		return []ModUpdateBackup{}
	}

	entries, err := os.ReadDir(backupRoot)
	if err != nil {
		return []ModUpdateBackup{}
	}

	results := make([]ModUpdateBackup, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dir := filepath.Join(backupRoot, entry.Name())
		meta, err := readBackupMeta(dir)
		if err != nil {
			info, statErr := entry.Info()
			if statErr != nil {
				continue
			}
			meta = ModUpdateBackup{
				Name:        entry.Name(),
				OriginalDir: entry.Name(),
				CreateTime:  info.ModTime(),
			}
		}
		size, err := CalcDirSize(dir)
		if err == nil {
			meta.Size = size
		}
		results = append(results, meta)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CreateTime.After(results[j].CreateTime)
	})

	return results
}

// RollbackModUpdate 将指定备份恢复到原 mod 目录，覆盖当前内容
func (a *App) RollbackModUpdate(backupName string) RollbackModUpdateResult {
	if backupName == "" {
		return RollbackModUpdateResult{Success: false, Message: "缺少备份名称"}
	}

	backupRoot, err := modUpdateBackupRoot(a)
	if err != nil || backupRoot == "" {
		return RollbackModUpdateResult{Success: false, Message: "未配置备份目录"}
	}

	backupDir := filepath.Join(backupRoot, backupName)
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return RollbackModUpdateResult{Success: false, Message: "备份不存在"}
	}

	meta, err := readBackupMeta(backupDir)
	if err != nil || meta.ModPath == "" {
		gamePath, _ := a.ReadConfig("game_path").(string)
		if gamePath == "" {
			return RollbackModUpdateResult{Success: false, Message: "无法解析备份的目标路径，请手动恢复"}
		}
		dirName := strings.TrimSuffix(backupName, filepath.Ext(backupName))
		meta.ModPath = filepath.Join(gamePath, "Mods", dirName)
	}

	if _, err := os.Stat(meta.ModPath); err == nil {
		if err := os.RemoveAll(meta.ModPath); err != nil {
			return RollbackModUpdateResult{Success: false, Message: fmt.Sprintf("清理目标目录失败: %v", err)}
		}
	}

	if err := copyDirExcludingMeta(backupDir, meta.ModPath); err != nil {
		return RollbackModUpdateResult{Success: false, Message: fmt.Sprintf("恢复失败: %v", err)}
	}

	a.LoadEnabledMods(false)

	log.Printf("成功回滚 Mod 到备份 %s", backupName)

	return RollbackModUpdateResult{Success: true, Message: "已回滚到旧版"}
}

// RemoveModUpdateBackup 删除一个 mod 更新备份
func (a *App) RemoveModUpdateBackup(backupName string) RollbackModUpdateResult {
	if backupName == "" {
		return RollbackModUpdateResult{Success: false, Message: "缺少备份名称"}
	}
	backupRoot, err := modUpdateBackupRoot(a)
	if err != nil || backupRoot == "" {
		return RollbackModUpdateResult{Success: false, Message: "未配置备份目录"}
	}
	backupDir := filepath.Join(backupRoot, backupName)
	if err := os.RemoveAll(backupDir); err != nil {
		return RollbackModUpdateResult{Success: false, Message: fmt.Sprintf("删除备份失败: %v", err)}
	}
	return RollbackModUpdateResult{Success: true, Message: "已删除备份"}
}

func modUpdateBackupRoot(a *App) (string, error) {
	backupPath, _ := a.ReadConfig("backup_path").(string)
	if backupPath == "" {
		return "", fmt.Errorf("未配置备份目录")
	}
	return filepath.Join(backupPath, ModUpdateBackupSubdir), nil
}

func ensureModUpdateBackupRoot(a *App) (string, error) {
	root, err := modUpdateBackupRoot(a)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(root, fs.FileMode(0755)); err != nil {
		return "", err
	}
	return root, nil
}

func writeBackupMeta(backupDir string, meta ModUpdateBackup) error {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(backupDir, ModUpdateBackupMetaFile), data, 0644)
}

func readBackupMeta(backupDir string) (ModUpdateBackup, error) {
	var meta ModUpdateBackup
	data, err := os.ReadFile(filepath.Join(backupDir, ModUpdateBackupMetaFile))
	if err != nil {
		return meta, err
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return meta, err
	}
	return meta, nil
}

func copyDirExcludingMeta(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == ModUpdateBackupMetaFile {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}
