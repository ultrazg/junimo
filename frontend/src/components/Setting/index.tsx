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
} from '@/utils'
import { Button, Radio, RadioGroup, Typography, useColorScheme } from '@mui/joy'
import styles from './index.module.scss'
import FolderTwoToneIcon from '@mui/icons-material/FolderTwoTone'
import BuildTwoToneIcon from '@mui/icons-material/BuildTwoTone'

type IProps = {
  open: boolean
  onClose: () => void
}

const THEME_LISTS = ['light', 'dark']

const Setting: React.FC<IProps> = ({ open, onClose }) => {
  const [currentTheme, setCurrentTheme] = useState<'light' | 'dark'>('light')
  const [gamePath, setGamePath] = useState<string>('')
  const { setMode } = useColorScheme()

  const init = () => {
    ReadConfig('theme').then((res: 'light' | 'dark') => {
      setCurrentTheme(res)
    })

    ReadConfig('game_path').then((res: string) => {
      setGamePath(res)
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

  const handleSaveGamePath = async () => {
    const { success, message, path } = await onSaveGamePath()

    if (success) {
      snackbar.show(`游戏所在路径已设置为: ${path}`, {
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
        </div>
      </Modal>
    </React.Fragment>
  )
}

export default Setting
