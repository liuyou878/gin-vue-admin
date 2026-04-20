/**
 * 网站配置文件
 */
import packageInfo from '../../package.json'

const greenText = (text) => `\x1b[32m${text}\x1b[0m`

export const config = {
  appName: 'Alpha',
  showViteLogo: false,
  keepAliveTabs: false,
  logs: []
}

export const viteLogo = (env) => {
  console.log(env);

}

export default config
