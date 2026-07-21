<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">发布版本</el-button>
      </div>
      <el-table :data="list">
        <el-table-column label="版本号" prop="version" width="120" />
        <el-table-column label="下载地址" prop="downloadUrl" min-width="240" show-overflow-tooltip />
        <el-table-column label="校验值" prop="checksum" min-width="200" show-overflow-tooltip />
        <el-table-column label="稳定版" width="80"><template #default="s"><el-tag v-if="s.row.stable" type="success">稳定</el-tag><span v-else>-</span></template></el-table-column>
        <el-table-column label="强制升级" width="90"><template #default="s"><el-tag v-if="s.row.force" type="danger">强制</el-tag><span v-else>-</span></template></el-table-column>
        <el-table-column label="发布时间" width="170"><template #default="s">{{ s.row.CreatedAt }}</template></el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="s">
            <el-button link type="primary" @click="openDialog(s.row)">编辑</el-button>
            <el-button link type="primary" @click="del(s.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination"><el-pagination :current-page="page" :page-size="pageSize" :total="total" @current-change="(v) => { page = v; load() }" /></div>
    </div>

    <el-drawer v-model="visible" size="560" :show-close="false" :before-close="close">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-base">{{ form.ID ? '编辑版本' : '发布版本' }}</span>
          <div>
            <el-button type="primary" @click="enter">确定</el-button>
            <el-button @click="close">取消</el-button>
          </div>
        </div>
      </template>
      <el-form :model="form" label-width="100px">
        <el-form-item label="版本号"><el-input v-model="form.version" placeholder="如 1.0.0" /></el-form-item>
        <el-form-item label="下载地址"><el-input v-model="form.downloadUrl" placeholder="https://.../pcdn-agent" /></el-form-item>
        <el-form-item label="SHA256校验"><el-input v-model="form.checksum" placeholder="可空" /></el-form-item>
        <el-form-item label="稳定版"><el-switch v-model="form.stable" /></el-form-item>
        <el-form-item label="强制升级"><el-switch v-model="form.force" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" /></el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
  import { ref, reactive } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { getReleaseList, createRelease, updateRelease, deleteRelease } from '@/plugin/pcdn/api/release'

  defineOptions({ name: 'PcdnRelease' })

  const list = ref([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(10)
  const load = async () => {
    const r = await getReleaseList({ page: page.value, pageSize: pageSize.value })
    if (r.code === 0) { list.value = r.data.list || []; total.value = r.data.total }
  }
  load()

  const visible = ref(false)
  const form = reactive({ version: '', downloadUrl: '', checksum: '', stable: false, force: false, remark: '', ID: 0 })
  const openDialog = (row) => {
    if (row) { Object.assign(form, row) } else { Object.assign(form, { version: '', downloadUrl: '', checksum: '', stable: false, force: false, remark: '', ID: 0 }) }
    visible.value = true
  }
  const close = () => { visible.value = false }
  const enter = async () => {
    const fn = form.ID ? updateRelease : createRelease
    const r = await fn(form)
    if (r.code === 0) { ElMessage.success('保存成功'); close(); load() }
  }
  const del = (row) => ElMessageBox.confirm('删除该版本?', '提示', { type: 'warning' }).then(async () => { const r = await deleteRelease({ id: row.ID }); if (r.code === 0) { ElMessage.success('已删除'); load() } })
</script>
