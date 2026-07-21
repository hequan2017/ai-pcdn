<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="节点SN">
          <el-input v-model="searchInfo.nodeSn" placeholder="节点SN" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
            <el-option label="待上机" value="pending" />
            <el-option label="异常" value="abnormal" />
            <el-option label="已停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="归属">
          <el-input v-model="searchInfo.ownerName" placeholder="归属用户" clearable />
        </el-form-item>
        <el-form-item label="大厂">
          <el-select v-model="searchInfo.platform" placeholder="全部" clearable style="width: 120px">
            <el-option label="抖音" value="douyin" />
            <el-option label="腾讯" value="tencent" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button icon="delete" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="节点SN" prop="nodeSn" min-width="180" />
        <el-table-column label="归属" prop="ownerName" width="120" />
        <el-table-column label="地域" prop="region" width="120" />
        <el-table-column label="运营商" prop="isp" width="90" />
        <el-table-column label="大厂" prop="platform" width="90" />
        <el-table-column label="状态" prop="status" width="100">
          <template #default="scope">
            <el-tag :type="statusType(scope.row.status)">{{ statusLabel(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="最后心跳" prop="lastHeartbeatAt" width="170">
          <template #default="scope">{{ scope.row.lastHeartbeatAt ? formatDate(scope.row.lastHeartbeatAt) : '-' }}</template>
        </el-table-column>
        <el-table-column label="agent版本" prop="agentVersion" width="100" />
        <el-table-column label="操作" fixed="right" width="240">
          <template #default="scope">
            <el-button type="primary" link icon="view" @click="openTraffic(scope.row)">流量</el-button>
            <el-button type="primary" link icon="edit" @click="openEdit(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 编辑抽屉 -->
    <el-drawer v-model="editVisible" size="560" :before-close="closeEdit" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-base">编辑节点</span>
          <div>
            <el-button type="primary" @click="enterEdit">确定</el-button>
            <el-button @click="closeEdit">取消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="editFormRef" :model="formData" label-width="100px">
        <el-form-item label="节点SN"><el-input v-model="formData.nodeSn" disabled /></el-form-item>
        <el-form-item label="归属用户"><el-input v-model="formData.ownerName" /></el-form-item>
        <el-form-item label="地域"><el-input v-model="formData.region" /></el-form-item>
        <el-form-item label="运营商">
          <el-select v-model="formData.isp" style="width: 100%">
            <el-option label="电信" value="电信" /><el-option label="联通" value="联通" /><el-option label="移动" value="移动" />
          </el-select>
        </el-form-item>
        <el-form-item label="接入大厂">
          <el-select v-model="formData.platform" style="width: 100%">
            <el-option label="抖音" value="douyin" /><el-option label="腾讯" value="tencent" /><el-option label="无" value="" />
          </el-select>
        </el-form-item>
        <el-form-item label="计费模式">
          <el-select v-model="formData.billingMode" style="width: 100%">
            <el-option label="包月" value="monthly" /><el-option label="95计费" value="p95" />
          </el-select>
        </el-form-item>
        <el-form-item label="包月价"><el-input-number v-model="formData.monthlyPrice" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="formData.status" style="width: 100%">
            <el-option label="在线" value="online" /><el-option label="离线" value="offline" />
            <el-option label="异常" value="abnormal" /><el-option label="已停用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" /></el-form-item>
      </el-form>
    </el-drawer>

    <!-- 流量与95值抽屉 -->
    <el-drawer v-model="trafficVisible" size="780" :title="`节点流量 · ${currentNode?.nodeSn || ''}`">
      <div class="mb-3 flex items-center gap-2">
        <el-radio-group v-model="trafficRange" @change="loadTraffic">
          <el-radio-button label="6h">近6小时</el-radio-button>
          <el-radio-button label="24h">近24小时</el-radio-button>
          <el-radio-button label="7d">近7天</el-radio-button>
        </el-radio-group>
        <el-select v-model="trafficIface" placeholder="全部网卡" clearable style="width: 160px" @change="loadTraffic">
          <el-option v-for="f in ifaceOptions" :key="f" :label="f" :value="f" />
        </el-select>
      </div>
      <div ref="chartRef" style="width: 100%; height: 320px"></div>
      <div class="mt-4">
        <div class="mb-2 text-base font-bold">95值</div>
        <el-table :data="n95Data" size="small">
          <el-table-column label="周期类型" prop="periodType" width="100">
            <template #default="scope">{{ scope.row.periodType === 'month' ? '月' : '日' }}</template>
          </el-table-column>
          <el-table-column label="周期开始" prop="periodStart" width="170">
            <template #default="scope">{{ formatDate(scope.row.periodStart) }}</template>
          </el-table-column>
          <el-table-column label="下行95(Mbps)" width="130">
            <template #default="scope">{{ bpsToMbps(scope.row.rx95Bps) }}</template>
          </el-table-column>
          <el-table-column label="上行95(Mbps)" width="130">
            <template #default="scope">{{ bpsToMbps(scope.row.tx95Bps) }}</template>
          </el-table-column>
          <el-table-column label="合计95(Mbps)">
            <template #default="scope">{{ bpsToMbps(scope.row.combined95Bps) }}</template>
          </el-table-column>
          <el-table-column label="状态" prop="status" width="90">
            <template #default="scope">
              <el-tag :type="scope.row.status === 'frozen' ? 'success' : 'info'">{{ scope.row.status === 'frozen' ? '已冻结' : '滚动' }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
  import { getNodeList, updateNode, deleteNode, deleteNodeByIds, getNodeTraffic, getNode95 } from '@/plugin/pcdn/api/node'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, nextTick, onBeforeUnmount } from 'vue'
  import * as echarts from 'echarts'

  defineOptions({ name: 'PcdnNode' })

  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})
  const multipleSelection = ref([])

  const getTableData = async () => {
    const res = await getNodeList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total
    }
  }
  getTableData()

  const onSubmit = () => { page.value = 1; getTableData() }
  const onReset = () => { searchInfo.value = {}; getTableData() }
  const handleSizeChange = (v) => { pageSize.value = v; getTableData() }
  const handleCurrentChange = (v) => { page.value = v; getTableData() }
  const handleSelectionChange = (v) => { multipleSelection.value = v }

  const statusType = (s) => ({ online: 'success', offline: 'danger', pending: 'info', abnormal: 'warning', disabled: 'info' }[s] || 'info')
  const statusLabel = (s) => ({ online: '在线', offline: '离线', pending: '待上机', abnormal: '异常', disabled: '已停用' }[s] || s)
  const bpsToMbps = (bps) => (bps ? (bps / 1e6).toFixed(2) : '0.00')

  // 删除
  const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除该节点吗?', '提示', { type: 'warning' }).then(async () => {
      const res = await deleteNode({ id: row.ID })
      if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
    })
  }
  const onDelete = () => {
    if (!multipleSelection.value.length) return
    ElMessageBox.confirm('确定批量删除选中节点?', '提示', { type: 'warning' }).then(async () => {
      const ids = multipleSelection.value.map((i) => i.ID)
      const res = await deleteNodeByIds({ ids })
      if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
    })
  }

  // 编辑
  const editVisible = ref(false)
  const formData = ref({})
  const openEdit = (row) => { formData.value = { ...row }; editVisible.value = true }
  const closeEdit = () => { editVisible.value = false }
  const enterEdit = async () => {
    const res = await updateNode(formData.value)
    if (res.code === 0) { ElMessage.success('更新成功'); closeEdit(); getTableData() }
  }

  // 流量与95值
  const trafficVisible = ref(false)
  const currentNode = ref(null)
  const trafficRange = ref('24h')
  const trafficIface = ref('')
  const ifaceOptions = ref([])
  const chartRef = ref(null)
  const n95Data = ref([])
  let chart = null

  const rangeMs = { '6h': 6 * 3600e3, '24h': 24 * 3600e3, '7d': 7 * 24 * 3600e3 }
  const fmtDateTime = (d) => {
    const p = (n) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
  }

  const openTraffic = async (row) => {
    currentNode.value = row
    trafficVisible.value = true
    trafficIface.value = ''
    await loadTraffic()
    await load95()
  }

  const loadTraffic = async () => {
    if (!currentNode.value) return
    const end = new Date()
    const start = new Date(end.getTime() - rangeMs[trafficRange.value])
    const res = await getNodeTraffic({
      nodeId: currentNode.value.ID,
      iface: trafficIface.value,
      start: fmtDateTime(start),
      end: fmtDateTime(end)
    })
    const list = res.code === 0 ? res.data || [] : []
    const ifaces = [...new Set(list.map((p) => p.ifaceName))]
    ifaceOptions.value = ifaces
    drawChart(list)
  }

  const drawChart = (list) => {
    nextTick(() => {
      if (!chartRef.value) return
      chart = chart || echarts.init(chartRef.value)
      const times = [...new Set(list.map((p) => p.windowStart))].sort()
      const rx = times.map((t) => {
        const p = list.find((x) => x.windowStart === t)
        return p ? +(p.rxMaxBps / 1e6).toFixed(2) : 0
      })
      const tx = times.map((t) => {
        const p = list.find((x) => x.windowStart === t)
        return p ? +(p.txMaxBps / 1e6).toFixed(2) : 0
      })
      chart.setOption({
        tooltip: { trigger: 'axis' },
        legend: { data: ['下行Mbps', '上行Mbps'] },
        grid: { left: 50, right: 20, top: 40, bottom: 30 },
        xAxis: { type: 'category', data: times },
        yAxis: { type: 'value', name: 'Mbps' },
        series: [
          { name: '下行Mbps', type: 'line', smooth: true, data: rx },
          { name: '上行Mbps', type: 'line', smooth: true, data: tx }
        ]
      }, true)
    })
  }

  const load95 = async () => {
    const res = await getNode95({ nodeId: currentNode.value.ID })
    if (res.code === 0) n95Data.value = res.data || []
  }

  onBeforeUnmount(() => { chart && chart.dispose(); chart = null })
</script>
