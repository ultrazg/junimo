import React from 'react'
import { Modal, ModalClose, ModalDialog, Typography } from '@mui/joy'

type IProps = {
  title?: string
  open: boolean
  onClose: () => void
  children: React.ReactNode
}

const IModal: React.FC<IProps> = ({ title, open, onClose, children }) => {
  return (
    <Modal
      open={open}
      onClose={onClose}
      hideBackdrop
    >
      <ModalDialog>
        <ModalClose
          variant={'plain'}
          size={'sm'}
        />

        {title && (
          <Typography
            component="h3"
            id="modal-title"
            level="h4"
            textColor="inherit"
            sx={{ fontWeight: 'lg', userSelect: 'none' }}
          >
            {title}
          </Typography>
        )}

        {children}
      </ModalDialog>
    </Modal>
  )
}

export default IModal
