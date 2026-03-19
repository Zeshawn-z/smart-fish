<template>
  <div class="not-found">
    <el-empty description="页面不存在" :image-size="200">
      <el-button type="primary" @click="router.push('/')">返回首页</el-button>
    </el-empty>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<style scoped>
.not-found {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 200px);
}
</style>
