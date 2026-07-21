<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="search">
        <el-form-item label="账期"><el-input v-model="search.period" placeholder="2026-07" style="width: 120px" /></el-form-item>
        <el-form-item label="状态">
          <el-select v-model="search.status" clearable style="width: 120px">
            <el-option label="待审核" value="draft" /><el-option label="已审核" value="approved" /><el-option label="已付款" value="paid" /><el-option label="已驳回" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item label="归属"><el-input v-model="search.ownerName" /></el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="load">查询</el-button>
          <el-button type="success" @click="showGen = true">生成账单</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="list">
        <el-table-column label="账期" prop="period" width="100" />
        <el-table-column label="归属" prop="ownerName" width="120" />
        <el-table-column label="节点数" prop="nodeCount" width="80" />
        <el-table-column label="应付总额" width="120"><template #default="s">¥{{ (s.row.totalAmount || 0).toFixed(2) }}</template></el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="s"><el-tag :type="statusType(s.row.status)">{{ statusLabel(s.row.status) }}</el-tag></template>
        </el-table-column>
        <el-table-column label="操作" width="300">
          <template #default="s">
            <el-button link type="primary" @click="openDetail(s.row)">明细</el-button>
            <el-button v-if="s.row.status === 'draft'" link type="success" @click="approve(s.row)">审核</el-button>
            <el-button v-if="s.row.status === 'draft'" link type="warning" @click="reject(s.row)">驳回</el-button>
            <el-button v-if="s.row.status === 'approved'" link type="primary" @click="openPay(s.row)">付款</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination"><el-pagination :current-page="page" :page-size="pageSize" :total="total" @current-change="(v) => { page = v; load() }" /></div>
    </div>

    <el-dialog v-model="showGen" title="生成账单" width="400">
      <el-form label-width="80px"><el-form-item label="账期"><el-input v-model="genPeriod" placeholder="2026-07" /></el-form-item></el-form>
      <template #footer><el-button type="primary" :loading="gening" @click="doGen">生成</el-button><el-button @click="showGen = false">取消</el-button></template>
    </el-dialog>

    <el-dialog v-model="showPay" title="付款" width="460">
      <el-form :model="payForm" label-width="100px">
        <el-form-item label="应付总额">¥{{ (payForm.totalAmount || 0).toFixed(2) }}</el-form-item>
        <el-form-item label="实付金额"><el-input-number v-model="payForm.paidAmount" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="付款方式"><el-input v-model="payForm.payMethod" placeholder="银行/支付宝/微信" /></el-form-item>
        <el-form-item label="流水号"><el-input v-model="payForm.payNo" /></el-form-item>
      </el-form>
      <template #footer><el-button type="primary" :loading="paying" @click="doPay">确认付款</el-button><el-button @click="showPay = false">取消</el-button></template>
    </el-dialog>

    <el-drawer v-model="showDetail" size="640" :title="`账单明细 · ${current?.period || ''}`">
      <el-table :data="currentDetails" size="small">
        <el-table-column label="节点SN" prop="nodeSn" min-width="180" />
        <el-table-column label="计费" width="80"><template #default="s">{{ s.row.billingMode === 'monthly' ? '包月' : '95' }}</template></el-table-column>
        <el-table-column label="数值" width="120"><template #default="s">{{ s.row.billingMode === 'p95' ? (s.row.value || 0).toFixed(2) + ' Mbps' : '-' }}</template></el-table-column>
        <el-table-column label="单价" width="120"><template #default="s">{{ s.row.billingMode === 'p95' ? '¥' + s.row.unitPrice + '/Mbps' : '-' }}</template></el-table-column>
        <el-table-column label="金额" width="100"><template #default="s">¥{{ (s.row.amount || 0).toFixed(2) }}</template></el-table-column>
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup>
  import { ref, reactive } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { generateBill, getBillList, approveBill, rejectBill, payBill } from '@/plugin/pcdn/api/bill'

  defineOptions({ name: 'PcdnBill' })

  const list = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const search = reactive({ period: '', status: '', ownerName: '' })
  const load = async () => {
    const r = await getBillList({ page: page.value, pageSize: pageSize.value, ...search })
    if (r.code === 0) { list.value = r.data.list || []; total.value = r.data.total }
  }
  load()

  const statusType = (s) => ({ draft: 'warning', approved: 'primary', paid: 'success', rejected: 'info' }[s] || 'info')
  const statusLabel = (s) => ({ draft: '待审核', approved: '已审核', paid: '已付款', rejected: '已驳回' }[s] || s)

  const showGen = ref(false)
  const genPeriod = ref('')
  const gening = ref(false)
  const doGen = async () => {
    if (!genPeriod.value) { ElMessage.warning('请填账期'); return }
    gening.value = true
    try { const r = await generateBill({ period: genPeriod.value }); if (r.code === 0) { ElMessage.success('生成 ' + r.data.created + ' 条'); showGen.value = false; load() } } finally { gening.value = false }
  }

  const approve = (row) => ElMessageBox.confirm('审核通过该账单?', '提示').then(async () => { const r = await approveBill({ id: row.ID }); if (r.code === 0) { ElMessage.success('已审核'); load() } })
  const reject = (row) => ElMessageBox.prompt('驳回原因', '驳回').then(async ({ value }) => { const r = await rejectBill({ id: row.ID, remark: value }); if (r.code === 0) { ElMessage.success('已驳回'); load() } }).catch(() => {})

  const showPay = ref(false)
  const paying = ref(false)
  const payForm = reactive({ id: 0, paidAmount: 0, payMethod: '', payNo: '', totalAmount: 0 })
  const openPay = (row) => { payForm.id = row.ID; payForm.paidAmount = row.totalAmount; payForm.payMethod = ''; payForm.payNo = ''; payForm.totalAmount = row.totalAmount; showPay.value = true }
  const doPay = async () => {
    paying.value = true
    try { const r = await payBill({ id: payForm.id, paidAmount: payForm.paidAmount, payMethod: payForm.payMethod, payNo: payForm.payNo }); if (r.code === 0) { ElMessage.success('已付款'); showPay.value = false; load() } } finally { paying.value = false }
  }

  const showDetail = ref(false)
  const current = ref(null)
  const currentDetails = ref([])
  const openDetail = (row) => { current.value = row; currentDetails.value = row.details || []; showDetail.value = true }
</script>
