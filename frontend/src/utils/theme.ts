import { WindowSetLightTheme, WindowSetDarkTheme } from 'wailsjs/runtime'
import { ReadConfig } from './index'

const setLightMode = () => {
  WindowSetLightTheme()
}

const setDarkMode = () => {
  WindowSetDarkTheme()
}

const getStyleMode = async () => {
  const { theme } = await ReadConfig()

  return theme || 'light'
}

export { setLightMode, setDarkMode, getStyleMode }
