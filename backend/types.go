package backend

type Config struct {
	Theme        string            `json:"theme"`
	GamePath     string            `json:"game_path"`
	Mods         []ModManifestJson `json:"mods"`
	SMAPIVersion string            `json:"smapi_version"`
}

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
}

type LoadModsOptions struct {
	Mods  []ModManifestJson `json:"mods"`
	Total int               `json:"total"`
}
