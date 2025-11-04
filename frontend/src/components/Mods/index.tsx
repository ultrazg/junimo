import React, { useEffect, useState } from 'react'
import { Tabs, TabList, TabPanel, Tab } from '@mui/joy'
import { EventsOn } from 'wailsjs/runtime'
import { ReadModConfigFile } from '@/utils'
import styles from './index.module.scss'
import { ModConfigModal } from '@/components'
import EnabledMods from './components/EnabledMods'
import DisabledMods from './components/DisabledMods'
import { ModManifestType } from '@/types'

const Mods = () => {
  const [enabledMods, setEnabledMods] = useState<ModManifestType[]>([])
  const [enabledModsTotal, setEnabledModsTotal] = useState<number>(0)
  const [disabledMods, setDisabledMods] = useState<ModManifestType[]>([])
  const [disabledModsTotal, setDisabledModsTotal] = useState<number>(0)
  const [modConfigModal, setModConfigModal] = useState<{
    modName: string
    open: boolean
    configStr: string
    path: string
  }>({
    modName: '',
    open: false,
    configStr: '',
    path: '',
  })

  useEffect(() => {
    const onSyncMods = EventsOn(
      'loadActiveMods',
      (data: { mods: ModManifestType[]; total: number }) => {
        setEnabledMods(data.mods)
        setEnabledModsTotal(data.total)
      },
    )

    return () => {
      onSyncMods()
    }
  })

  const onEditConfigFile = (modName: string, configPath: string) => {
    ReadModConfigFile(configPath).then((configStr) => {
      setModConfigModal({
        modName: modName,
        open: true,
        configStr,
        path: configPath,
      })
    })
  }

  return (
    <div className={styles['mods-wrapper']}>
      <Tabs defaultValue={0}>
        <TabList
          variant={'plain'}
          size={'sm'}
        >
          <Tab>已启用({enabledModsTotal})</Tab>
          <Tab>已禁用</Tab>
        </TabList>
        <TabPanel
          value={0}
          style={{ padding: 0 }}
        >
          <EnabledMods
            mods={enabledMods}
            onEditConfigFileFunc={onEditConfigFile}
          />
        </TabPanel>
        <TabPanel
          value={1}
          style={{ padding: 0 }}
        >
          <DisabledMods />
        </TabPanel>
      </Tabs>

      <ModConfigModal
        modName={modConfigModal.modName}
        path={modConfigModal.path}
        configStr={modConfigModal.configStr}
        open={modConfigModal.open}
        onClose={() =>
          setModConfigModal({
            modName: '',
            open: false,
            configStr: '',
            path: '',
          })
        }
      />
    </div>
  )
}

export default Mods
