import React, { useEffect, useState } from 'react'
import { Table } from '@mui/joy'
import { EventsOn } from 'wailsjs/runtime'
import styles from './index.module.scss'

type ModManifestType = {
  name: string
  author: string
  version: string
  minimumApiVersion: string
  description: string
  uniqueID: string
  entryDll: string
  updateKeys: string[]
}

const Mods = () => {
  const [mods, setMods] = useState<ModManifestType[]>()
  const [total, setTotal] = useState<number>(0)

  useEffect(() => {
    const onSyncMods = EventsOn(
      'loadMods',
      (data: { mods: ModManifestType[]; total: number }) => {
        setMods(data.mods)
        setTotal(data.total)
      },
    )

    return () => {
      onSyncMods()
    }
  })

  return (
    <div className={styles['mods-wrapper']}>
      <Table
        stickyFooter={false}
        stickyHeader
        stripe="even"
        variant="soft"
      >
        <thead>
          <tr>
            <th style={{ width: 50 }}>序号</th>
            <th style={{ width: 200 }}>Mod 名称</th>
            <th style={{ width: 100 }}>版本</th>
            <th style={{ width: 150 }}>作者</th>
            <th>描述</th>
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
            </tr>
          ))}
        </tbody>
      </Table>
    </div>
  )
}

export default Mods
