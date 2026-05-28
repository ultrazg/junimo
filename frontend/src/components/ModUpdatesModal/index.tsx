import React, { useEffect, useState } from 'react'
import {
  Button,
  CircularProgress,
  Table,
  Tooltip,
  Typography,
  IconButton,
  Chip,
} from '@mui/joy'
import LinkIcon from '@mui/icons-material/Link'
import AccessTimeOutlinedIcon from '@mui/icons-material/AccessTimeOutlined'
import CloudDownloadOutlinedIcon from '@mui/icons-material/CloudDownloadOutlined'
import RestoreOutlinedIcon from '@mui/icons-material/RestoreOutlined'
import HistoryOutlinedIcon from '@mui/icons-material/HistoryOutlined'
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined'
import { Modal } from '@/components'
import {
  BrowserOpenURL,
  ViewModChangelog,
  UpdateMod,
  ListModUpdateBackups,
  RollbackModUpdate,
  RemoveModUpdateBackup,
  formatBytes,
  snackbar,
} from '@/utils'
import {
  CheckForUpdatesResultType,
  ModChangelogEntryType,
  ModUpdateInfoType,
  ModUpdateBackupType,
} from '@/types'

type IProps = {
  open: boolean
  loading: boolean
  result?: CheckForUpdatesResultType
  onClose: () => void
  onRefresh: () => void
}

type ChangelogState = {
  open: boolean
  loading: boolean
  modName: string
  entries: ModChangelogEntryType[]
  message: string
}

type KeepConfigState = {
  open: boolean
  item?: ModUpdateInfoType
}

const initialChangelog: ChangelogState = {
  open: false,
  loading: false,
  modName: '',
  entries: [],
  message: '',
}

