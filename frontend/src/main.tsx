import React from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import { CssVarsProvider, extendTheme } from '@mui/joy/styles'
import { CssBaseline } from '@mui/joy'

const container = document.getElementById('root')

const root = createRoot(container!)

const theme = extendTheme({
  colorSchemes: {
    light: {
      palette: {
        primary: {
          50: '#f6ffed',
          100: '#d9f7be',
          200: '#b7eb8f',
          300: '#95de64',
          400: '#73d13d',
          500: '#52c41a',
          600: '#389e0d',
          700: '#237804',
          800: '#135200',
          900: '#092b00',
        },
      },
    },
    dark: {
      palette: {
        primary: {
          50: '#f6ffed',
          100: '#d9f7be',
          200: '#b7eb8f',
          300: '#95de64',
          400: '#73d13d',
          500: '#52c41a',
          600: '#389e0d',
          700: '#237804',
          800: '#135200',
          900: '#092b00',
        },
      },
    },
  },
  fontSize: {
    xs: '0.75rem',
    sm: '0.8125rem',
    md: '0.875rem',
    lg: '1rem',
    xl: '1.125rem',
  },
})

root.render(
  <React.StrictMode>
    <CssVarsProvider
      theme={theme}
      defaultColorScheme={'light'}
    >
      <CssBaseline />
      <App />
    </CssVarsProvider>
  </React.StrictMode>,
)
