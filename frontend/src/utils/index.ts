import { setLightMode, setDarkMode, getStyleMode } from './theme'
import { APP_NAME, APP_VERSION } from './env'
import { ReadConfig, UpdateConfig } from 'wailsjs/go/backend/App'
import { onSaveGamePath } from './fs'
import { snackbar } from '@/providers/SnackbarProvider'

export {
  setLightMode,
  setDarkMode,
  getStyleMode,
  APP_NAME,
  APP_VERSION,
  ReadConfig,
  UpdateConfig,
  onSaveGamePath,
  snackbar,
}
