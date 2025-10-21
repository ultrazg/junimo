import React from 'react'
import { Typography, Button } from '@mui/joy'
import { Modal } from '@/components'
import WarningTwoToneIcon from '@mui/icons-material/WarningTwoTone'

type IProps = {
  open: boolean
  loading: boolean
  onOk: () => void
  onClose: () => void
}

const DeleteModal: React.FC<IProps> = ({ open, loading, onOk, onClose }) => {
  return (
    <React.Fragment>
      <Modal
        title={'警告'}
        open={open}
        onClose={onClose}
      >
        <Typography
          level={'body-md'}
          color={'danger'}
        >
          确定要执行删除操作吗？
        </Typography>

        <div style={{ textAlign: 'right' }}>
          <Button
            size={'sm'}
            variant={'soft'}
            color={'danger'}
            onClick={() => {
              onOk()
            }}
            loading={loading}
          >
            <WarningTwoToneIcon style={{ marginRight: 6 }} />
            确定
          </Button>
        </div>
      </Modal>
    </React.Fragment>
  )
}

export default DeleteModal
