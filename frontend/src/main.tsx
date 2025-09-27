import React from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import { CssVarsProvider } from '@mui/joy/styles'
import { CssBaseline } from '@mui/joy'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
  <React.StrictMode>
    <CssVarsProvider>
      <CssBaseline />
      <App />
    </CssVarsProvider>
  </React.StrictMode>,
)
