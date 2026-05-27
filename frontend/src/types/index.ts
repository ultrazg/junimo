export type ModManifestType = {
  name: string
  author: string
  version: string
  minimumApiVersion: string
  description: string
  uniqueID: string
  entryDll: string
  updateKeys: string[]
  nexusKey: number
  manifestPath: string
  modPath: string
  configPath: string
}

export type ImportModPreviewType = {
  zipPath: string
  manifest: ModManifestType
  exists: boolean
  existingPath: string
  error: string
}

export type ModUpdateInfoType = {
  name: string
  nexusKey: number
  currentVersion: string
  latestVersion: string
  hasUpdate: boolean
  error: string
}

export type CheckForUpdatesResultType = {
  success: boolean
  message: string
  total: number
  items: ModUpdateInfoType[]
}

export type SnackbarOptions = {
  color?: 'primary' | 'neutral' | 'danger' | 'success' | 'warning'
  variant?: 'soft' | 'solid' | 'outlined' | 'plain'
  autoHideDuration?: number
  showIcon?: boolean
}
