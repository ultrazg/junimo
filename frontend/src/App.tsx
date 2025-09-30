import React, { Fragment, useEffect } from 'react'
import { MenuBar } from '@/components'
import { setLightMode, setDarkMode, getStyleMode } from '@/utils'

const App = () => {
  const initStyleMode = () => {
    getStyleMode().then((style) => {
      if (style === 'light') {
        setLightMode()
      } else {
        setDarkMode()
      }
    })
  }

  useEffect(() => {
    initStyleMode()
  }, [])

  return (
    <Fragment>
      <MenuBar />
    </Fragment>
  )
}

export default App
