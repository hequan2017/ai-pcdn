<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="search">
        <el-form-item label="状态">
          <el-select v-model="search.status" clearable style="width: 120px">
            <el-option label="触发中" value="firing" />
            <el-option label="已恢复" value="resolved" />
          </el-select>
        </el-form-item>
        <el-form-item label="节点ID"><el-input v-model="search.nodeId" clearable style="width: 120px" /></el-form-item>
        <el-form-item><el-button type="primary" icon="search" @click="load">查询</el-button></el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="list">
        <el-table-column label="规则" prop="ruleName" min-width="140" />
        <el-table-column label="节点SN" prop="nodeSn" min-width="180" />
        <el-table-column label="指标" width="120">
          <template #default="s">{{ ({ offline: '离线', bandwidth_low: '带宽低', p95_high: '95值高', agent_down: '上报中断' }[s.row.metric]) || s.row.metric }}</template>
        </el-table-column>
        <el-table-column label="触发值" width="140">
          <template #default="s">{{ (s.row.metric === 'bandwidth_low' || s.row.metric === 'p95_high') ? (s.row.triggerValue / 1e6).toFixed(2) + ' Mbps' : s.row.triggerValue }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="s">
            <el-tag :type="s.row.status === 'firing' ? 'danger' : 'success'">{{ s.row.status === 'firing' ? '触发中' : '已恢复' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="触发时间" prop="firedAt" width="170">
          <template #default="s">{{ formatDate(s.row.firedAt) }}</template>
        </el-table-column>
        <el-table-column label="恢复时间" prop="resolvedAt" width="170">
          <template #default="s">{{ s.row.resolvedAt ? formatDate(s.row.resolvedAt) : '-' }}</template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination :current-page="page" :page-size="pageSize" :total="total" @current-change="(v) => { page = v; load() }" />
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, reactive } from 'vue'
  import { getAlarmRecordList } from '@/plugin/pcdn/api/alarm'
  import { formatDate } from '@/utils/format'

  defineOptions({ name: 'PcdnAlarmRecord' })

  const list = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const search = reactive({ status: '', nodeId: '' })

  const load = async () => {
    const r = await getAlarmRecordList({ page: page.value, pageSize: pageSize.value, ...search })
    if (r.code === 0) { list.value = r.data.list || []; total.value = r.data.total }
  }
  load()
</script>
