package backend

import "time"

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
	ManifestPath      string   `json:"manifestPath"`
	ModPath           string   `json:"modPath"`
	ConfigPath        string   `json:"configPath"`
}

type LoadModsOptions struct {
	Mods  []ModManifestJson `json:"mods"`
	Total int               `json:"total"`
}

type ListBackupDirsResult struct {
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	CreateTime time.Time `json:"createTime"`
}
