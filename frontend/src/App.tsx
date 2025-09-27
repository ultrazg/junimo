import React from 'react'
import { Button } from '@mui/joy'
import { useColorScheme } from '@mui/joy/styles'

function App() {
  const { mode, setMode } = useColorScheme()

  return (
    <div>
      hello world
      <Button
        size={'sm'}
        variant={'soft'}
        onClick={() => setMode(mode === 'light' ? 'dark' : 'light')}
      >
        {mode === 'light' ? 'dark' : 'light'}
      </Button>
    </div>
  )
}

export default App
