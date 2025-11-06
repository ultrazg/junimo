import React from 'react'
import { ModManifestType } from '@/types'
import styles from './index.module.scss'
import { IconButton, Tooltip, Table } from '@mui/joy'
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined'
import { OpenModDir, EnableMod } from '@/utils'
import FolderOpenOutlinedIcon from '@mui/icons-material/FolderOpenOutlined'
import CheckCircleOutlinedIcon from '@mui/icons-material/CheckCircleOutlined'

type IProps = {
  mods: ModManifestType[]
}

const DisabledMods: React.FC<IProps> = ({ mods }) => {
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
            <th style={{ width: 140 }}>操作</th>
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
                  >
                    <InfoOutlinedIcon />
                  </IconButton>
                </Tooltip>

                <Tooltip
                  title={'启用 Mod'}
                  placement={'top'}
                  variant={'outlined'}
                >
                  <IconButton
                    variant={'plain'}
                    size={'sm'}
                    onClick={() => {
                      EnableMod(mod.modPath).then()
                    }}
                    color={'primary'}
                  >
                    <CheckCircleOutlinedIcon />
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
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
  )
}

export default DisabledMods
