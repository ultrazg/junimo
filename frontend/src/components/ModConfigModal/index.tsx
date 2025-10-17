import React, { useEffect, useState } from 'react'
import { Modal } from '@/components'
import { ReadConfig, UpdateModConfigFile } from '@/utils'
import AceEditor from 'react-ace'
import { Button, Typography } from '@mui/joy'
import SaveOutlinedIcon from '@mui/icons-material/SaveOutlined'
import styles from './index.module.scss'

import 'ace-builds/src-noconflict/mode-json'
import 'ace-builds/src-noconflict/theme-github'
import 'ace-builds/src-noconflict/theme-twilight'

type IProps = {
  modName: string
  path: string
  configStr: string
  open: boolean
  onClose: () => void
}

const ModConfigModal: React.FC<IProps> = ({
  modName,
  path,
  configStr,
  open,
  onClose,
}) => {
  const [theme, setTheme] = useState<'light' | 'dark'>('light')
  const [configValue, setConfigValue] = useState<string>('')
  const [submitLoading, setSubmitLoading] = useState<boolean>(false)

  const init = () => {
    ReadConfig('theme').then((res: 'light' | 'dark') => {
      setTheme(res)
    })

    setConfigValue(configStr)
  }

  const onChangeConfigValue = (value: string) => {
    setConfigValue(value)
  }

  const onUpdateConfig = () => {
    setSubmitLoading(true)
    UpdateModConfigFile(path, configValue)
      .then()
      .finally(() => {
        setSubmitLoading(false)
      })
  }

  useEffect(() => {
    if (open) {
      init()
    }

    return () => {
      setConfigValue('')
    }
  }, [open])

  return (
    <React.Fragment>
      <Modal
        open={open}
        onClose={onClose}
        title={`${modName} 的配置`}
      >
        <AceEditor
          mode={'json'}
          theme={theme === 'light' ? 'github' : 'twilight'}
          highlightActiveLine={true}
          showGutter={true}
          value={configValue}
          onChange={onChangeConfigValue}
        />

        <div className={styles['mod-config-modal-option-button']}>
          <div className={styles['tip']}>
            <Typography
              color={'danger'}
              level={'body-sm'}
            >
              注意：请确保配置文件格式，修改前关闭游戏
            </Typography>
          </div>

          <div className={styles['option-button']}>
            <Button
              size={'sm'}
              variant={'soft'}
              loading={submitLoading}
              onClick={onUpdateConfig}
            >
              <SaveOutlinedIcon />
              修改配置
            </Button>
          </div>
        </div>
      </Modal>
    </React.Fragment>
  )
}

export default ModConfigModal
