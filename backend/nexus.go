package backend

import (
	"encoding/json"
	"fmt"
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
			}

			res, err := n.ViewSpecifiedModFile(strconv.Itoa(mod.NexusKey))
			if err != nil {
				info.Error = err.Error()
				results[i] = info
				return
			}

			latest := latestVersionFromFiles(res.Files)
			info.LatestVersion = latest
			info.HasUpdate = latest != "" && compareVersion(latest, mod.Version) > 0

			results[i] = info
		}(i, mod)
	}

	wg.Wait()
	return results
}

// latestVersionFromFiles 从 Nexus 返回的文件列表里取主文件版本，否则取上传时间最大的一项
func latestVersionFromFiles(files []Files) string {
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
		return modVersionOf(primary)
	}
	if newest != nil {
		return modVersionOf(newest)
	}
	return ""
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
