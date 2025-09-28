import { version } from '../../package.json'
import { Environment } from 'wailsjs/runtime'
import { envType } from '@/types/env'

const APP_VERSION: string = version
const APP_NAME: string = 'Junimo'
const ENV: envType = await Environment()
const APP_INFO = {
  name: APP_NAME,
  version: APP_VERSION,
  env: ENV,
}

export { APP_INFO }
