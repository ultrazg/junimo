import { SaveGamePath } from 'wailsjs/go/backend/App'

export const onSaveGamePath = async () => {
  const { success, message, path } = await SaveGamePath()

  return {
    success,
    message,
    path,
  }
}

export const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  const value = bytes / Math.pow(k, i)
  return `${value.toFixed(2)} ${sizes[i]}`
}
