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

export type SnackbarOptions = {
  color?: 'primary' | 'neutral' | 'danger' | 'success' | 'warning'
  variant?: 'soft' | 'solid' | 'outlined' | 'plain'
  autoHideDuration?: number
  showIcon?: boolean
}
