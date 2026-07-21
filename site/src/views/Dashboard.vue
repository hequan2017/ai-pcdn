<template>
  <div class="dash">
    <header class="nav">
      <div class="logo">PCDN 控制台</div>
      <div class="right">
        <span class="user">{{ user?.nickName || user?.userName || '用户' }}</span>
        <el-button link @click="logout">退出</el-button>
      </div>
    </header>
    <main class="content">
      <div class="toolbar">
        <h3>我的节点</h3>
        <el-button type="primary" @click="showAdd = true">+ 添加节点（上机）</el-button>
      </div>
      <el-table :data="nodes" border>
        <el-table-column label="节点SN" prop="nodeSn" min-width="200" />
        <el-table-column label="地域" prop="region" width="120" />
        <el-table-column label="运营商" prop="isp" width="90" />
        <el-table-column label="大厂" prop="platform" width="90" />
        <el-table-column label="状态" width="100">
          <template #default="s">
            <el-tag :type="statusType(s.row.status)">{{ statusLabel(s.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="s">
            <el-button link type="primary" @click="openTraffic(s.row)">流量</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="section-title">我的账单</div>
      <el-table :data="bills" border size="small">
        <el-table-column label="账期" prop="period" width="100" />
        <el-table-column label="节点数" prop="nodeCount" width="80" />
        <el-table-column label="应付总额" width="120"><template #default="s">¥{{ (s.row.totalAmount || 0).toFixed(2) }}</template></el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="s">
            <el-tag :type="({ draft: 'warning', approved: 'primary', paid: 'success', rejected: 'info' }[s.row.status])">{{ ({ draft: '待审核', approved: '已审核', paid: '已付款', rejected: '已驳回' }[s.row.status]) || s.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="付款时间" prop="paidAt" min-width="160"><template #default="s">{{ s.row.paidAt || '-' }}</template></el-table-column>
      </el-table>
    </main>

    <!-- 添加节点 / 上机 -->
    <el-dialog v-model="showAdd" title="添加节点（自助上机）" width="560">
      <el-form :model="addForm" label-width="80px" v-if="!credential">
        <el-form-item label="地域"><el-input v-model="addForm.region" placeholder="省/市" /></el-form-item>
        <el-form-item label="运营商">
          <el-select v-model="addForm.isp">
            <el-option label="电信" value="电信" /><el-option label="联通" value="联通" /><el-option label="移动" value="移动" />
          </el-select>
        </el-form-item>
        <el-form-item label="接入大厂">
          <el-select v-model="addForm.platform">
            <el-option label="抖音" value="douyin" /><el-option label="腾讯" value="tencent" />
          </el-select>
        </el-form-item>
      </el-form>
      <div v-else class="cred">
        <el-alert type="warning" :closable="false" title="凭证仅展示一次，请立即复制保存！" show-icon />
        <div class="line">节点SN：<code>{{ credential.nodeSn }}</code></div>
        <div class="line">Token：<code>{{ credential.token }}</code></div>
        <div class="line">一键安装命令：</div>
        <pre>{{ credential.installScript }}</pre>
        <el-button size="small" @click="copy(credential.installScript)">复制安装命令</el-button>
      </div>
      <template #footer>
        <el-button v-if="!credential" type="primary" :loading="adding" @click="onAdd">生成凭证</el-button>
        <el-button @click="closeAdd">{{ credential ? '完成' : '取消' }}</el-button>
      </template>
    </el-dialog>

    <!-- 流量 -->
    <el-drawer v-model="trafficVisible" size="680" :title="`流量详情 · ${current?.nodeSn || ''}`">
      <el-radio-group v-model="range" @change="loadTraffic">
        <el-radio-button label="6h">近6小时</el-radio-button>
        <el-radio-button label="24h">近24小时</el-radio-button>
        <el-radio-button label="7d">近7天</el-radio-button>
      </el-radio-group>
      <el-table :data="traffic" size="small" style="margin-top: 12px">
        <el-table-column label="时间" prop="windowStart" min-width="160" />
        <el-table-column label="网卡" prop="ifaceName" width="100" />
        <el-table-column label="下行Mbps" width="110">
          <template #default="s">{{ (s.row.rxMaxBps / 1e6).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="上行Mbps" width="110">
          <template #default="s">{{ (s.row.txMaxBps / 1e6).toFixed(2) }}</template>
        </el-table-column>
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup>
  import { ref, reactive, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import { ElMessage } from 'element-plus'
  import { myNodes, myNodeTraffic, addNode, myBills } from '../api/portal'
  import { getUser, clearAuth } from '../store/auth'

  const router = useRouter()
  const user = ref(getUser())
  const nodes = ref([])

  const loadNodes = async () => {
    const r = await myNodes({ page: 1, pageSize: 100 })
    nodes.value = r.data.list || []
  }
  const bills = ref([])
  const loadBills = async () => {
    const r = await myBills({ page: 1, pageSize: 20 })
    bills.value = r.data.list || []
  }
  onMounted(() => { loadNodes(); loadBills() })

  const logout = () => { clearAuth(); router.push('/login') }

  const statusType = (s) => ({ online: 'success', offline: 'danger', pending: 'info', abnormal: 'warning', disabled: 'info' }[s] || 'info')
  const statusLabel = (s) => ({ online: '在线', offline: '离线', pending: '待上机', abnormal: '异常', disabled: '已停用' }[s] || s)

  // 添加节点
  const showAdd = ref(false)
  const adding = ref(false)
  const addForm = reactive({ region: '', isp: '电信', platform: 'douyin' })
  const credential = ref(null)
  const onAdd = async () => {
    adding.value = true
    try {
      const r = await addNode(addForm)
      credential.value = r.data
      ElMessage.success('凭证已生成')
      loadNodes()
    } finally {
      adding.value = false
    }
  }
  const closeAdd = () => { showAdd.value = false; credential.value = null; addForm.region = '' }
  const copy = (t) => { navigator.clipboard.writeText(t); ElMessage.success('已复制') }

  // 流量
  const trafficVisible = ref(false)
  const current = ref(null)
  const traffic = ref([])
  const range = ref('24h')
  const openTraffic = async (row) => { current.value = row; trafficVisible.value = true; await loadTraffic() }
  const loadTraffic = async () => {
    const end = new Date()
    const ms = { '6h': 6 * 3600e3, '24h': 24 * 3600e3, '7d': 7 * 24 * 3600e3 }[range.value]
    const start = new Date(end.getTime() - ms)
    const r = await myNodeTraffic({ nodeId: current.value.ID, start: fmt(start), end: fmt(end) })
    traffic.value = r.data || []
  }
  const fmt = (d) => {
    const p = (n) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
  }
</script>

<style scoped>
  .dash { font-family: system-ui, sans-serif; }
  .nav { display: flex; justify-content: space-between; align-items: center; padding: 14px 32px; border-bottom: 1px solid #eee; }
  .logo { font-weight: 700; }
  .right { display: flex; align-items: center; gap: 12px; }
  .content { padding: 24px 32px; }
  .toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
  .cred { line-height: 2; }
  .cred .line { margin: 4px 0; }
  .cred code { background: #f3f4f6; padding: 2px 6px; border-radius: 4px; }
  .cred pre { background: #1f2937; color: #fff; padding: 12px; border-radius: 8px; overflow-x: auto; }
</style>
