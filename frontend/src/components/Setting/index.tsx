import React from 'react'
import { Modal } from '@/components'
import { Button, Radio, RadioGroup, Typography } from '@mui/joy'
import styles from './index.module.scss'
import FolderTwoToneIcon from '@mui/icons-material/FolderTwoTone'
import BuildTwoToneIcon from '@mui/icons-material/BuildTwoTone'

type IProps = {
  open: boolean
  onClose: () => void
}

const THEME_LISTS = ['system', 'light', 'dark']

const Setting: React.FC<IProps> = ({ open, onClose }) => {
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
              <RadioGroup orientation={'horizontal'}>
                {THEME_LISTS.map((item) => (
                  <Radio
                    key={item}
                    value={item}
                    label={
                      {
                        system: <Typography>跟随系统</Typography>,
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
