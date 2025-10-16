import React, { useState } from 'react'
import { Sheet, ButtonGroup, IconButton, Typography } from '@mui/joy'
import SettingsTwoToneIcon from '@mui/icons-material/SettingsTwoTone'
import InfoTwoToneIcon from '@mui/icons-material/InfoTwoTone'
import SyncTwoToneIcon from '@mui/icons-material/SyncTwoTone'
import styles from './index.module.scss'
import { About, Setting } from '@/components'
import { LoadMods } from '@/utils'
import SMAPI_ICON from '@/assets/images/smapi_icon.png'

const MenuBar = () => {
  const [aboutOpen, setAboutOpen] = useState<boolean>(false)
  const [settingOpen, setSettingOpen] = useState<boolean>(false)

  const onLoadMods = () => {
    LoadMods(true).then()
  }

  return (
    <React.Fragment>
      <Sheet className={styles['menu-bar-wrapper']}>
        <ButtonGroup
          variant={'plain'}
          spacing={1}
          color={'primary'}
        >
          <IconButton>
            <img
              src={SMAPI_ICON}
              alt={SMAPI_ICON}
            />

            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              启动游戏
            </Typography>
          </IconButton>

          <IconButton onClick={() => setSettingOpen(true)}>
            <SettingsTwoToneIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              设置
            </Typography>
          </IconButton>

          <IconButton onClick={() => onLoadMods()}>
            <SyncTwoToneIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              同步 Mod 列表
            </Typography>
          </IconButton>

          <IconButton onClick={() => setAboutOpen(true)}>
            <InfoTwoToneIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              关于...
            </Typography>
          </IconButton>
        </ButtonGroup>
      </Sheet>

      <About
        open={aboutOpen}
        onClose={() => setAboutOpen(false)}
      />

      <Setting
        open={settingOpen}
        onClose={() => setSettingOpen(false)}
      />
    </React.Fragment>
  )
}

export default MenuBar
