import React, { useState } from 'react'
import { Sheet, ButtonGroup, IconButton, Typography } from '@mui/joy'
import SettingsTwoToneIcon from '@mui/icons-material/SettingsTwoTone'
import InfoTwoToneIcon from '@mui/icons-material/InfoTwoTone'
import SyncTwoToneIcon from '@mui/icons-material/SyncTwoTone'
import ContentCopyOutlinedIcon from '@mui/icons-material/ContentCopyOutlined'
import SourceOutlinedIcon from '@mui/icons-material/SourceOutlined'
import styles from './index.module.scss'
import { About, Setting, BackupModal } from '@/components'
import { LoadMods, BackupModDir } from '@/utils'
import SMAPI_ICON from '@/assets/images/smapi_icon.png'

const MenuBar = () => {
  const [aboutOpen, setAboutOpen] = useState<boolean>(false)
  const [settingOpen, setSettingOpen] = useState<boolean>(false)
  const [syncLoading, setSyncLoading] = useState<boolean>(false)
  const [backupLoading, setBackupLoading] = useState<boolean>(false)
  const [backupModalOpen, setBackupModalOpen] = useState<boolean>(false)

  const onLoadMods = () => {
    setSyncLoading(true)
    LoadMods(true)
      .then()
      .finally(() => {
        setSyncLoading(false)
      })
  }

  const onBackupMod = () => {
    setBackupLoading(true)
    BackupModDir()
      .then()
      .finally(() => {
        setBackupLoading(false)
      })
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
              alt={'SMAPI_ICON'}
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

          <IconButton
            onClick={() => onLoadMods()}
            loading={syncLoading}
          >
            <SyncTwoToneIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              同步 Mod 列表
            </Typography>
          </IconButton>

          <IconButton
            onClick={() => onBackupMod()}
            loading={backupLoading}
          >
            <ContentCopyOutlinedIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              备份 Mod 文件
            </Typography>
          </IconButton>

          <IconButton
            onClick={() => setBackupModalOpen(true)}
          >
            <SourceOutlinedIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              查看备份
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

      <BackupModal
        open={backupModalOpen}
        onClose={() => setBackupModalOpen(false)}
      />
    </React.Fragment>
  )
}

export default MenuBar
