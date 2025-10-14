import { WindowSetLightTheme, WindowSetDarkTheme } from 'wailsjs/runtime'
import { ReadConfig } from './index'

const setLightMode = () => {
  WindowSetLightTheme()
}

const setDarkMode = () => {
  WindowSetDarkTheme()
}

export { setLightMode, setDarkMode }
