/*
 * ai-pcdn web 框架初始化
 * Modified for ai-pcdn on 2026-07-21: 已获商业授权，移除上游品牌与版权横幅。
 */
import { register } from './global'
import packageInfo from '../../package.json'

export default {
  install: (app) => {
    register(app)
    console.log(`ai-pcdn v${packageInfo.version}`)
  }
}
