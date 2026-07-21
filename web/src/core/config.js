/**
 * 网站配置文件
 * Modified for ai-pcdn on 2026-07-21: 已获商业授权，移除上游品牌与版权横幅。
 */
import packageInfo from '../../package.json'

const greenText = (text) => `\x1b[32m${text}\x1b[0m`

export const config = {
  appName: 'ai-pcdn',
  showViteLogo: true,
  keepAliveTabs: false,
  logs: []
}

export const viteLogo = () => {
  if (!config.showViteLogo) return

  console.log(greenText('> 欢迎使用 ai-pcdn'))
  console.log(greenText(`> 当前版本:v${packageInfo.version}`))
}

export default config
