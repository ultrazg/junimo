import { SaveGamePath } from 'wailsjs/go/backend/App'

export const onSaveGamePath = async () => {
  const { success, message, path } = await SaveGamePath()

  return {
    success,
    message,
    path,
  }
}
