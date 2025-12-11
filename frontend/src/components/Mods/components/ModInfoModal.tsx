import React, { useEffect, useState } from 'react'
import { Modal } from '@/components'
import { ModManifestType } from '@/types'
import { BrowserOpenURL } from 'wailsjs/runtime'
import { ViewModSize, formatBytes } from '@/utils'
import { Link } from '@mui/joy'

type IProps = {
  open: boolean
  onClose: () => void
  mod: ModManifestType
}

const ModInfoModal = ({ open, onClose, mod }: IProps) => {
  const [size, setSize] = useState(0)
  const [loading, setLoading] = useState<boolean>(false)

  const calcModSize = () => {
    setLoading(true)
    ViewModSize(mod.modPath)
      .then((res) => {
        setSize(res)
      })
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    if (open) {
      calcModSize()
    }

    return () => {
      setSize(0)
      setLoading(false)
    }
  }, [open])

  return (
    <React.Fragment>
      <Modal
        title={mod.name}
        open={open}
        onClose={onClose}
      >
        <div>
          <span style={{ color: '#ccc' }}>作者：</span>
          {mod.author}
        </div>
        <div>
          <span style={{ color: '#ccc' }}>已安装版本：</span>v{mod.version}
        </div>
        <div>
          <span style={{ color: '#ccc' }}>文件大小：</span>
          {loading ? '计算中...' : formatBytes(size)}
        </div>
        <div>
          <span style={{ color: '#ccc' }}>描述：</span>
          {mod.description}
        </div>
        {mod.nexusKey !== 0 && (
          <React.Fragment>
            <div>
              <span style={{ color: '#ccc' }}>Nexus 尾号：</span>
              {mod.nexusKey}
            </div>
            <div style={{ display: 'flex' }}>
              <span style={{ color: '#ccc' }}>Nexus 主页：</span>
              <Link
                style={{ flex: 1 }}
                title={'点击访问'}
                onClick={() =>
                  BrowserOpenURL(
                    `https://www.nexusmods.com/stardewvalley/mods/${mod.nexusKey}`,
                  )
                }
              >
                {`https://www.nexusmods.com/stardewvalley/mods/${mod.nexusKey}`}
              </Link>
            </div>
          </React.Fragment>
        )}
      </Modal>
    </React.Fragment>
  )
}

export default ModInfoModal
