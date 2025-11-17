package backend

import (
	"fmt"
	"log"
)

func (a *App) ViewSpecifiedModFile(modID string) *NexusViewSpecifiedModFileResult {
	r, err := a.nexus.ViewSpecifiedModFile(modID)
	if err != nil {
		log.Printf("查看 Mod 文件失败：%v", err)

		SnackbarShow(a.ctx, &SnackbarShowOptions{
			Message:          fmt.Sprintf("查看 Mod 文件失败：%v", err),
			ShowIcon:         true,
			AutoHideDuration: 6000,
			Color:            SnackbarColorDanger,
			Variant:          SnackbarVariantSoft,
		})

		return &NexusViewSpecifiedModFileResult{
			Flag: false,
		}
	}

	log.Println("查看 Mod 文件成功")

	return &NexusViewSpecifiedModFileResult{
		Flag:        true,
		Files:       r.Files,
		FileUpdates: r.FileUpdates,
	}
}
