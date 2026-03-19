<template>
  <el-footer class="footer">
    <div class="footer-content">
      <span>智钓蓝海信息平台 &copy; {{ year }}</span>
      <span class="separator">|</span>
      <span>智钓盈联——边缘-云协同的全生态智能垂钓系统</span>
    </div>
  </el-footer>
</template>

<script setup lang="ts">
const year = new Date().getFullYear()
</script>

<style scoped>
.footer {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  border-top: 1px solid #ebeef5;
}

.footer-content {
  font-size: 13px;
  color: #909399;
}

.separator {
  margin: 0 8px;
}
</style>
