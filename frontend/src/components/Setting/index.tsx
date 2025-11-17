import React, { useEffect, useState } from 'react'
import { Modal } from '@/components'
import {
  UpdateConfig,
  setLightMode,
  setDarkMode,
  ReadConfig,
  snackbar,
  onSaveGamePath,
  OpenGameDir,
  OpenAppDir,
  LoadEnabledMods,
  OpenLogDir,
  LogDirSize,
  ValidateUser,
} from '@/utils'
import {
  Button,
  Radio,
  RadioGroup,
  Typography,
  useColorScheme,
  Switch,
  Input,
  Link,
  Avatar,
} from '@mui/joy'
import styles from './index.module.scss'
import FolderTwoToneIcon from '@mui/icons-material/FolderTwoTone'
import BuildTwoToneIcon from '@mui/icons-material/BuildTwoTone'
import VerifiedUserTwoToneIcon from '@mui/icons-material/VerifiedUserTwoTone'

type IProps = {
  open: boolean
  onClose: () => void
}

const THEME_LISTS = ['light', 'dark']

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  const value = bytes / Math.pow(k, i)
  return `${value.toFixed(2)} ${sizes[i]}`
}

const Setting: React.FC<IProps> = ({ open, onClose }) => {
  const [currentTheme, setCurrentTheme] = useState<'light' | 'dark'>('light')
  const [gamePath, setGamePath] = useState<string>('')
  const [autoCheckForUpdate, setAutoCheckForUpdate] = useState<boolean>(true)
  const [logDirSize, setLogDirSize] = useState<number>(0)
  const [nexusModsApiKey, setNexusModsApiKey] = useState<string>('')
  const [nexusUserInfo, setNexusUserInfo] = useState<{
    avatar: string
    name: string
  }>({
    avatar: '',
    name: '',
  })
  const [verifyLoading, setVerifyLoading] = useState<boolean>(false)
  const [isFocus, setIsFocus] = useState<boolean>(false)
  const { setMode } = useColorScheme()

  const init = () => {
    ReadConfig('theme').then((res: 'light' | 'dark') => {
      setCurrentTheme(res)
    })

    ReadConfig('game_path').then((res: string) => {
      setGamePath(res)
    })

    ReadConfig('nexus_user_avatar').then((res: string) => {
      setNexusUserInfo({
        ...nexusUserInfo,
        avatar: res,
      })
    })

    ReadConfig('nexus_user_name').then((res: string) => {
      setNexusUserInfo({
        ...nexusUserInfo,
        name: res,
      })
    })

    ReadConfig('nexus_api_key').then((res: string) => {
      setNexusModsApiKey(res)
    })

    ReadConfig('auto_check_for_update').then((res: boolean) => {
      setAutoCheckForUpdate(res)
    })

    LogDirSize().then((res: number) => {
      setLogDirSize(res)
    })
  }

  const handleChangeNexusModsApiKey = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const key = event.target.value
    setNexusModsApiKey(key)
    UpdateConfig('nexus_api_key', key).then()
  }

  const handleVerifyNexusModsApiKey = () => {
    setVerifyLoading(true)

    ValidateUser(nexusModsApiKey)
      .then((res) => {
        if (res.flag) {
          UpdateConfig('nexus_user_avatar', res.profile_url).then()
          UpdateConfig('nexus_user_name', res.name).then()
          setNexusUserInfo({
            name: res.name,
            avatar: res.profile_url,
          })
        }
      })
      .finally(() => {
        setVerifyLoading(false)
      })
  }

  const handleChangeTheme = (event: React.ChangeEvent<HTMLInputElement>) => {
    const theme = event.target.value
    setCurrentTheme(theme as 'light' | 'dark')
    setMode(theme === 'light' ? 'light' : 'dark')
    UpdateConfig('theme', theme).then()

    if (theme === 'light') {
      setLightMode()
    } else {
      setDarkMode()
    }
  }

  const handleChangeAutoCheckForUpdate = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const checked = event.target.checked
    setAutoCheckForUpdate(checked)
    UpdateConfig('auto_check_for_update', checked).then()
  }

  const handleSaveGamePath = async () => {
    const { success, message, path } = await onSaveGamePath()

    if (success) {
      setGamePath(path)

      LoadEnabledMods(true).then()

      snackbar.show(`游戏目录已设置为: ${path}`, {
        showIcon: true,
        color: 'success',
        variant: 'soft',
      })
    } else {
      snackbar.show(message, {
        showIcon: true,
        color: 'danger',
        variant: 'soft',
      })
    }
  }

  useEffect(() => {
    if (open) {
      init()
    }
  }, [open])

  return (
    <React.Fragment>
      <Modal
        open={open}
        onClose={onClose}
        title={'设置'}
      >
        <div className={styles['setting-wrapper']}>
          <div className={styles['setting-item']}>
            <div className={styles['label']}>
              <Typography level={'title-lg'}>主题</Typography>
            </div>
            <div className={styles['value']}>
              <RadioGroup
                value={currentTheme}
                orientation={'horizontal'}
              >
                {THEME_LISTS.map((item) => (
                  <Radio
                    onChange={handleChangeTheme}
                    key={item}
                    value={item}
                    label={
                      {
                        light: <Typography>浅色模式</Typography>,
                        dark: <Typography>深色模式</Typography>,
                      }[item]
                    }
                  />
                ))}
              </RadioGroup>
            </div>
          </div>

          <div className={styles['setting-item']}>
            <div className={styles['label']}>
              <Typography level={'title-lg'}>
                游戏目录
                <Typography level={'body-sm'}>
                  （Stardew Valley.exe 所在目录）
                </Typography>
              </Typography>
            </div>
            <div className={styles['value']}>
              {gamePath !== '' && (
                <Typography
                  level={'body-md'}
                  sx={{ mb: 1.5 }}
                  color={'neutral'}
                >
                  当前游戏目录：{gamePath}
                </Typography>
              )}
              <Button
                size={'sm'}
                variant={'soft'}
                onClick={handleSaveGamePath}
              >
                <BuildTwoToneIcon />
                {gamePath === '' ? '设置游戏目录' : '重新设置'}
              </Button>

              {gamePath !== '' && (
                <Button
                  size={'sm'}
                  variant={'soft'}
                  sx={{ ml: 1.5 }}
                  onClick={() => OpenGameDir()}
                >
                  <FolderTwoToneIcon />
                  打开游戏目录
                </Button>
              )}
            </div>
          </div>

          <div className={styles['setting-item']}>
            <div className={styles['label']}>
              <Typography level={'title-lg'}>应用程序目录</Typography>
            </div>
            <div className={styles['value']}>
              <Button
                size={'sm'}
                variant={'soft'}
                onClick={() => OpenAppDir()}
              >
                <FolderTwoToneIcon />
                打开应用程序目录
              </Button>
            </div>
          </div>

          <div className={styles['setting-item']}>
            <div className={styles['label']}>
              <Typography level={'title-lg'}>应用程序日志目录</Typography>
            </div>
            <div className={styles['value']}>
              {logDirSize !== 0 && (
                <Typography
                  level={'body-md'}
                  sx={{ mb: 1.5 }}
                  color={'neutral'}
                >
                  日志文件已占用 {formatBytes(logDirSize)}{' '}
                  的磁盘空间，可手动清除
                </Typography>
              )}

              <Button
                size={'sm'}
                variant={'soft'}
                onClick={() => OpenLogDir()}
              >
                <FolderTwoToneIcon />
                打开应用程序日志目录
              </Button>
            </div>
          </div>

          <div className={styles['setting-item']}>
            <div className={styles['label']}>
              <Typography level={'title-lg'}>启动时检查更新</Typography>
            </div>
            <div className={styles['value']}>
              <Switch
                checked={autoCheckForUpdate}
                onChange={handleChangeAutoCheckForUpdate}
              />
            </div>
          </div>

          <div className={styles['setting-item']}>
            <div className={styles['label']}>
              <Typography level={'title-lg'}>
                Nexus Mods API Key
                <Typography level={'body-sm'}>
                  （<Link>如何获取？</Link>）
                </Typography>
              </Typography>
            </div>
            <div className={styles['value']}>
              <Input
                className={isFocus ? '' : styles['blur-text']}
                size={'sm'}
                value={nexusModsApiKey}
                onChange={handleChangeNexusModsApiKey}
                placeholder={'请输入 Nexus Mods API Key'}
                onFocus={() => setIsFocus(true)}
                onBlur={() => setIsFocus(false)}
              />
              <div className={styles['nexus-user-info']}>
                <Button
                  size={'sm'}
                  variant={'soft'}
                  disabled={!nexusModsApiKey}
                  loading={verifyLoading}
                  onClick={handleVerifyNexusModsApiKey}
                >
                  <VerifiedUserTwoToneIcon />
                  验证
                </Button>
                {nexusUserInfo.name !== '' && (
                  <React.Fragment>
                    <Avatar
                      size={'sm'}
                      sx={{ ml: 1, mr: 1 }}
                      alt={'nexus user avatar'}
                      src={nexusUserInfo.avatar}
                    />
                    <span>{nexusUserInfo.name}</span>
                  </React.Fragment>
                )}
              </div>
            </div>
          </div>
        </div>
      </Modal>
    </React.Fragment>
  )
}

export default Setting
