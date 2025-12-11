import React, { useEffect, useState } from 'react'
import { Tabs, TabList, TabPanel, Tab } from '@mui/joy'
import { EventsOn } from 'wailsjs/runtime'
import { ReadModConfigFile } from '@/utils'
import styles from './index.module.scss'
import { ModConfigModal } from '@/components'
import EnabledMods from './components/EnabledMods'
import DisabledMods from './components/DisabledMods'
import ModInfoModal from './components/ModInfoModal'
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
  const [modInfoModal, setModInfoModal] = useState<{
    mod: ModManifestType
    open: boolean
  }>({
    mod: {
      name: '',
      author: '',
      version: '',
      minimumApiVersion: '',
      description: '',
      uniqueID: '',
      entryDll: '',
      updateKeys: [],
      nexusKey: 0,
      manifestPath: '',
      modPath: '',
      configPath: '',
    },
    open: false,
  })

  useEffect(() => {
    const onLoadEnabledMods = EventsOn(
      'mod:loadEnabled',
      (data: { mods: ModManifestType[]; total: number }) => {
        setEnabledMods(data.mods)
        setEnabledModsTotal(data.total)
      },
    )

    const onDisabledMods = EventsOn(
      'mod:loadDisabled',
      (data: { mods: ModManifestType[]; total: number }) => {
        setDisabledMods(data.mods)
        setDisabledModsTotal(data.total)
      },
    )

    return () => {
      onLoadEnabledMods()
      onDisabledMods()
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

  const onModInfoShow = (mod: ModManifestType) => {
    setModInfoModal({
      mod,
      open: true,
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
          <Tab>已禁用({disabledModsTotal})</Tab>
        </TabList>
        <TabPanel
          value={0}
          style={{ padding: 0 }}
        >
          <EnabledMods
            mods={enabledMods}
            onEditConfigFileFunc={onEditConfigFile}
            onModInfoFunc={onModInfoShow}
          />
        </TabPanel>
        <TabPanel
          value={1}
          style={{ padding: 0 }}
        >
          <DisabledMods
            mods={disabledMods}
            onModInfoFunc={onModInfoShow}
          />
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

      <ModInfoModal
        mod={modInfoModal.mod}
        open={modInfoModal.open}
        onClose={() =>
          setModInfoModal({
            mod: {
              name: '',
              author: '',
              version: '',
              minimumApiVersion: '',
              description: '',
              uniqueID: '',
              entryDll: '',
              updateKeys: [],
              nexusKey: 0,
              manifestPath: '',
              modPath: '',
              configPath: '',
            },
            open: false,
          })
        }
      />
    </div>
  )
}

export default Mods
