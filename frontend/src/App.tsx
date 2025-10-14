import React, { Fragment, useEffect } from 'react'
import { MenuBar, Mods } from '@/components'
import {
  setLightMode,
  setDarkMode,
  ReadConfig,
  EventsOn,
  snackbar,
} from '@/utils'
import { SnackbarOptions } from '@/providers/SnackbarProvider'

const App = () => {
  const initStyleMode = () => {
    ReadConfig('theme').then((res) => {
      if (res === 'light') {
        setLightMode()
      } else {
        setDarkMode()
      }
    })
  }

  useEffect(() => {
    initStyleMode()

    const showSnackbar = EventsOn(
      'snackbarShow',
      (options: SnackbarOptions & { message: string }) => {
        snackbar.show(options.message, {
          showIcon: options.showIcon,
          autoHideDuration: options.autoHideDuration,
          color: options.color,
          variant: options.variant,
        })
      },
    )

    return () => {
      showSnackbar()
    }
  }, [])

  return (
    <Fragment>
      <MenuBar />
      <Mods />
    </Fragment>
  )
}

export default App
