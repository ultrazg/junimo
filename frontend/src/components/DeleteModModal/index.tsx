import React from 'react'
import { Typography, Button } from '@mui/joy'
import { Modal } from '@/components'
import WarningTwoToneIcon from '@mui/icons-material/WarningTwoTone'
import { RemoveModDir } from '@/utils'

type IProps = {
  path: string
  open: boolean
  onClose: () => void
}

const DeleteModModal: React.FC<IProps> = ({ path, open, onClose }) => {
  const onRemoveMod = () => {
    RemoveModDir(path).then(() => {
      onClose()
    })
  }

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
          确定要删除该 Mod 吗？
        </Typography>

        <div style={{ textAlign: 'right' }}>
          <Button
            size={'sm'}
            variant={'soft'}
            color={'danger'}
            onClick={() => {
              onRemoveMod()
            }}
          >
            <WarningTwoToneIcon style={{ marginRight: 6 }} />
            确定
          </Button>
        </div>
      </Modal>
    </React.Fragment>
  )
}

export default DeleteModModal
