import React, { useState } from 'react'
import {
  Button,
  CircularProgress,
  Table,
  Tooltip,
  Typography,
  IconButton,
} from '@mui/joy'
import LinkIcon from '@mui/icons-material/Link'
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined'
import { Modal } from '@/components'
import { BrowserOpenURL, ViewModChangelog } from '@/utils'
import {
  CheckForUpdatesResultType,
  ModChangelogEntryType,
  ModUpdateInfoType,
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

const initialChangelog: ChangelogState = {
  open: false,
  loading: false,
  modName: '',
  entries: [],
  message: '',
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
        <div style={{ width: 720, minHeight: 120 }}>
          {loading ? (
            <div style={{ textAlign: 'center', padding: 24 }}>
              <CircularProgress />
              <div style={{ marginTop: 12 }}>
                正在通过 Nexus Mods API 检查更新...
              </div>
            </div>
          ) : !result ? (
            <div style={{ textAlign: 'center', padding: 24, fontStyle: 'italic' }}>暂无数据</div>
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
                    <th style={{ width: 120 }}>当前版本</th>
                    <th style={{ width: 120 }}>最新版本</th>
                    <th style={{ width: 80 }}>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {updatable.map((it, i) => (
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
                          title={'查看更新详情'}
                          placement={'top'}
                          variant={'outlined'}
                        >
                          <IconButton
                            variant={'plain'}
                            size={'sm'}
                            color={'primary'}
                            onClick={() => onViewChangelog(it)}
                          >
                            <InfoOutlinedIcon />
                          </IconButton>
                        </Tooltip>
                      </td>
                    </tr>
                  ))}
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
              justifyContent: 'flex-end',
              gap: 8,
            }}
          >
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
            <div style={{ textAlign: 'center', padding: 24, fontStyle: 'italic' }}>暂无更新日志</div>
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
    </React.Fragment>
  )
}

export default ModUpdatesModal
