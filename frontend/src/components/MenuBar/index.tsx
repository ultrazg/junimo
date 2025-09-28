import React, { useState } from 'react'
import { Sheet, ButtonGroup, IconButton, Typography } from '@mui/joy'
import SettingsTwoToneIcon from '@mui/icons-material/SettingsTwoTone'
import InfoTwoToneIcon from '@mui/icons-material/InfoTwoTone'
import styles from './index.module.scss'
import { About, Setting } from '@/components'

const MenuBar = () => {
  const [aboutOpen, setAboutOpen] = useState<boolean>(false)
  const [settingOpen, setSettingOpen] = useState<boolean>(false)

  return (
    <React.Fragment>
      <Sheet className={styles['menu-bar-wrapper']}>
        <ButtonGroup
          variant={'plain'}
          spacing={1}
          color={'primary'}
        >
          <IconButton
            style={{ cursor: 'default' }}
            onClick={() => setSettingOpen(true)}
          >
            <SettingsTwoToneIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              设置
            </Typography>
          </IconButton>

          <IconButton
            style={{ cursor: 'default' }}
            onClick={() => setAboutOpen(true)}
          >
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
