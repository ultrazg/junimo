package backend

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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

func (a *App) RemoveModDir(path string) {
	fmt.Println("删除 Mod 目录:", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Mod 目录 %s 不存在\n", path)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("Mod 目录 %s 不存在，删除失败", path),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("删除 Mod 目录 %s 失败: %v\n", path, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("删除 Mod 失败：%v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	fmt.Printf("目录 %s 删除成功\n", path)

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:  "成功删除",
		ShowIcon: true,
		Color:    SnackbarColorSuccess,
		Variant:  SnackbarVariantSoft,
	})

	a.LoadMods(false)
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

func (a *App) UpdateModConfigFile(path, configStr string) {
	err := os.WriteFile(path, []byte(configStr), 0644)
	if err != nil {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("写入 Mod 配置文件失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:  "成功更新 Mod 配置文件",
		ShowIcon: true,
		Color:    SnackbarColorSuccess,
		Variant:  SnackbarVariantSoft,
	})
}

func FixDrivePath(p string) string {
	if len(p) >= 2 && p[1] == ':' && (len(p) == 2 || (len(p) > 2 && p[2] != '\\' && p[2] != '/')) {
		return p[:2] + string(filepath.Separator) + p[2:]
	}

	return p
}

func (a *App) BackupModDir() {
	gamePath := a.ReadConfig("game_path").(string)
	backupPath := a.ReadConfig("backup_path").(string)

	if gamePath == "" || backupPath == "" {
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  "无法备份，请先在设置中配置或检查游戏目录",
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	modDir := filepath.Join(gamePath, "Mods")

	if _, err := os.Stat(modDir); os.IsNotExist(err) {
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
		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("备份失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:          fmt.Sprintf("Mods 目录已成功备份至 %s", targetDir),
		ShowIcon:         true,
		Color:            SnackbarColorSuccess,
		Variant:          SnackbarVariantSoft,
		AutoHideDuration: 6000,
	})
}

func copyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}
