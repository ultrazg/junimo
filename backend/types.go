package backend

type Config struct {
	Theme    string `json:"theme"`
	GamePath string `json:"game_path"`
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
