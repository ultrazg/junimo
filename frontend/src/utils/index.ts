import { setLightMode, setDarkMode } from './theme'
import { APP_NAME, APP_VERSION } from './env'
import {
  ReadConfig,
  UpdateConfig,
  LoadEnabledMods,
  OpenGameDir,
  OpenAppDir,
  OpenModDir,
  ReadModConfigFile,
  RemoveModDir,
  RemoveBackupDir,
  BackupModDir,
  UpdateModConfigFile,
  ListBackupDirs,
  OpenBackupDir,
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
  LoadEnabledMods,
  EventsOn,
  OpenGameDir,
  OpenAppDir,
  OpenModDir,
  ReadModConfigFile,
  RemoveModDir,
  RemoveBackupDir,
  BackupModDir,
  UpdateModConfigFile,
  ListBackupDirs,
  OpenBackupDir,
}
