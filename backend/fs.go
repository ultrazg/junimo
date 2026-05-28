package backend

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func openDir(path string) error {
	log.Printf("打开目录: %s", path)

	cmd := exec.Command("explorer", path)

	return cmd.Start()
}

func (a *App) OpenGameDir() {
	gamePath := a.ReadConfig("game_path")

	err := openDir(gamePath.(string))
	if err != nil {
		log.Printf("打开游戏目录 %s 失败: %v", gamePath, err)

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
		log.Printf("打开应用目录 %s 失败: %v", appPath, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("打开应用目录失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})
	}
}

func (a *App) OpenLogDir() {
	appPath, _ := os.Getwd()
	logDir := filepath.Join(appPath, "logs")

	err := openDir(logDir)
	if err != nil {
		log.Printf("打开日志目录 %s 失败: %v", logDir, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("打开日志目录 %s 失败: %v", logDir, err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})
	}
}

func (a *App) OpenModDir(path string) {
	err := openDir(path)
	if err != nil {
		log.Printf("打开 Mod 目录 %s 失败: %v", path, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("打开 Mod 目录失败: %v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})
	}
}

func (a *App) OpenBackupDir(backupName string) {
	backupPath := a.ReadConfig("backup_path").(string)
	backupPath = filepath.Join(backupPath, backupName)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		log.Printf("备份目录 %s 不存在", backupPath)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("备份目录 %s 不存在，打开失败", backupPath),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	err := openDir(backupPath)
	if err != nil {
		log.Printf("打开备份目录 %s 失败: %v", backupPath, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("打开备份目录 %s 失败: %v", backupPath, err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}
}

func (a *App) RemoveModDir(path string) {
	log.Println("删除 Mod 目录:", path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Mod 目录 %s 不存在", path)

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
		log.Printf("删除 Mod 目录 %s 失败: %v", path, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("删除 Mod 失败：%v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	log.Printf("目录 %s 删除成功", path)

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:  "成功删除",
		ShowIcon: true,
		Color:    SnackbarColorSuccess,
		Variant:  SnackbarVariantSoft,
	})

	a.LoadEnabledMods(false)
}

func (a *App) RemoveBackupDir(path string) {
	backupPath := a.ReadConfig("backup_path").(string)
	removeBackupPath := filepath.Join(backupPath, path)

	log.Println("删除备份目录:", removeBackupPath)

	if _, err := os.Stat(removeBackupPath); os.IsNotExist(err) {
		log.Printf("备份目录 %s 不存在", removeBackupPath)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("备份目录 %s 不存在，删除失败", removeBackupPath),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	err := os.RemoveAll(removeBackupPath)
	if err != nil {
		log.Printf("删除备份目录 %s 失败: %v", removeBackupPath, err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:  fmt.Sprintf("删除备份失败：%v", err),
			ShowIcon: true,
			Color:    SnackbarColorDanger,
			Variant:  SnackbarVariantSoft,
		})

		return
	}

	log.Printf("备份 %s 删除成功", path)

	SnackbarShow(a.ctx, &SnackbarShowOptions{
		Message:  "成功删除",
		ShowIcon: true,
		Color:    SnackbarColorSuccess,
		Variant:  SnackbarVariantSoft,
	})
}

func (a *App) ReadModConfigFile(path string) string {
	configStr, err := os.ReadFile(path)
	if err != nil {
		log.Printf("读取 Mod 配置文件 %s 失败: %v", path, err)

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
		log.Printf("写入 Mod 配置文件 %s 失败: %v", path, err)

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

func copyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		log.Print(err)

		return err
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		log.Print(err)

		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		log.Print(err)

		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				log.Print(err)

				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				log.Print(err)

				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Print(err)

		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		log.Print(err)

		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		log.Print(err)

		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Print(err)

		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

func CalcDirSize(path string) (int64, error) {
	var totalSize int64

	err := filepath.WalkDir(path, func(entryPath string, d os.DirEntry, err error) error {
		if err != nil {
			log.Print(err)

			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				log.Print(err)

				return err
			}
			totalSize += info.Size()
		}

		return nil
	})

	return totalSize, err
}

func MoveDir(src, dst string) error {
	if err := copyDir(src, dst); err != nil {
		log.Print(src, dst, err)

		return err
	}

	return os.RemoveAll(src)
}

func clearDir(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Print(err)

		return err
	}

	for _, entry := range entries {
		path := filepath.Join(path, entry.Name())
		err = os.RemoveAll(path)
		if err != nil {
			log.Print(err)

			return err
		}
	}

	log.Printf("清空目录 %s 成功", path)

	return nil
}

func Unzip(zipPath, dstPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	dstPath, err = filepath.Abs(dstPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	for _, file := range reader.File {
		err := extractFile(file, dstPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractFile(file *zip.File, dstPath string) error {
	targetPath := filepath.Join(dstPath, file.Name)

	if !strings.HasPrefix(filepath.Clean(targetPath)+string(os.PathSeparator), dstPath+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path (zip slip?): %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open zipped file %s: %w", file.Name, err)
	}
	defer rc.Close()

	wc, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", targetPath, err)
	}

	_, err = io.Copy(wc, rc)
	closeErr := wc.Close()
	if closeErr != nil {
		log.Printf("warning: failed to close file %s: %v", targetPath, closeErr)
	}

	if err != nil {
		return fmt.Errorf("failed to copy data to file %s: %w", targetPath, err)
	}

	return nil
}

// DownloadFile 下载远程文件到 dstPath，使用更长的超时以兼容大文件
func DownloadFile(url, dstPath string) error {
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: HTTP %d %s", resp.StatusCode, resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("创建下载文件失败: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("写入下载文件失败: %w", err)
	}

	return nil
}
