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

const (
	ApiUsersValidate = "https://api.nexusmods.com/v1/users/validate.json"
)

type Client struct {
	http   *http.Client
	apiKey string
}

type Nexus struct {
}

type NexusUserValidateResult struct {
	Flag       bool   `json:"flag"`
	UserID     int    `json:"user_id"`
	Key        string `json:"key"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfileUrl string `json:"profile_url"`
}
