import { WindowSetLightTheme, WindowSetDarkTheme } from 'wailsjs/runtime'

const setLightMode = () => {
  WindowSetLightTheme()
}

const setDarkMode = () => {
  WindowSetDarkTheme()
}

export { setLightMode, setDarkMode }
