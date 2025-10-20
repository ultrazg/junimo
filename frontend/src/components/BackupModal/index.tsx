import React, { useEffect, useState } from 'react'
import { Modal } from '@/components'
import { IconButton, Table, Tooltip, CircularProgress } from '@mui/joy'
import { ListBackupDirs } from '@/utils'
import FolderOpenOutlinedIcon from '@mui/icons-material/FolderOpenOutlined'
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined'

type IProps = {
  open: boolean
  onClose: () => void
}

const formatTime = (TimeStr: string) => {
  const date: any = new Date(TimeStr)
  if (isNaN(date)) return TimeStr

  const pad = (n: any) => n.toString().padStart(2, '0')

  const yyyy = date.getFullYear()
  const mm = pad(date.getMonth() + 1)
  const dd = pad(date.getDate())
  const hh = pad(date.getHours())
  const mi = pad(date.getMinutes())
  const ss = pad(date.getSeconds())

  return `${yyyy}-${mm}-${dd} ${hh}:${mi}:${ss}`
}

const BackupModal: React.FC<IProps> = ({ open, onClose }) => {
  const [loading, setLoading] = useState<boolean>(true)
  const [dirs, setDirs] = useState<
    { name: string; size: number; createTime: string }[]
  >([])

  const onListBackupDirs = () => {
    setLoading(true)
    ListBackupDirs()
      .then((dirs) => {
        const temp: { name: string; size: number; createTime: string }[] = []

        dirs.map((dir) => {
          temp.push({
            name: dir.name,
            size: dir.size,
            createTime: dir.createTime.toLocaleString(),
          })
        })

        setDirs(temp)
      })
      .finally(() => {
        setLoading(false)
      })
  }

  useEffect(() => {
    if (open) {
      onListBackupDirs()
    }

    return () => {
      setDirs([])
    }
  }, [open])

  return (
    <Modal
      title={`备份${dirs.length > 0 ? `(${dirs.length})`: ''}`}
      open={open}
      onClose={onClose}
    >
      {loading ? (
        <div style={{ textAlign: 'center' }}>
          <CircularProgress />
        </div>
      ) : dirs.length > 0 ? (
        <div
          style={{
            width: '700px',
            cursor: 'default',
            maxHeight: '500px',
            overflowY: 'auto',
          }}
        >
          <Table
            stickyHeader
            stripe="even"
            variant="soft"
          >
            <thead>
              <tr>
                <th style={{ width: 50 }}>序号</th>
                <th style={{ width: 250 }}>备份名称</th>
                <th style={{ width: 100 }}>备份大小</th>
                <th style={{ width: 150 }}>备份时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {dirs.map((dir, index) => (
                <tr key={dir.name}>
                  <td>{index + 1}</td>
                  <td>{dir.name}</td>
                  <td>{(dir.size / 1024 / 1024).toFixed(2)} MB</td>
                  <td>{formatTime(dir.createTime)}</td>
                  <td>
                    <Tooltip
                      title={'打开备份目录'}
                      placement={'top'}
                      variant={'outlined'}
                    >
                      <IconButton
                        variant={'plain'}
                        size={'sm'}
                        onClick={() => {}}
                      >
                        <FolderOpenOutlinedIcon />
                      </IconButton>
                    </Tooltip>

                    <Tooltip
                      title={'删除备份'}
                      placement={'top'}
                      variant={'outlined'}
                    >
                      <IconButton
                        variant={'plain'}
                        size={'sm'}
                        color={'danger'}
                        onClick={() => {}}
                      >
                        <DeleteOutlineOutlinedIcon />
                      </IconButton>
                    </Tooltip>
                  </td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      ) : (
        <div style={{ textAlign: 'center' }}>暂无备份</div>
      )}
    </Modal>
  )
}

export default BackupModal
