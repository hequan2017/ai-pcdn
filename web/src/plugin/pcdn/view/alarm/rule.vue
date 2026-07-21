<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增规则</el-button>
        <el-button icon="delete" :disabled="!sel.length" @click="onDelete">删除</el-button>
      </div>
      <el-table :data="list" row-key="ID" @selection-change="(v) => (sel = v)">
        <el-table-column type="selection" width="55" />
        <el-table-column label="规则名" prop="name" min-width="140" />
        <el-table-column label="范围" width="100">
          <template #default="s">{{ scopeLabel(s.row.scopeType) }}</template>
        </el-table-column>
        <el-table-column label="指标" width="130">
          <template #default="s">{{ metricLabel(s.row.metric) }}</template>
        </el-table-column>
        <el-table-column label="阈值(Mbps)" width="120">
          <template #default="s">{{ s.row.threshold ? (s.row.threshold / 1e6).toFixed(2) : '-' }}</template>
        </el-table-column>
        <el-table-column label="启用" width="80">
          <template #default="s">
            <el-switch v-model="s.row.enabled" @change="onToggle(s.row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="s">
            <el-button link type="primary" @click="openDialog(s.row)">编辑</el-button>
            <el-button link type="primary" @click="deleteRow(s.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination :current-page="page" :page-size="pageSize" :total="total" @current-change="(v) => { page = v; load() }" />
      </div>
    </div>

    <el-drawer v-model="visible" size="560" :show-close="false" :before-close="close">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-base">{{ form.ID ? '编辑规则' : '新增规则' }}</span>
          <div>
            <el-button type="primary" @click="enter">确定</el-button>
            <el-button @click="close">取消</el-button>
          </div>
        </div>
      </template>
      <el-form :model="form" label-width="100px">
        <el-form-item label="规则名"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="范围">
          <el-select v-model="form.scopeType">
            <el-option label="全部节点" value="all" />
            <el-option label="按分组" value="group" />
            <el-option label="单节点" value="node" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.scopeType !== 'all'" label="范围值">
          <el-input v-model="form.scopeValue" :placeholder="form.scopeType === 'group' ? '分组ID' : '节点ID'" />
        </el-form-item>
        <el-form-item label="指标">
          <el-select v-model="form.metric">
            <el-option label="节点离线" value="offline" />
            <el-option label="带宽低于阈值" value="bandwidth_low" />
            <el-option label="95值高于阈值" value="p95_high" />
            <el-option label="agent上报中断" value="agent_down" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.metric === 'bandwidth_low' || form.metric === 'p95_high'" label="阈值(Mbps)">
          <el-input-number v-model="thresholdMbps" :min="0" />
        </el-form-item>
        <el-form-item v-if="form.metric === 'agent_down'" label="持续秒数">
          <el-input-number v-model="form.durationSec" :min="0" />
        </el-form-item>
        <el-form-item label="Webhook"><el-input v-model="webhookUrl" placeholder="钉钉/企微 webhook 地址" /></el-form-item>
        <el-form-item label="@手机号"><el-input v-model="atMobilesStr" placeholder="多个用逗号分隔" /></el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import { ref, reactive } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { getAlarmRuleList, createAlarmRule, updateAlarmRule, deleteAlarmRule, deleteAlarmRuleByIds } from '@/plugin/pcdn/api/alarm'

  defineOptions({ name: 'PcdnAlarmRule' })

  const list = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const sel = ref([])

  const load = async () => {
    const r = await getAlarmRuleList({ page: page.value, pageSize: pageSize.value })
    if (r.code === 0) { list.value = r.data.list || []; total.value = r.data.total }
  }
  load()

  const scopeLabel = (s) => ({ all: '全部', group: '分组', node: '单节点' }[s] || s)
  const metricLabel = (m) => ({ offline: '节点离线', bandwidth_low: '带宽低', p95_high: '95值高', agent_down: '上报中断' }[m] || m)

  const visible = ref(false)
  const form = reactive({ name: '', scopeType: 'all', scopeValue: '', metric: 'offline', threshold: 0, durationSec: 300, enabled: true, notifyConfig: null, ID: 0 })
  const webhookUrl = ref('')
  const atMobilesStr = ref('')
  const thresholdMbps = ref(0)

  const parseNotify = (v) => {
    if (!v) return {}
    if (typeof v === 'string') { try { return JSON.parse(v) } catch (e) { return {} } }
    return v
  }

  const openDialog = (row) => {
    if (row) {
      Object.assign(form, row)
      const c = parseNotify(row.notifyConfig)
      webhookUrl.value = c.webhookUrl || ''
      atMobilesStr.value = (c.atMobiles || []).join(',')
      thresholdMbps.value = row.threshold ? row.threshold / 1e6 : 0
    } else {
      Object.assign(form, { name: '', scopeType: 'all', scopeValue: '', metric: 'offline', threshold: 0, durationSec: 300, enabled: true, notifyConfig: null, ID: 0 })
      webhookUrl.value = ''
      atMobilesStr.value = ''
      thresholdMbps.value = 0
    }
    visible.value = true
  }
  const close = () => { visible.value = false }

  const enter = async () => {
    form.threshold = Math.round(thresholdMbps.value * 1e6)
    form.notifyConfig = { webhookUrl: webhookUrl.value, atMobiles: atMobilesStr.value ? atMobilesStr.value.split(',').map((s) => s.trim()) : [] }
    const fn = form.ID ? updateAlarmRule : createAlarmRule
    const r = await fn(form)
    if (r.code === 0) { ElMessage.success('保存成功'); close(); load() }
  }

  const onToggle = async (row) => { await updateAlarmRule(row); ElMessage.success(row.enabled ? '已启用' : '已停用') }
  const deleteRow = (row) => ElMessageBox.confirm('确定删除该规则?', '提示', { type: 'warning' }).then(async () => { const r = await deleteAlarmRule({ id: row.ID }); if (r.code === 0) { ElMessage.success('已删除'); load() } })
  const onDelete = () => { if (!sel.value.length) return; ElMessageBox.confirm('批量删除选中规则?', '提示', { type: 'warning' }).then(async () => { const r = await deleteAlarmRuleByIds({ ids: sel.value.map((i) => i.ID) }); if (r.code === 0) { ElMessage.success('已删除'); load() } }) }
</script>
