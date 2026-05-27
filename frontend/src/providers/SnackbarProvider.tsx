import React, {
  useState,
  useCallback,
  ReactNode,
  forwardRef,
  useImperativeHandle,
} from 'react'
import Snackbar from '@mui/joy/Snackbar'
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined'
import CheckCircleOutlinedIcon from '@mui/icons-material/CheckCircleOutlined'
import ErrorOutlineOutlinedIcon from '@mui/icons-material/ErrorOutlineOutlined'
import HighlightOffOutlinedIcon from '@mui/icons-material/HighlightOffOutlined'
import { SnackbarOptions } from '@/types'

export type SnackbarHandler = {
  show: (message: string, options?: SnackbarOptions) => void
}

export const snackbar: SnackbarHandler = {
  show: () => {
    throw new Error('SnackbarProvider not mounted yet.')
  },
}

type SnackbarProviderProps = {
  children: ReactNode
}

const renderSnackbarIcon = (color: SnackbarOptions['color']) => {
  switch (color) {
    case 'primary':
      return <InfoOutlinedIcon />
    case 'neutral':
      return <InfoOutlinedIcon />
    case 'danger':
      return <HighlightOffOutlinedIcon />
    case 'success':
      return <CheckCircleOutlinedIcon />
    case 'warning':
      return <ErrorOutlineOutlinedIcon />
    default:
      return null
  }
}

const SnackbarProviderInner = forwardRef<
  SnackbarHandler,
  SnackbarProviderProps
>(function SnackbarProviderInner({ children }, ref) {
  const [open, setOpen] = useState(false)
  const [message, setMessage] = useState('')
  const [options, setOptions] = useState<SnackbarOptions>({
    color: 'neutral',
    variant: 'outlined',
    autoHideDuration: 3000,
    showIcon: true,
  })

  const show = useCallback((msg: string, opts?: SnackbarOptions) => {
    setMessage(msg)
    setOptions({
      color: opts?.color || 'neutral',
      variant: opts?.variant || 'outlined',
      autoHideDuration: opts?.autoHideDuration || 3000,
      showIcon: opts?.showIcon || false,
    })
    setOpen(true)
  }, [])

  const handleClose = useCallback(() => setOpen(false), [])

  useImperativeHandle(ref, () => ({ show }))

  return (
    <>
      {children}
      <Snackbar
        open={open}
        color={options.color}
        variant={options.variant}
        autoHideDuration={options.autoHideDuration}
        onClose={handleClose}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
        size={'sm'}
      >
        {options.showIcon && renderSnackbarIcon(options.color)}
        <pre style={{ fontFamily: 'inherit', padding: 0, margin: 0 }}>
          {message}
        </pre>
      </Snackbar>
    </>
  )
})

export const SnackbarProvider = ({ children }: SnackbarProviderProps) => {
  const ref = React.useRef<SnackbarHandler>(null)

  React.useEffect(() => {
    if (ref.current) {
      snackbar.show = ref.current.show
    }
  }, [])

  return <SnackbarProviderInner ref={ref}>{children}</SnackbarProviderInner>
}
