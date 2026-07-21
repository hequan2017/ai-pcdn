<!-- Modified for ai-pcdn on 2026-07-21: 重构关于页为 ai-pcdn 产品信息。 -->
<template>
  <div class="mt-2">
    <el-card>
      <section class="mx-auto max-w-4xl px-4 py-10 sm:px-8 sm:py-14">
        <div
          class="mb-5 inline-flex rounded-full bg-blue-50 px-3 py-1 text-xs font-semibold tracking-[0.16em] text-blue-600 dark:bg-blue-950/40 dark:text-blue-300"
        >
          AI-PCDN
        </div>
        <h1 class="text-3xl font-bold text-base-text sm:text-4xl">ai-pcdn</h1>
        <p class="mt-4 max-w-2xl text-base leading-7 text-muted-foreground">
          面向 PCDN 业务的管理平台，提供统一、清晰的后台操作入口。
        </p>

        <dl class="mt-10 grid gap-4 sm:grid-cols-2">
          <div class="rounded-xl border border-border bg-main p-5">
            <dt class="text-sm text-muted-foreground">产品名称</dt>
            <dd class="mt-2 text-lg font-semibold text-base-text">ai-pcdn</dd>
          </div>
          <div class="rounded-xl border border-border bg-main p-5">
            <dt class="text-sm text-muted-foreground">产品定位</dt>
            <dd class="mt-2 text-lg font-semibold text-base-text">PCDN 管理平台</dd>
          </div>
          <div class="rounded-xl border border-border bg-main p-5">
            <dt class="text-sm text-muted-foreground">前端技术</dt>
            <dd class="mt-2 text-lg font-semibold text-base-text">Vue 3 / Vite</dd>
          </div>
          <div class="rounded-xl border border-border bg-main p-5">
            <dt class="text-sm text-muted-foreground">后端与数据</dt>
            <dd class="mt-2 text-lg font-semibold text-base-text">Go / Gin / SQLite</dd>
          </div>
        </dl>
      </section>
    </el-card>
  </div>
</template>

<script setup>
  defineOptions({
    name: 'About'
  })
</script>
