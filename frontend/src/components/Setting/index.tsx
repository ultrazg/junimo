import React, { useEffect, useState } from 'react'
import { Modal } from '@/components'
import { UpdateConfig, setLightMode, setDarkMode, getStyleMode } from '@/utils'
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
  const [currentTheme, setCurrentTheme] = useState<string>('light')
  const { setMode } = useColorScheme()

  const init = async () => {
    const theme = await getStyleMode()

    setCurrentTheme(theme)
  }

  const handleChangeTheme = (event: React.ChangeEvent<HTMLInputElement>) => {
    const theme = event.target.value
    setCurrentTheme(theme)
    setMode(theme === 'light' ? 'light' : 'dark')
    UpdateConfig('theme', theme).then()

    if (theme === 'light') {
      setLightMode()
    } else {
      setDarkMode()
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
              <Typography level={'title-lg'}>游戏所在路径</Typography>
            </div>
            <div className={styles['value']}>
              <Button
                size={'sm'}
                variant={'soft'}
              >
                <BuildTwoToneIcon />
                设置游戏所在路径
              </Button>
            </div>
          </div>

          <div className={styles['setting-item']}>
            <div className={styles['label']}>
              <Typography level={'title-lg'}>应用程序所在路径</Typography>
            </div>
            <div className={styles['value']}>
              <Button
                size={'sm'}
                variant={'soft'}
              >
                <FolderTwoToneIcon />
                打开应用程序所在路径
              </Button>
            </div>
          </div>
        </div>
      </Modal>
    </React.Fragment>
  )
}

export default Setting
