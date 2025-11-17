import React, { useState } from 'react'
import { ModManifestType } from '@/types'
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined'
import DoDisturbOnOutlinedIcon from '@mui/icons-material/DoDisturbOnOutlined'
import BuildCircleOutlinedIcon from '@mui/icons-material/BuildCircleOutlined'
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined'
import FolderOpenOutlinedIcon from '@mui/icons-material/FolderOpenOutlined'
import { IconButton, Table, Tooltip } from '@mui/joy'
import styles from './index.module.scss'
import {
  OpenModDir,
  RemoveModDir,
  DisableMod,
  ViewSpecifiedModFile,
} from '@/utils'
import { DeleteModal } from '@/components'

type IProps = {
  mods: ModManifestType[]
  onEditConfigFileFunc: (modName: string, configPath: string) => void
}

const EnabledMods: React.FC<IProps> = ({ mods, onEditConfigFileFunc }) => {
  const [deleteModal, setDeleteModal] = useState<{
    loading: boolean
    open: boolean
    path: string
  }>({
    loading: false,
    open: false,
    path: '',
  })

  const onRemoveModFunc = () => {
    setDeleteModal({
      ...deleteModal,
      loading: true,
    })

    RemoveModDir(deleteModal.path)
      .then()
      .finally(() => {
        setDeleteModal({
          open: false,
          path: '',
          loading: false,
        })
      })
  }

  return (
    <div className={styles['mods-wrapper']}>
      <Table
        className={styles['table']}
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
              <td title={mod.name}>{mod.name}</td>
              <td title={mod.version}>v{mod.version}</td>
              <td title={mod.author}>{mod.author}</td>
              <td title={mod.description}>{mod.description}</td>
              <td>
                <Tooltip
                  title={'Mod 信息'}
                  placement={'top'}
                  variant={'outlined'}
                >
                  <IconButton
                    variant={'plain'}
                    size={'sm'}
                    color={'primary'}
                    onClick={() => {
                      // TODO
                      if (mod.updateKeys) {
                        const ids = mod.updateKeys
                          .map((key) => key.match(/^Nexus:(\d+)$/))
                          .filter(Boolean)
                          .map((match) => match?.[1])

                        ViewSpecifiedModFile(String(ids)).then((res) => {
                          console.log(res)
                        })
                      }
                    }}
                  >
                    <InfoOutlinedIcon />
                  </IconButton>
                </Tooltip>

                <Tooltip
                  title={'修改 Mod 配置文件'}
                  placement={'top'}
                  variant={'outlined'}
                >
                  <IconButton
                    variant={'plain'}
                    size={'sm'}
                    disabled={mod.configPath === ''}
                    onClick={() => {
                      onEditConfigFileFunc(mod.name, mod.configPath)
                    }}
                    color={'primary'}
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
                      OpenModDir(mod.manifestPath).then()
                    }}
                    color={'primary'}
                  >
                    <FolderOpenOutlinedIcon />
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
                    onClick={() => {
                      DisableMod(mod.modPath).then()
                    }}
                    color={'primary'}
                  >
                    <DoDisturbOnOutlinedIcon />
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
                      setDeleteModal({
                        open: true,
                        path: mod.modPath,
                        loading: false,
                      })
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

      <DeleteModal
        open={deleteModal.open}
        onClose={() =>
          setDeleteModal({ open: false, path: '', loading: false })
        }
        loading={deleteModal.loading}
        onOk={onRemoveModFunc}
      />
    </div>
  )
}

export default EnabledMods
