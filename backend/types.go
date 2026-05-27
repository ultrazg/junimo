package backend

import (
	"net/http"
	"time"
)

type SaveGamePathResultFlag struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

type SnackbarVariant string

type SnackbarColor string

const (
	SnackbarVariantSoft     = "soft"
	SnackbarVariantSolid    = "solid"
	SnackbarVariantOutlined = "outlined"
	SnackbarVariantPlain    = "plain"
)

const (
	SnackbarColorPrimary = "primary"
	SnackbarColorNeutral = "neutral"
	SnackbarColorDanger  = "danger"
	SnackbarColorSuccess = "success"
	SnackbarColorWarning = "warning"
)

type SnackbarShowOptions struct {
	Message          string          `json:"message"`
	ShowIcon         bool            `json:"showIcon"`
	AutoHideDuration int             `json:"autoHideDuration"`
	Color            SnackbarColor   `json:"color"`
	Variant          SnackbarVariant `json:"variant"`
}

type ModManifestJson struct {
	Name              string   `json:"name"`
	Author            string   `json:"author"`
	Version           string   `json:"version"`
	MinimumApiVersion string   `json:"minimumApiVersion"`
	Description       string   `json:"description"`
	UniqueID          string   `json:"uniqueID"`
	EntryDll          string   `json:"entryDll"`
	UpdateKeys        []string `json:"updateKeys"`
	NexusKey          int      `json:"nexusKey"`
	ManifestPath      string   `json:"manifestPath"`
	ModPath           string   `json:"modPath"`
	ConfigPath        string   `json:"configPath"`
}

type LoadModsOptions struct {
	Mods  []ModManifestJson `json:"mods"`
	Total int               `json:"total"`
}

type ImportModPreview struct {
	ZipPath      string          `json:"zipPath"`
	Manifest     ModManifestJson `json:"manifest"`
	Exists       bool            `json:"exists"`
	ExistingPath string          `json:"existingPath"`
	Error        string          `json:"error"`
}

type ConfirmImportModResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ModUpdateInfo struct {
	Name           string `json:"name"`
	NexusKey       int    `json:"nexusKey"`
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	HasUpdate      bool   `json:"hasUpdate"`
	Error          string `json:"error"`
}

type CheckForUpdatesResult struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Total   int             `json:"total"`
	Items   []ModUpdateInfo `json:"items"`
}

type ListBackupDirsResult struct {
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	CreateTime time.Time `json:"createTime"`
}

const (
	// 检查 API 密钥是否有效并返回用户的详细信息
	ApiUsersValidate = "https://api.nexusmods.com/v1/users/validate.json"

	// 根据 ModID 获取 Mod 的所有文件
	ApiViewSpecifiedModFile = "https://api.nexusmods.com/v1/games/stardewvalley/mods/%s/files.json"
)

type Client struct {
	http   *http.Client
	apiKey string
}

type Nexus struct {
}

type NexusUserValidateResult struct {
	Flag       bool              `json:"flag"`
	UserID     int               `json:"user_id"`
	Key        string            `json:"key"`
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	ProfileUrl string            `json:"profile_url"`
	RateLimit  NexusAPIRateLimit `json:"rate_limit"`
}

type NexusAPIRateLimit struct {
	HourlyLimit     string `json:"hourly_limit"`
	HourlyRemaining string `json:"hourly_remaining"`
	HourlyReset     string `json:"hourly_reset"`
	DailyLimit      string `json:"daily_limit"`
	DailyRemaining  string `json:"daily_remaining"`
	DailyReset      string `json:"daily_reset"`
}

type NexusViewSpecifiedModFileResult struct {
	Flag        bool          `json:"flag"`
	Files       []Files       `json:"files"`
	FileUpdates []FileUpdates `json:"file_updates"`
}

type Files struct {
	Id                   []int  `json:"id"`
	Uid                  int    `json:"uid"`
	FileID               int    `json:"file_id"`
	Name                 string `json:"name"`
	Version              string `json:"version"`
	CategoryID           int    `json:"category_id"`
	CategoryName         string `json:"category_name"`
	IsPrimary            bool   `json:"is_primary"`
	Size                 int64  `json:"size"`
	FileName             string `json:"file_name"`
	UploadedTimestamp    int64  `json:"uploaded_timestamp"`
	UploadedTime         string `json:"uploaded_time"`
	ModVersion           string `json:"mod_version"`
	ExternalVirusScanUrl string `json:"external_virus_scan_url"`
	Description          string `json:"description"`
	SizeKB               int64  `json:"size_kb"`
	SizeInBytes          int64  `json:"size_in_bytes"`
	ChangelogHtml        string `json:"changelog_html"`
	ContentPreviewLink   string `json:"content_preview_link"`
}

type FileUpdates struct {
	OldFileID         int    `json:"old_file_id"`
	NewFileID         int    `json:"new_file_id"`
	OldFileName       string `json:"old_file_name"`
	NewFileName       string `json:"new_file_name"`
	UploadedTimestamp int64  `json:"uploaded_timestamp"`
	UploadedTime      string `json:"uploaded_time"`
}
