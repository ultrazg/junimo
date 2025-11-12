import React from 'react'
import { Modal } from '@/components'
import { Typography, Button, Link, Divider } from '@mui/joy'
import { APP_NAME, APP_VERSION, BrowserOpenURL } from '@/utils'
import styles from './index.module.scss'
import LoopTwoToneIcon from '@mui/icons-material/LoopTwoTone'
import BalanceTwoToneIcon from '@mui/icons-material/BalanceTwoTone'
import GitHubIcon from '@mui/icons-material/GitHub'
import BugReportTwoToneIcon from '@mui/icons-material/BugReportTwoTone'

type IProps = {
  open: boolean
  onClose: () => void
}

const About: React.FC<IProps> = ({ open, onClose }) => {
  return (
    <React.Fragment>
      <Modal
        open={open}
        onClose={onClose}
        title={'关于'}
      >
        <div className={styles['about-wrapper']}>
          <div className={styles['app-logo']}>
            <div className={styles['logo']} />
          </div>

          <div className={styles['app-name']}>
            <Typography level="h4">{APP_NAME}</Typography>
          </div>

          <div>
            <Typography>v{APP_VERSION}</Typography>
          </div>

          <div>
            <Typography>
              {APP_NAME} 是星露谷物语（
              <Link
                onClick={() => BrowserOpenURL('https://www.stardewvalley.net/')}
              >
                Stardew Valley
              </Link>
              ）Mod 管理软件
            </Typography>
            <Typography>
              使用{' '}
              <Link onClick={() => BrowserOpenURL('https://wails.io/')}>
                wails
              </Link>{' '}
              +{' '}
              <Link onClick={() => BrowserOpenURL('https://react.dev/')}>
                react
              </Link>{' '}
              构建
            </Typography>
          </div>

          <div className={styles['check-for-update-button']}>
            <Button
              size={'sm'}
              variant={'soft'}
            >
              <LoopTwoToneIcon />
              检查更新...
            </Button>
          </div>

          <div className={styles['extra-button']}>
            <Link>
              <GitHubIcon />
              GitHub
            </Link>

            <Divider
              orientation={'vertical'}
              style={{ margin: '0 8px' }}
            />

            <Link>
              <BalanceTwoToneIcon />
              GPL-3.0 License
            </Link>

            <Divider
              orientation={'vertical'}
              style={{ margin: '0 8px' }}
            />

            <Link>
              <BugReportTwoToneIcon />
              Bug report
            </Link>
          </div>
        </div>
      </Modal>
    </React.Fragment>
  )
}

export default About
