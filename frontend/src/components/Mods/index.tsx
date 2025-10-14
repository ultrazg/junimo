import React from 'react'
import { LoadMods } from '@/utils'
import { Button } from '@mui/joy'
import styles from './index.module.scss'

const Mods = () => {
  const loadMods = () => {
    LoadMods().then()
  }

  return (
    <div className={styles['mods-wrapper']}>
      <Button onClick={loadMods}>加载 mods</Button>
    </div>
  )
}

export default Mods
