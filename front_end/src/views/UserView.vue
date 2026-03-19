<template>
  <div class="user-page">
    <el-card shadow="hover">
      <template #header>
        <span class="card-title"><el-icon><User /></el-icon> 个人中心</span>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="个人信息" name="info">
          <ProfileForm />
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="password">
          <PasswordForm />
        </el-tab-pane>

        <el-tab-pane label="我的收藏" name="favorites">
          <FavoriteList />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useFishingStore } from '@/stores/fishing'
import { User } from '@element-plus/icons-vue'

// 子组件
import ProfileForm from './user/ProfileForm.vue'
import PasswordForm from './user/PasswordForm.vue'
import FavoriteList from './user/FavoriteList.vue'

const fishingStore = useFishingStore()
const activeTab = ref('info')

onMounted(() => {
  fishingStore.fetchFavorites()
})
</script>

<style scoped>
.user-page {
  max-width: 800px;
  margin: 0 auto;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  font-size: 16px;
}
</style>
