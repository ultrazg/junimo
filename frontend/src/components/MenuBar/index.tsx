import React, { useState } from 'react'
import {
  Sheet,
  ButtonGroup,
  IconButton,
  Typography,
  Menu,
  MenuItem,
} from '@mui/joy'
import SettingsTwoToneIcon from '@mui/icons-material/SettingsTwoTone'
import InfoTwoToneIcon from '@mui/icons-material/InfoTwoTone'
import SyncTwoToneIcon from '@mui/icons-material/SyncTwoTone'
import ContentCopyOutlinedIcon from '@mui/icons-material/ContentCopyOutlined'
import SourceOutlinedIcon from '@mui/icons-material/SourceOutlined'
import SearchTwoToneIcon from '@mui/icons-material/SearchTwoTone'
import AddTwoToneIcon from '@mui/icons-material/AddTwoTone'
import styles from './index.module.scss'
import {
  About,
  Setting,
  BackupModal,
  ImportModModal,
  ModUpdatesModal,
} from '@/components'
import {
  LoadEnabledMods,
  LoadDisabledMods,
  BackupModDir,
  ImportMod,
  CheckForUpdatesBySMAPI,
  CheckForUpdatesByNexus,
} from '@/utils'
import SMAPI_ICON from '@/assets/images/smapi_icon.png'
import { ImportModPreviewType, CheckForUpdatesResultType } from '@/types'

const MenuOptions = [
  {
    menu_name: '通过 SMAPI',
    menu_description: '通过 SMAPI 来检查 Mod 的更新，需要启动游戏',
  },
  {
    menu_name: '通过 Nexus Mods API',
    menu_description:
      '通过 Nexus Mods API 来检查 Mod 的更新，需要提供 Nexus Mods API Key',
  },
]

const MenuBar = () => {
  const actionRef = React.useRef<() => void>(null)
  const anchorRef = React.useRef<any>(null)
  const [open, setOpen] = React.useState(false)
  const [aboutOpen, setAboutOpen] = useState<boolean>(false)
  const [settingOpen, setSettingOpen] = useState<boolean>(false)
  const [syncLoading, setSyncLoading] = useState<boolean>(false)
  const [backupLoading, setBackupLoading] = useState<boolean>(false)
  const [backupModalOpen, setBackupModalOpen] = useState<boolean>(false)
  const [importLoading, setImportLoading] = useState<boolean>(false)
  const [importModal, setImportModal] = useState<{
    open: boolean
    items: ImportModPreviewType[]
  }>({ open: false, items: [] })
  const [updatesModal, setUpdatesModal] = useState<{
    open: boolean
    loading: boolean
    result?: CheckForUpdatesResultType
  }>({ open: false, loading: false })

  const onCheckUpdatesByNexus = () => {
    setUpdatesModal({ open: true, loading: true, result: undefined })
    CheckForUpdatesByNexus()
      .then((res) => {
        setUpdatesModal({ open: true, loading: false, result: res })
      })
      .catch((err) => {
        setUpdatesModal({
          open: true,
          loading: false,
          result: {
            success: false,
            message: String(err),
            total: 0,
            items: [],
          },
        })
      })
  }

  const handleMenuItemClick = (index: number) => {
    if (index === 0) {
      CheckForUpdatesBySMAPI().then()
    } else if (index === 1) {
      onCheckUpdatesByNexus()
    }
    setOpen(false)
  }

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

  const onImportMod = () => {
    setImportLoading(true)
    ImportMod()
      .then((items) => {
        if (!items || items.length === 0) return
        setImportModal({ open: true, items })
      })
      .finally(() => setImportLoading(false))
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
            title={'导入 ZIP 格式的 Mod 文件'}
            loading={importLoading}
            onClick={onImportMod}
          >
            <AddTwoToneIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              导入 Mod
            </Typography>
          </IconButton>

          <IconButton
            ref={anchorRef}
            title={'检查 Mods 是否有更新'}
            onMouseDown={() => {
              // @ts-ignore
              actionRef.current = () => setOpen(!open)
            }}
            onClick={() => {
              actionRef.current?.()
            }}
          >
            <SearchTwoToneIcon />
            <Typography
              className={styles['button-text']}
              level="title-md"
            >
              检查 Mods 更新
            </Typography>
          </IconButton>

          <Menu
            open={open}
            onClose={() => setOpen(false)}
            anchorEl={anchorRef.current}
          >
            {MenuOptions.map((option, index) => (
              <MenuItem
                key={option.menu_name}
                title={option.menu_description}
                onClick={() => handleMenuItemClick(index)}
              >
                {option.menu_name}
              </MenuItem>
            ))}
          </Menu>

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

      <ImportModModal
        open={importModal.open}
        items={importModal.items}
        onClose={() => setImportModal({ open: false, items: [] })}
        onFinished={() => {
          LoadEnabledMods(false).then()
          LoadDisabledMods().then()
        }}
      />

      <ModUpdatesModal
        open={updatesModal.open}
        loading={updatesModal.loading}
        result={updatesModal.result}
        onClose={() =>
          setUpdatesModal({ open: false, loading: false, result: undefined })
        }
        onRefresh={onCheckUpdatesByNexus}
      />
    </React.Fragment>
  )
}

export default MenuBar
