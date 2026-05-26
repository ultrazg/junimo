import React, { useEffect, useMemo, useState } from 'react'
import { Button, Chip, Typography } from '@mui/joy'
import { Modal } from '@/components'
import { ConfirmImportMod, snackbar } from '@/utils'
import { backend } from 'wailsjs/go/models'

type IProps = {
  open: boolean
  items: backend.ImportModPreview[]
  onClose: () => void
  onFinished: () => void
}

const ImportModModal: React.FC<IProps> = ({
  open,
  items,
  onClose,
  onFinished,
}) => {
  const [index, setIndex] = useState(0)
  const [loading, setLoading] = useState(false)
  const [touched, setTouched] = useState(false)

  const current = items[index]
  const total = items.length

  const fileName = useMemo(() => {
    if (!current?.zipPath) return ''
    const parts = current.zipPath.split(/[\\/]/)
    return parts[parts.length - 1]
  }, [current?.zipPath])

  useEffect(() => {
    if (open) {
      setIndex(0)
      setLoading(false)
      setTouched(false)
    }
  }, [open])

  if (!open || !current) return null

  const finishIfLast = () => {
    if (index + 1 >= total) {
      if (touched) onFinished()
      onClose()
    } else {
      setIndex(index + 1)
    }
  }

  const onSkip = () => {
    finishIfLast()
  }

  const onConfirm = (replace: boolean) => {
    setLoading(true)
    ConfirmImportMod(current.zipPath, replace)
      .then((res) => {
        if (res.success) {
          setTouched(true)
          snackbar.show(
            `成功${replace ? '替换' : '导入'} Mod：${current.manifest.name || fileName}`,
            { showIcon: true, color: 'success', variant: 'soft' },
          )
          finishIfLast()
        } else {
          snackbar.show(`导入失败：${res.message}`, {
            showIcon: true,
            color: 'danger',
            variant: 'soft',
          })
        }
      })
      .finally(() => setLoading(false))
  }

  if (current.error) {
    return (
      <Modal
        title={`导入 Mod (${index + 1}/${total})`}
        open={open}
        onClose={onClose}
      >
        <div style={{ width: 480 }}>
          <Typography
            level={'body-md'}
            color={'danger'}
          >
            读取 {fileName} 失败：{current.error}
          </Typography>

          <div style={{ marginTop: 16, textAlign: 'right' }}>
            <Button
              size={'sm'}
              variant={'soft'}
              onClick={onSkip}
            >
              {index + 1 >= total ? '关闭' : '跳过'}
            </Button>
          </div>
        </div>
      </Modal>
    )
  }

  const m = current.manifest

  return (
    <Modal
      title={`导入 Mod (${index + 1}/${total})`}
      open={open}
      onClose={onClose}
    >
      <div style={{ width: 480 }}>
        <div style={{ marginBottom: 8 }}>
          <Chip
            size={'sm'}
            variant={'soft'}
            color={current.exists ? 'warning' : 'primary'}
          >
            {current.exists ? '该 Mod 已存在' : '新安装'}
          </Chip>
          <span style={{ marginLeft: 8, color: '#999', fontSize: 12 }}>
            {fileName}
          </span>
        </div>

        <div>
          <span style={{ color: '#999' }}>名称：</span>
          {m.name || '-'}
        </div>
        <div>
          <span style={{ color: '#999' }}>作者：</span>
          {m.author || '-'}
        </div>
        <div>
          <span style={{ color: '#999' }}>版本：</span>v{m.version || '-'}
        </div>
        {/*{m.uniqueID && (*/}
        {/*  <div>*/}
        {/*    <span style={{ color: '#999' }}>UniqueID：</span>*/}
        {/*    {m.uniqueID}*/}
        {/*  </div>*/}
        {/*)}*/}
        {m.description && (
          <div>
            <span style={{ color: '#999' }}>描述：</span>
            {m.description}
          </div>
        )}

        {current.exists && (
          <Typography
            level={'body-sm'}
            color={'warning'}
            sx={{ marginTop: 1 }}
          >
            已存在路径：{current.existingPath}
            <br />
            替换会删除当前已安装版本后再写入新版本。
          </Typography>
        )}

        <div
          style={{
            marginTop: 16,
            textAlign: 'right',
            display: 'flex',
            gap: 8,
            justifyContent: 'flex-end',
          }}
        >
          <Button
            size={'sm'}
            variant={'plain'}
            color={'neutral'}
            disabled={loading}
            onClick={onSkip}
          >
            {index + 1 >= total ? '取消' : '跳过'}
          </Button>

          {current.exists ? (
            <Button
              size={'sm'}
              variant={'soft'}
              color={'warning'}
              loading={loading}
              onClick={() => onConfirm(true)}
            >
              替换
            </Button>
          ) : (
            <Button
              size={'sm'}
              variant={'soft'}
              color={'primary'}
              loading={loading}
              onClick={() => onConfirm(false)}
            >
              添加
            </Button>
          )}
        </div>
      </div>
    </Modal>
  )
}

export default ImportModModal
