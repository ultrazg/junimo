import React, { useState } from 'react'
import { Sheet, ButtonGroup, IconButton, Typography } from '@mui/joy'
import SettingsTwoToneIcon from '@mui/icons-material/SettingsTwoTone'
import InfoTwoToneIcon from '@mui/icons-material/InfoTwoTone'
import SyncTwoToneIcon from '@mui/icons-material/SyncTwoTone'
import ContentCopyOutlinedIcon from '@mui/icons-material/ContentCopyOutlined'
import SourceOutlinedIcon from '@mui/icons-material/SourceOutlined'
import styles from './index.module.scss'
import { About, Setting, BackupModal } from '@/components'
import { LoadEnabledMods, LoadDisabledMods, BackupModDir } from '@/utils'
import SMAPI_ICON from '@/assets/images/smapi_icon.png'

const MenuBar = () => {
  const [aboutOpen, setAboutOpen] = useState<boolean>(false)
  const [settingOpen, setSettingOpen] = useState<boolean>(false)
  const [syncLoading, setSyncLoading] = useState<boolean>(false)
  const [backupLoading, setBackupLoading] = useState<boolean>(false)
  const [backupModalOpen, setBackupModalOpen] = useState<boolean>(false)

  const onLoadMods = () => {
    setSyncLoading(true)

    Promise.allSettled([LoadEnabledMods(true), LoadDisabledMods()])
      .catch((err) => {
        console.error(err)
      })
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
          size={'sm'}
        >
          <IconButton title={'启动 SMAPI'}>
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

          <IconButton
            onClick={() => setSettingOpen(true)}
            title={'打开设置'}
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
            onClick={() => onLoadMods()}
            loading={syncLoading}
            title={'立即刷新 Mod 列表'}
          >
            <SyncTwoToneIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              刷新 Mod
            </Typography>
          </IconButton>

          <IconButton
            onClick={() => onBackupMod()}
            loading={backupLoading}
            title={'立即备份当前已启用的 Mod 的程序、资源和配置文件'}
          >
            <ContentCopyOutlinedIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              备份 Mod
            </Typography>
          </IconButton>

          <IconButton
            onClick={() => setBackupModalOpen(true)}
            title={'查看已备份的列表'}
          >
            <SourceOutlinedIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              查看备份
            </Typography>
          </IconButton>

          <IconButton
            onClick={() => setAboutOpen(true)}
            title={'关于 Junimo'}
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

      <BackupModal
        open={backupModalOpen}
        onClose={() => setBackupModalOpen(false)}
      />
    </React.Fragment>
  )
}

export default MenuBar
