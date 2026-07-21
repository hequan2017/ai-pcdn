<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="search">
        <el-form-item label="账期"><el-input v-model="search.period" style="width: 110px" /></el-form-item>
        <el-form-item label="大厂">
          <el-select v-model="search.platform" clearable style="width: 110px"><el-option label="抖音" value="douyin" /><el-option label="腾讯" value="tencent" /></el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="search.status" clearable style="width: 110px"><el-option label="待核对" value="pending" /><el-option label="已核对" value="matched" /><el-option label="有差异" value="diff" /></el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="load">查询</el-button>
          <el-button type="success" @click="showImport = true">导入结算单</el-button>
          <el-button @click="recheck">重新核对</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div v-if="summary" class="summary-cards">
        <div class="card"><div class="label">应收总额</div><div class="value">¥{{ (summary.totalRevenue || 0).toFixed(2) }}</div></div>
        <div class="card"><div class="label">已核对</div><div class="value">{{ summary.matchedCount }}</div></div>
        <div class="card"><div class="label">有差异</div><div class="value warn">{{ summary.diffCount }}</div></div>
        <div class="card"><div class="label">待核对</div><div class="value">{{ summary.pendingCount }}</div></div>
      </div>
      <el-table :data="list">
        <el-table-column label="账期" prop="period" width="100" />
        <el-table-column label="大厂" prop="platform" width="90" />
        <el-table-column label="节点SN" prop="nodeSn" min-width="160" />
        <el-table-column label="收入" width="110"><template #default="s">¥{{ (s.row.revenue || 0).toFixed(2) }}</template></el-table-column>
        <el-table-column label="大厂流量Mbps" width="120"><template #default="s">{{ s.row.trafficBps ? (s.row.trafficBps / 1e6).toFixed(2) : '-' }}</template></el-table-column>
        <el-table-column label="自采集Mbps" width="120"><template #default="s">{{ s.row.ourTrafficBps ? (s.row.ourTrafficBps / 1e6).toFixed(2) : '-' }}</template></el-table-column>
        <el-table-column label="差异%" width="90"><template #default="s">{{ s.row.diffPercent ? (s.row.diffPercent * 100).toFixed(1) + '%' : '-' }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="s"><el-tag :type="statusType(s.row.status)">{{ statusLabel(s.row.status) }}</el-tag></template></el-table-column>
        <el-table-column label="操作" width="90"><template #default="s"><el-button link type="primary" @click="del(s.row)">删除</el-button></template></el-table-column>
      </el-table>
      <div class="gva-pagination"><el-pagination :current-page="page" :page-size="pageSize" :total="total" @current-change="(v) => { page = v; load() }" /></div>
    </div>

    <el-dialog v-model="showImport" title="导入结算单" width="520">
      <el-form :model="form" label-width="110px">
        <el-form-item label="账期"><el-input v-model="form.period" placeholder="2026-07" /></el-form-item>
        <el-form-item label="大厂"><el-select v-model="form.platform"><el-option label="抖音" value="douyin" /><el-option label="腾讯" value="tencent" /></el-select></el-form-item>
        <el-form-item label="节点SN"><el-input v-model="form.nodeSn" /></el-form-item>
        <el-form-item label="收入(元)"><el-input-number v-model="form.revenue" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="大厂流量Mbps"><el-input-number v-model="trafficMbps" :min="0" /></el-form-item>
      </el-form>
      <template #footer><el-button type="primary" :loading="imp" @click="doImport">导入</el-button><el-button @click="showImport = false">取消</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { ref, reactive } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { importSettlement, recheckSettlement, getSettlementList, getRevenueSummary, deleteSettlement } from '@/plugin/pcdn/api/settlement'

  defineOptions({ name: 'PcdnSettlement' })

  const list = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const search = reactive({ period: '', platform: '', status: '' })
  const summary = ref(null)

  const load = async () => {
    const r = await getSettlementList({ page: page.value, pageSize: pageSize.value, ...search })
    if (r.code === 0) { list.value = r.data.list || []; total.value = r.data.total }
    const s = await getRevenueSummary({ period: search.period, platform: search.platform })
    if (s.code === 0) summary.value = s.data
  }
  load()

  const statusType = (s) => ({ pending: 'info', matched: 'success', diff: 'danger' }[s] || 'info')
  const statusLabel = (s) => ({ pending: '待核对', matched: '已核对', diff: '有差异' }[s] || s)

  const showImport = ref(false)
  const imp = ref(false)
  const form = reactive({ period: '', platform: 'douyin', nodeSn: '', revenue: 0 })
  const trafficMbps = ref(0)
  const doImport = async () => {
    imp.value = true
    try { const r = await importSettlement({ ...form, trafficBps: Math.round(trafficMbps.value * 1e6) }); if (r.code === 0) { ElMessage.success('已导入'); showImport.value = false; load() } } finally { imp.value = false }
  }

  const recheck = async () => {
    if (!search.period) { ElMessage.warning('请先填账期'); return }
    const r = await recheckSettlement({ period: search.period })
    if (r.code === 0) { ElMessage.success('已核对 ' + r.data.count + ' 条'); load() }
  }
  const del = (row) => ElMessageBox.confirm('删除该结算单?', '提示', { type: 'warning' }).then(async () => { const r = await deleteSettlement({ id: row.ID }); if (r.code === 0) { ElMessage.success('已删除'); load() } })
</script>

<style scoped>
  .summary-cards { display: flex; gap: 16px; margin-bottom: 16px; }
  .card { flex: 1; padding: 16px; border: 1px solid #eee; border-radius: 8px; text-align: center; }
  .card .label { color: #6b7280; font-size: 13px; }
  .card .value { font-size: 22px; font-weight: 700; margin-top: 4px; }
  .card .value.warn { color: #f56c6c; }
</style>
