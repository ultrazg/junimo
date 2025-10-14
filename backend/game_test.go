package backend

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const gamePath = "D:\\steam\\steamapps\\common\\Stardew Valley"
const modManifestFileName = "manifest.json"
const modMainDir = "D:\\steam\\steamapps\\common\\Stardew Valley\\Mods\\LookupAnything"

func TestLoadMods(t *testing.T) {
	modsPath := filepath.Join(gamePath, "Mods")

	mods, err := os.ReadDir(modsPath)
	if err != nil {
		t.Errorf("读取Mods目录失败: %v", err)
		return
	}

	for _, mod := range mods {
		if mod.IsDir() {
			findModManifestFileTest(t, filepath.Join(modsPath, mod.Name()))
		}
	}
}

func findModManifestFileTest(t *testing.T, path string) {
	manifestFilePath := filepath.Join(path, modManifestFileName)
	if _, err := os.Stat(manifestFilePath); os.IsNotExist(err) {
		// t.Errorf("× 目录 %s 中不存在 %s 文件", path, modManifestFileName)

		dirs, err := os.ReadDir(path)
		if err != nil {
			t.Errorf("读取目录 %s 失败: %v", path, err)
			return
		}

		for _, dir := range dirs {
			if dir.IsDir() {
				findModManifestFileTest(t, filepath.Join(path, dir.Name()))
			}
		}

	} else {
		t.Logf("√ 目录 %s 中存在 %s 文件", path, modManifestFileName)
	}
}

func TestLoadModManifest(t *testing.T) {
	manifestFile := filepath.Join(modMainDir, modManifestFileName)

	var modManifest ModManifestJson

	jsonData, err := os.ReadFile(manifestFile)
	if err != nil {
		t.Errorf("读取 %s 文件失败: %v", manifestFile, err)

		return
	}

	jsonData = bytes.TrimPrefix(jsonData, []byte("\xef\xbb\xbf"))

	err = json.Unmarshal(jsonData, &modManifest)
	if err != nil {
		t.Errorf("解析 %s 文件失败: %v", manifestFile, err)

		return
	}

	t.Logf("成功加载 %s 文件", manifestFile)
	t.Logf("Mod Name: %s", modManifest.Name)
	t.Logf("Mod Author: %s", modManifest.Author)
	t.Logf("Mod Version: %s", modManifest.Version)
	t.Logf("Mod Description: %s", modManifest.Description)
}
