import React from 'react'
import {
  Button,
  CircularProgress,
  Table,
  Tooltip,
  Typography,
  IconButton,
} from '@mui/joy'
import LinkIcon from '@mui/icons-material/Link'
import { Modal } from '@/components'
import { BrowserOpenURL } from '@/utils'
import { CheckForUpdatesResultType } from '@/types'

type IProps = {
  open: boolean
  loading: boolean
  result?: CheckForUpdatesResultType
  onClose: () => void
  onRefresh: () => void
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

  return (
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
          <div style={{ textAlign: 'center', padding: 24 }}>暂无数据</div>
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
            variant={'plain'}
            color={'neutral'}
            onClick={onClose}
            disabled={loading}
          >
            关闭
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
  )
}

export default ModUpdatesModal
