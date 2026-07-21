<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true">
        <el-form-item label="账期"><el-input v-model="period" style="width: 120px" /></el-form-item>
        <el-form-item><el-button type="primary" @click="loadAll">查询</el-button></el-form-item>
      </el-form>
    </div>
    <div class="cards">
      <div class="card"><div class="label">大厂收入</div><div class="value">¥{{ (summary.revenue || 0).toFixed(2) }}</div></div>
      <div class="card"><div class="label">采购成本</div><div class="value">¥{{ (summary.cost || 0).toFixed(2) }}</div></div>
      <div class="card"><div class="label">利润</div><div class="value" :class="{ pos: summary.profit >= 0, neg: summary.profit < 0 }">¥{{ (summary.profit || 0).toFixed(2) }}</div></div>
      <div class="card"><div class="label">利润率</div><div class="value">{{ ((summary.profitMargin || 0) * 100).toFixed(1) }}%</div></div>
    </div>
    <div class="row">
      <div class="col">
        <div class="title">平台收入</div>
        <el-table :data="platforms" size="small">
          <el-table-column label="平台" prop="platform" /><el-table-column label="收入"><template #default="s">¥{{ (s.row.revenue || 0).toFixed(2) }}</template></el-table-column><el-table-column label="单数" prop="count" width="80" />
        </el-table>
      </div>
      <div class="col">
        <div class="title">贡献者成本 Top</div>
        <el-table :data="owners" size="small">
          <el-table-column label="贡献者" prop="ownerName" /><el-table-column label="成本"><template #default="s">¥{{ (s.row.cost || 0).toFixed(2) }}</template></el-table-column><el-table-column label="单数" prop="count" width="80" />
        </el-table>
      </div>
    </div>
    <div class="title">利润趋势（近6月）</div>
    <el-table :data="trend" size="small">
      <el-table-column label="账期" prop="period" width="100" /><el-table-column label="收入"><template #default="s">¥{{ (s.row.revenue || 0).toFixed(2) }}</template></el-table-column><el-table-column label="成本"><template #default="s">¥{{ (s.row.cost || 0).toFixed(2) }}</template></el-table-column><el-table-column label="利润"><template #default="s">¥{{ (s.row.profit || 0).toFixed(2) }}</template></el-table-column>
    </el-table>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { getProfitSummary, getRevenueByPlatform, getCostByOwner, getProfitTrend } from '@/plugin/pcdn/api/profit'

  defineOptions({ name: 'PcdnProfit' })

  const now = new Date()
  const period = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)
  const summary = ref({})
  const platforms = ref([])
  const owners = ref([])
  const trend = ref([])

  const loadAll = async () => {
    const [s, p, o, t] = await Promise.all([
      getProfitSummary({ period: period.value }),
      getRevenueByPlatform({ period: period.value }),
      getCostByOwner({ period: period.value }),
      getProfitTrend({ months: 6 })
    ])
    if (s.code === 0) summary.value = s.data
    if (p.code === 0) platforms.value = p.data || []
    if (o.code === 0) owners.value = o.data || []
    if (t.code === 0) trend.value = t.data || []
  }
  onMounted(loadAll)
</script>

<style scoped>
  .cards { display: flex; gap: 16px; margin-bottom: 20px; }
  .card { flex: 1; padding: 20px; border: 1px solid #eee; border-radius: 8px; text-align: center; }
  .card .label { color: #6b7280; font-size: 13px; }
  .card .value { font-size: 26px; font-weight: 700; margin-top: 6px; }
  .card .value.pos { color: #67c23a; }
  .card .value.neg { color: #f56c6c; }
  .row { display: flex; gap: 16px; margin-bottom: 20px; }
  .col { flex: 1; }
  .title { font-weight: 700; margin-bottom: 8px; }
</style>
