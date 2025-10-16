import React, { useEffect, useState } from 'react'
import { Table, IconButton, Tooltip } from '@mui/joy'
import { EventsOn } from 'wailsjs/runtime'
import { OpenModDir, ReadModConfigFile } from '@/utils'
import styles from './index.module.scss'
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined'
import DoDisturbOnOutlinedIcon from '@mui/icons-material/DoDisturbOnOutlined'
import BuildCircleOutlinedIcon from '@mui/icons-material/BuildCircleOutlined'
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined'
import FolderOpenOutlinedIcon from '@mui/icons-material/FolderOpenOutlined'
import { ModConfigModal, DeleteModModal } from '@/components'

type ModManifestType = {
  name: string
  author: string
  version: string
  minimumApiVersion: string
  description: string
  uniqueID: string
  entryDll: string
  updateKeys: string[]
  modPath: string
  configPath: string
}

const Mods = () => {
  const [mods, setMods] = useState<ModManifestType[]>()
  const [modConfigModal, setModConfigModal] = useState<{
    modName: string
    open: boolean
    configStr: string
  }>({
    modName: '',
    open: false,
    configStr: '',
  })
  const [deleteModModal, setDeleteModModal] = useState<{
    open: boolean
    path: string
  }>({
    open: false,
    path: '',
  })

  useEffect(() => {
    const onSyncMods = EventsOn(
      'loadMods',
      (data: { mods: ModManifestType[]; total: number }) => {
        setMods(data.mods)
      },
    )

    return () => {
      onSyncMods()
    }
  })

  const onOpenModDir = (path: string) => {
    OpenModDir(path).then()
  }

  const onEditConfigFile = (modName: string, configPath: string) => {
    ReadModConfigFile(configPath).then((configStr) => {
      setModConfigModal({
        modName: modName,
        open: true,
        configStr,
      })
    })
  }

  return (
    <div className={styles['mods-wrapper']}>
      <Table
        className={styles['table-wrapper']}
        stickyFooter={false}
        stickyHeader
        stripe="even"
        variant="soft"
      >
        <thead>
          <tr>
            <th style={{ width: 50 }}>序号</th>
            <th style={{ width: 200 }}>Mod 名称</th>
            <th style={{ width: 80 }}>版本</th>
            <th style={{ width: 150 }}>作者</th>
            <th>描述</th>
            <th style={{ width: 180 }}>操作</th>
          </tr>
        </thead>

        <tbody>
          {mods?.map((mod, index) => (
            <tr key={mod.name + mod.version}>
              <td>{index + 1}</td>
              <td>{mod.name}</td>
              <td>v{mod.version}</td>
              <td>{mod.author}</td>
              <td>{mod.description}</td>
              <td>
                <Tooltip
                  title={'Mod 信息'}
                  placement={'top'}
                  variant={'outlined'}
                >
                  <IconButton
                    variant={'plain'}
                    size={'sm'}
                  >
                    <InfoOutlinedIcon />
                  </IconButton>
                </Tooltip>

                <Tooltip
                  title={'禁用 Mod'}
                  placement={'top'}
                  variant={'outlined'}
                >
                  <IconButton
                    variant={'plain'}
                    size={'sm'}
                  >
                    <DoDisturbOnOutlinedIcon />
                  </IconButton>
                </Tooltip>

                <Tooltip
                  title={'打开 Mod 配置文件'}
                  placement={'top'}
                  variant={'outlined'}
                >
                  <IconButton
                    variant={'plain'}
                    size={'sm'}
                    disabled={mod.configPath === ''}
                    onClick={() => {
                      onEditConfigFile(mod.name, mod.configPath)
                    }}
                  >
                    <BuildCircleOutlinedIcon />
                  </IconButton>
                </Tooltip>

                <Tooltip
                  title={'打开 Mod 目录'}
                  placement={'top'}
                  variant={'outlined'}
                >
                  <IconButton
                    variant={'plain'}
                    size={'sm'}
                    onClick={() => {
                      onOpenModDir(mod.modPath)
                    }}
                  >
                    <FolderOpenOutlinedIcon />
                  </IconButton>
                </Tooltip>

                <Tooltip
                  title={'删除 Mod'}
                  placement={'top'}
                  variant={'outlined'}
                >
                  <IconButton
                    variant={'plain'}
                    size={'sm'}
                    color={'danger'}
                    onClick={() => {
                      setDeleteModModal({ open: true, path: mod.modPath })
                    }}
                  >
                    <DeleteOutlineOutlinedIcon />
                  </IconButton>
                </Tooltip>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>

      <ModConfigModal
        modName={modConfigModal.modName}
        configStr={modConfigModal.configStr}
        open={modConfigModal.open}
        onClose={() =>
          setModConfigModal({ modName: '', open: false, configStr: '' })
        }
      />

      <DeleteModModal
        path={deleteModModal.path}
        open={deleteModModal.open}
        onClose={() => setDeleteModModal({ open: false, path: '' })}
      />
    </div>
  )
}

export default Mods
