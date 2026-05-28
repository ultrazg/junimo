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
  modPath: string
  configPath: string
  uniqueID: string
  latestFileID: number
  latestFileName: string
}

export type CheckForUpdatesResultType = {
  success: boolean
  message: string
  total: number
  items: ModUpdateInfoType[]
}

export type ModChangelogEntryType = {
  version: string
  changes: string[]
}

export type ViewModChangelogResultType = {
  success: boolean
  message: string
  entries: ModChangelogEntryType[]
}

export type UpdateModResultType = {
  success: boolean
  message: string
  backupName: string
  newModPath: string
}

export type RollbackModUpdateResultType = {
  success: boolean
  message: string
}

export type ModUpdateBackupType = {
  name: string
  modName: string
  modPath: string
  oldVersion: string
  newVersion: string
  uniqueID: string
  size: number
  createTime: string
  keptConfig: boolean
  originalDir: string
}

export type SnackbarOptions = {
  color?: 'primary' | 'neutral' | 'danger' | 'success' | 'warning'
  variant?: 'soft' | 'solid' | 'outlined' | 'plain'
  autoHideDuration?: number
  showIcon?: boolean
}