const formatTime = (timeStr: string) => {
  const date: any = new Date(timeStr)
  if (isNaN(date)) return timeStr
  const pad = (n: any) => n.toString().padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(
    date.getDate(),
  )} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(
    date.getSeconds(),
  )}`
}

const ModUpdatesModal: React.FC<IProps> = ({
  open,
  loading,
  result,
  onClose,
  onRefresh,
}) => {
  const updatable = (result?.items || []).filter((it) => it.hasUpdate)
  const errors = (result?.items || []).filter((it) => !!it.error)

  const [changelog, setChangelog] = useState<ChangelogState>(initialChangelog)
  const [updatingKey, setUpdatingKey] = useState<number>(0)
  const [updatedKeys, setUpdatedKeys] = useState<Record<number, boolean>>({})
  const [keepConfig, setKeepConfig] = useState<KeepConfigState>({ open: false })
  const [rollbackOpen, setRollbackOpen] = useState<boolean>(false)
  const [rollbackLoading, setRollbackLoading] = useState<boolean>(false)
  const [backups, setBackups] = useState<ModUpdateBackupType[]>([])
  const [busyBackup, setBusyBackup] = useState<string>('')

  useEffect(() => {
    if (open) {
      setUpdatedKeys({})
    }
  }, [open])

  const onViewChangelog = (it: ModUpdateInfoType) => {
    setChangelog({
      open: true,
      loading: true,
      modName: it.name,
      entries: [],
      message: '',
    })

    ViewModChangelog(it.nexusKey)
      .then((res) => {
        setChangelog({
          open: true,
          loading: false,
          modName: it.name,
          entries: res.entries || [],
          message: res.success ? '' : res.message || '获取更新日志失败',
        })
      })
      .catch((err) => {
        setChangelog({
          open: true,
          loading: false,
          modName: it.name,
          entries: [],
          message: String(err),
        })
      })
  }

  const onClickUpdate = (it: ModUpdateInfoType) => {
    if (!it.modPath || !it.latestFileID) {
      snackbar.show('缺少更新所需信息（modPath/fileID）', {
        showIcon: true,
        color: 'danger',
        variant: 'soft',
      })
      return
    }

    if (it.configPath) {
      setKeepConfig({ open: true, item: it })
      return
    }

    runUpdate(it, false)
  }

  const runUpdate = (it: ModUpdateInfoType, keep: boolean) => {
    setKeepConfig({ open: false })
    setUpdatingKey(it.nexusKey)

    UpdateMod(it.modPath, it.nexusKey, it.latestFileID, keep)
      .then((res) => {
        if (res.success) {
          snackbar.show(res.message || `已更新 ${it.name}`, {
            showIcon: true,
            color: 'success',
            variant: 'soft',
          })
          setUpdatedKeys((prev) => ({ ...prev, [it.nexusKey]: true }))
        } else {
          snackbar.show(res.message || '更新失败', {
            showIcon: true,
            color: 'danger',
            variant: 'soft',
            autoHideDuration: 6000,
          })
        }
      })
      .catch((err) => {
        snackbar.show(`更新失败：${String(err)}`, {
          showIcon: true,
          color: 'danger',
          variant: 'soft',
          autoHideDuration: 6000,
        })
      })
      .finally(() => setUpdatingKey(0))
  }

  const loadBackups = () => {
    setRollbackLoading(true)
    ListModUpdateBackups()
      .then((res) => setBackups(res || []))
      .catch(() => setBackups([]))
      .finally(() => setRollbackLoading(false))
  }

  const onOpenRollback = () => {
    setRollbackOpen(true)
    loadBackups()
  }

  const onRollback = (b: ModUpdateBackupType) => {
    setBusyBackup(b.name)
    RollbackModUpdate(b.name)
      .then((res) => {
        snackbar.show(res.message || (res.success ? '已回滚' : '回滚失败'), {
          showIcon: true,
          color: res.success ? 'success' : 'danger',
          variant: 'soft',
        })
      })
      .catch((err) => {
        snackbar.show(`回滚失败：${String(err)}`, {
          showIcon: true,
          color: 'danger',
          variant: 'soft',
        })
      })
      .finally(() => setBusyBackup(''))
  }

  const onRemoveBackup = (b: ModUpdateBackupType) => {
    setBusyBackup(b.name)
    RemoveModUpdateBackup(b.name)
      .then((res) => {
        if (res.success) {
          setBackups((prev) => prev.filter((x) => x.name !== b.name))
        } else {
          snackbar.show(res.message || '删除备份失败', {
            showIcon: true,
            color: 'danger',
            variant: 'soft',
          })
        }
      })
      .finally(() => setBusyBackup(''))
  }

  return (
    <React.Fragment>
      <Modal
        title={`检查 Mod 更新${
          result && !loading
            ? `（共 ${result.total} 个，可更新 ${updatable.length} 个）`
            : ''
        }`}
        open={open}
        onClose={onClose}
      >
        <div style={{ width: 760, minHeight: 120 }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: 24 }}>
              <CircularProgress />
              <div style={{ marginTop: 12 }}>
                正在通过 Nexus Mods API 检查更新...
              </div>
            </div>
          ) : !result ? (
            <div
              style={{ textAlign: 'center', padding: 24, fontStyle: 'italic' }}
            >
              暂无数据
            </div>
          ) : !result.success ? (
            <Typography
              level={'body-md'}
              color={'danger'}
            >
              {result.message || '检查失败'}
            </Typography>
          ) : result.total === 0 ? (
            <div style={{ textAlign: 'center', padding: 24 }}>
              没有带 Nexus 标识的 Mod 可供检查
            </div>
          ) : updatable.length === 0 ? (
            <div style={{ textAlign: 'center', padding: 24 }}>
              所有 Mod 都已是最新版本
              {errors.length > 0 ? `（${errors.length} 个检查失败）` : ''}
            </div>
          ) : (
            <div style={{ maxHeight: 400, overflowY: 'auto' }}>
              <Table
                stickyHeader
                stripe="even"
                variant="soft"
              >
                <thead>
                  <tr>
                    <th style={{ width: 50 }}>序号</th>
                    <th>名称</th>
                    <th style={{ width: 110 }}>当前版本</th>
                    <th style={{ width: 110 }}>最新版本</th>
                    <th style={{ width: 160 }}>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {updatable.map((it, i) => {
                    const isUpdating = updatingKey === it.nexusKey
                    const isUpdated = !!updatedKeys[it.nexusKey]
                    return (
                      <tr key={`${it.nexusKey}-${i}`}>
                        <td>{i + 1}</td>
                        <td>{it.name}</td>
                        <td>
                          <Typography level={'body-sm'}>
                            v{it.currentVersion}
                          </Typography>
                        </td>
                        <td>
                          <Typography
                            level={'body-sm'}
                            color={'success'}
                          >
                            v{it.latestVersion}
                          </Typography>
                        </td>
                        <td>
                          <Tooltip
                            title={'访问 Nexus 页面'}
                            placement={'top'}
                            variant={'outlined'}
                          >
                            <IconButton
                              variant={'plain'}
                              size={'sm'}
                              color={'primary'}
                              onClick={() => {
                                BrowserOpenURL(
                                  `https://www.nexusmods.com/stardewvalley/mods/${it.nexusKey}`,
                                )
                              }}
                            >
                              <LinkIcon />
                            </IconButton>
                          </Tooltip>

                          <Tooltip
                            title={'查看更新日志'}
                            placement={'top'}
                            variant={'outlined'}
                          >
                            <IconButton
                              variant={'plain'}
                              size={'sm'}
                              color={'primary'}
                              onClick={() => onViewChangelog(it)}
                            >
                              <AccessTimeOutlinedIcon />
                            </IconButton>
                          </Tooltip>

                          <Tooltip
                            title={isUpdated ? '已更新' : '更新'}
                            placement={'top'}
                            variant={'outlined'}
                          >
                            <IconButton
                              variant={'plain'}
                              size={'sm'}
                              color={isUpdated ? 'success' : 'primary'}
                              loading={isUpdating}
                              disabled={
                                isUpdating || isUpdated || updatingKey !== 0
                              }
                              onClick={() => onClickUpdate(it)}
                            >
                              <CloudDownloadOutlinedIcon />
                            </IconButton>
                          </Tooltip>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </Table>
            </div>
          )}

          {result && !loading && errors.length > 0 && (
            <Typography
              level={'body-xs'}
              color={'warning'}
              sx={{ marginTop: 1 }}
            >
              {errors.length} 个 Mod 检查失败：
              {errors.map((e) => `${e.name}（${e.error}）`).join('；')}
            </Typography>
          )}

          <div
            style={{
              marginTop: 16,
              display: 'flex',
              justifyContent: 'space-between',
              gap: 8,
            }}
          >
            <Button
              size={'sm'}
              variant={'plain'}
              color={'neutral'}
              startDecorator={<HistoryOutlinedIcon />}
              onClick={onOpenRollback}
            >
              查看回滚记录
            </Button>

            <Button
              size={'sm'}
              variant={'soft'}
              color={'primary'}
              onClick={onRefresh}
              loading={loading}
            >
              重新检查
            </Button>
          </div>
        </div>
      </Modal>

      <Modal
        title={`更新日志${changelog.modName ? ` - ${changelog.modName}` : ''}`}
        open={changelog.open}
        onClose={() => setChangelog(initialChangelog)}
      >
        <div style={{ width: 560, maxHeight: 480, overflowY: 'auto' }}>
          {changelog.loading ? (
            <div style={{ textAlign: 'center', padding: 24 }}>
              <CircularProgress />
              <div style={{ marginTop: 12 }}>正在获取更新日志...</div>
            </div>
          ) : changelog.message ? (
            <Typography
              level={'body-md'}
              color={'danger'}
            >
              {changelog.message}
            </Typography>
          ) : changelog.entries.length === 0 ? (
            <div
              style={{ textAlign: 'center', padding: 24, fontStyle: 'italic' }}
            >
              暂无更新日志
            </div>
          ) : (
            changelog.entries.map((entry) => (
              <div
                key={entry.version}
                style={{ marginBottom: 16 }}
              >
                <Typography
                  level={'title-sm'}
                  color={'primary'}
                >
                  v{entry.version}
                </Typography>
                <ul style={{ margin: '4px 0 0', paddingLeft: 20 }}>
                  {entry.changes.map((change, i) => (
                    <li key={i}>
                      <Typography level={'body-sm'}>{change}</Typography>
                    </li>
                  ))}
                </ul>
              </div>
            ))
          )}
        </div>
      </Modal>

      <Modal
        title={'保留原配置文件？'}
        open={keepConfig.open}
        onClose={() => setKeepConfig({ open: false })}
      >
        <div style={{ width: 440 }}>
          <Typography level={'body-md'}>
            检测到 {keepConfig.item?.name} 目录下存在 config.json，
            是否将其保留并应用到新版 Mod？
          </Typography>
          <Typography
            level={'body-xs'}
            color={'neutral'}
            sx={{ marginTop: 1 }}
          >
            选择「保留」会在新版安装完成后把旧版 config.json 拷贝到新目录。
          </Typography>
          <div
            style={{
              marginTop: 16,
              display: 'flex',
              justifyContent: 'flex-end',
              gap: 8,
            }}
          >
            <Button
              size={'sm'}
              variant={'plain'}
              color={'neutral'}
              onClick={() => setKeepConfig({ open: false })}
            >
              取消
            </Button>
            <Button
              size={'sm'}
              variant={'soft'}
              color={'neutral'}
              onClick={() =>
                keepConfig.item && runUpdate(keepConfig.item, false)
              }
            >
              不保留
            </Button>
            <Button
              size={'sm'}
              variant={'soft'}
              color={'primary'}
              onClick={() =>
                keepConfig.item && runUpdate(keepConfig.item, true)
              }
            >
              保留
            </Button>
          </div>
        </div>
      </Modal>

      <Modal
        title={'回滚 Mod 更新'}
        open={rollbackOpen}
        onClose={() => setRollbackOpen(false)}
      >
        <div style={{ width: 720, minHeight: 120 }}>
          {rollbackLoading ? (
            <div style={{ textAlign: 'center', padding: 24 }}>
              <CircularProgress />
              <div style={{ marginTop: 12 }}>正在加载备份...</div>
            </div>
          ) : backups.length === 0 ? (
            <div
              style={{ textAlign: 'center', padding: 24, fontStyle: 'italic' }}
            >
              暂无更新备份
            </div>
          ) : (
            <div style={{ maxHeight: 420, overflowY: 'auto' }}>
              <Table
                stickyHeader
                stripe="even"
                variant="soft"
              >
                <thead>
                  <tr>
                    <th>Mod</th>
                    <th style={{ width: 130 }}>旧版 → 新版</th>
                    <th style={{ width: 90 }}>大小</th>
                    <th style={{ width: 160 }}>备份时间</th>
                    <th style={{ width: 130 }}>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {backups.map((b) => (
                    <tr key={b.name}>
                      <td>
                        <Typography level={'body-sm'}>
                          {b.modName || b.originalDir || b.name}
                        </Typography>
                        {b.keptConfig && (
                          <Chip
                            size={'sm'}
                            variant={'soft'}
                            color={'primary'}
                            sx={{ marginTop: 0.5 }}
                          >
                            含原配置
                          </Chip>
                        )}
                      </td>
                      <td>
                        <Typography level={'body-sm'}>
                          v{b.oldVersion || '-'} → v{b.newVersion || '-'}
                        </Typography>
                      </td>
                      <td>
                        <Typography level={'body-sm'}>
                          {formatBytes(b.size || 0)}
                        </Typography>
                      </td>
                      <td>
                        <Typography level={'body-sm'}>
                          {formatTime(b.createTime)}
                        </Typography>
                      </td>
                      <td>
                        <Tooltip
                          title={'回滚到旧版'}
                          placement={'top'}
                          variant={'outlined'}
                        >
                          <IconButton
                            variant={'plain'}
                            size={'sm'}
                            color={'primary'}
                            loading={busyBackup === b.name}
                            disabled={!!busyBackup && busyBackup !== b.name}
                            onClick={() => onRollback(b)}
                          >
                            <RestoreOutlinedIcon />
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
                            disabled={!!busyBackup && busyBackup !== b.name}
                            onClick={() => onRemoveBackup(b)}
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
          )}

          <div
            style={{
              marginTop: 16,
              display: 'flex',
              justifyContent: 'flex-end',
              gap: 8,
            }}
          >
            <Button
              size={'sm'}
              variant={'soft'}
              color={'primary'}
              onClick={loadBackups}
              loading={rollbackLoading}
            >
              刷新
            </Button>
          </div>
        </div>
      </Modal>
    </React.Fragment>
  )
}

export default ModUpdatesModal
