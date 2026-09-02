package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

func NewNexus() *Nexus {
	return &Nexus{}
}

// ValidateUser 检查 API 密钥是否有效并返回用户的详细信息
func (n *Nexus) ValidateUser(apiKey string) (*NexusUserValidateResult, error) {
	client, err := NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	result := &NexusUserValidateResult{}
	err = client.GetJSON(ApiUsersValidate, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ViewSpecifiedModFile 根据 ModID 获取 Mod 的所有文件
func (n *Nexus) ViewSpecifiedModFile(modID string) (*NexusViewSpecifiedModFileResult, error) {
	apiKey := viper.GetString("nexus_api_key")
	if apiKey == "" {
		return nil, fmt.Errorf("请检查 Nexus Mods API Key 是否正确")
	}

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	result := &NexusViewSpecifiedModFileResult{}
	err = client.GetJSON(fmt.Sprintf(ApiViewSpecifiedModFile, modID), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetDownloadLink 获取指定文件的下载链接，按 short_name 排序后取第一个可用 URI
func (n *Nexus) GetDownloadLink(modID, fileID string) (string, error) {
	apiKey := viper.GetString("nexus_api_key")
	if apiKey == "" {
		return "", fmt.Errorf("请检查 Nexus Mods API Key 是否正确")
	}

	client, err := NewClient(apiKey)
	if err != nil {
		return "", err
	}

	body, err := client.GetRaw(fmt.Sprintf(ApiDownloadLink, modID, fileID))
	if err != nil {
		var statusErr *StatusError
		if errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusForbidden {
			return "", fmt.Errorf("Nexus 拒绝访问（403）：API 获取下载链接需要 Premium 账户；Premium 用户请到 nexusmods.com/users/myaccount?tab=api 重新生成 API Key 并勾选「Request Download Permission」下载权限")
		}
		return "", err
	}

	var links []NexusDownloadLink
	if err := json.Unmarshal(body, &links); err != nil {
		return "", fmt.Errorf("解析下载链接失败: %w", err)
	}

	for _, l := range links {
		if l.URI != "" {
			return l.URI, nil
		}
	}

	return "", fmt.Errorf("未找到可用的下载链接（非 Premium 用户需要通过 Nexus 网站获取下载链接）")
}

// ViewModChangelog 根据 ModID 获取更新日志，返回按版本号降序排列的条目
func (n *Nexus) ViewModChangelog(modID string) ([]ModChangelogEntry, error) {
	apiKey := viper.GetString("nexus_api_key")
	if apiKey == "" {
		return nil, fmt.Errorf("请检查 Nexus Mods API Key 是否正确")
	}

	client, err := NewClient(apiKey)
	if err != nil {
		return nil, err
	}

	body, err := client.GetRaw(fmt.Sprintf(ApiViewModChangelog, modID))
	if err != nil {
		return nil, err
	}

	raw := map[string][]string{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	entries := make([]ModChangelogEntry, 0, len(raw))
	for version, changes := range raw {
		entries = append(entries, ModChangelogEntry{Version: version, Changes: changes})
	}

	sort.Slice(entries, func(i, j int) bool {
		return compareVersion(entries[i].Version, entries[j].Version) > 0
	})

	return entries, nil
}

const nexusMaxConcurrency = 5

// CheckForUpdate 并发检查传入 mods 的 Nexus 最新版本，返回每项的版本对比结果
func (n *Nexus) CheckForUpdate(mods []ModManifestJson) []ModUpdateInfo {
	results := make([]ModUpdateInfo, len(mods))
	var wg sync.WaitGroup
	sem := make(chan struct{}, nexusMaxConcurrency)

	for i, mod := range mods {
		wg.Add(1)
		go func(i int, mod ModManifestJson) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			info := ModUpdateInfo{
				Name:           mod.Name,
				NexusKey:       mod.NexusKey,
				CurrentVersion: mod.Version,
				ModPath:        mod.ModPath,
				ConfigPath:     mod.ConfigPath,
				UniqueID:       mod.UniqueID,
			}

			res, err := n.ViewSpecifiedModFile(strconv.Itoa(mod.NexusKey))
			if err != nil {
				info.Error = err.Error()
				results[i] = info
				return
			}

			latest, fileID, fileName := latestPrimaryFile(res.Files)
			info.LatestVersion = latest
			info.LatestFileID = fileID
			info.LatestFileName = fileName
			info.HasUpdate = latest != "" && compareVersion(latest, mod.Version) > 0

			results[i] = info
		}(i, mod)
	}

	wg.Wait()
	return results
}

// latestPrimaryFile 取主文件（is_primary）；若无主文件，则取上传时间最新的非过期文件
func latestPrimaryFile(files []Files) (version string, fileID int, fileName string) {
	var primary *Files
	var newest *Files
	for i := range files {
		f := &files[i]
		if f.CategoryName == "ARCHIVED" || f.CategoryName == "OLD_VERSION" {
			continue
		}
		if f.IsPrimary {
			if primary == nil || compareVersion(modVersionOf(f), modVersionOf(primary)) > 0 {
				primary = f
			}
		}
		if newest == nil || f.UploadedTimestamp > newest.UploadedTimestamp {
			newest = f
		}
	}
	if primary != nil {
		return modVersionOf(primary), primary.FileID, primary.FileName
	}
	if newest != nil {
		return modVersionOf(newest), newest.FileID, newest.FileName
	}
	return "", 0, ""
}

func modVersionOf(f *Files) string {
	if f.ModVersion != "" {
		return f.ModVersion
	}
	return f.Version
}

// compareVersion 比较两个点分版本号，返回 1/0/-1。解析失败时回退到字符串比较。
func compareVersion(a, b string) int {
	a = strings.TrimPrefix(strings.TrimSpace(a), "v")
	b = strings.TrimPrefix(strings.TrimSpace(b), "v")
	if a == b {
		return 0
	}
	pa := strings.Split(a, ".")
	pb := strings.Split(b, ".")
	n := len(pa)
	if len(pb) > n {
		n = len(pb)
	}
	for i := 0; i < n; i++ {
		var sa, sb string
		if i < len(pa) {
			sa = pa[i]
		}
		if i < len(pb) {
			sb = pb[i]
		}
		ia, errA := strconv.Atoi(stripNonDigits(sa))
		ib, errB := strconv.Atoi(stripNonDigits(sb))
		if errA != nil || errB != nil {
			if sa > sb {
				return 1
			}
			if sa < sb {
				return -1
			}
			continue
		}
		if ia > ib {
			return 1
		}
		if ia < ib {
			return -1
		}
	}
	return 0
}

func stripNonDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "0"
	}
	return b.String()
}
