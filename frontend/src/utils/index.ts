import { setLightMode, setDarkMode } from './theme'
import { APP_NAME, APP_VERSION } from './env'
import {
  ReadConfig,
  UpdateConfig,
  LoadMods,
  OpenGameDir,
  OpenAppDir,
  OpenModDir,
  ReadModConfigFile,
  RemoveModDir,
  BackupModDir,
  UpdateModConfigFile,
  ListBackupDirs,
} from 'wailsjs/go/backend/App'
import { EventsOn } from 'wailsjs/runtime'
import { onSaveGamePath } from './fs'
import { snackbar } from '@/providers/SnackbarProvider'

export {
  setLightMode,
  setDarkMode,
  APP_NAME,
  APP_VERSION,
  ReadConfig,
  UpdateConfig,
  onSaveGamePath,
  snackbar,
  LoadMods,
  EventsOn,
  OpenGameDir,
  OpenAppDir,
  OpenModDir,
  ReadModConfigFile,
  RemoveModDir,
  BackupModDir,
  UpdateModConfigFile,
  ListBackupDirs,
}
