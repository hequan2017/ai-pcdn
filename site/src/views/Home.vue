<template>
  <div class="home">
    <header class="nav">
      <div class="logo">PCDN 共享带宽平台</div>
      <nav>
        <router-link to="/login">登录</router-link>
        <router-link to="/register"><el-button type="primary">免费注册</el-button></router-link>
      </nav>
    </header>
    <section class="hero">
      <h1>把闲置带宽变成稳定收益</h1>
      <p class="sub">部署 PCDN 节点，接入抖音 / 腾讯，95 计费或包月透明结算</p>
      <div class="cta">
        <router-link to="/register"><el-button type="primary" size="large">立即上机赚收益</el-button></router-link>
        <router-link to="/login"><el-button size="large">登录控制台</el-button></router-link>
      </div>
    </section>
    <section class="features">
      <div class="feature"><h3>🚀 自助上机</h3><p>一键脚本，分钟级接入</p></div>
      <div class="feature"><h3>📊 实时监控</h3><p>带宽 / 在线率 / 95 值可视</p></div>
      <div class="feature"><h3>💰 透明结算</h3><p>每月账单，明明白白</p></div>
    </section>
  </div>
</template>

<script setup></script>

<style scoped>
  .home { font-family: system-ui, sans-serif; color: #1f2937; }
  .nav { display: flex; justify-content: space-between; align-items: center; padding: 16px 40px; border-bottom: 1px solid #eee; }
  .logo { font-weight: 700; font-size: 18px; }
  .nav nav { display: flex; align-items: center; gap: 16px; }
  .hero { text-align: center; padding: 100px 20px 80px; background: linear-gradient(135deg, #2563eb, #1e40af); color: #fff; }
  .hero h1 { font-size: 42px; margin: 0 0 16px; }
  .hero .sub { font-size: 18px; opacity: 0.9; margin: 0 0 32px; }
  .cta { display: flex; gap: 16px; justify-content: center; }
  .features { display: flex; gap: 24px; justify-content: center; padding: 60px 40px; flex-wrap: wrap; }
  .feature { width: 240px; text-align: center; padding: 24px; border: 1px solid #eee; border-radius: 12px; }
  .feature h3 { margin: 0 0 8px; }
  .feature p { margin: 0; color: #6b7280; }
</style>
