<template>
  <div class="banner">
    <!-- 背景装饰 -->
    <div class="banner-bg">
      <!-- 浮动气泡 -->
      <div class="bubble bubble-1"></div>
      <div class="bubble bubble-2"></div>
      <div class="bubble bubble-3"></div>
      <div class="bubble bubble-4"></div>
      <div class="bubble bubble-5"></div>

      <!-- SVG 鱼群装饰 -->
      <svg class="fish fish-left" viewBox="0 0 64 32" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M48 16C48 16 36 4 16 4C10 4 4 8 2 16C4 24 10 28 16 28C36 28 48 16 48 16Z" fill="currentColor"/>
        <path d="M56 16L62 8V24L56 16Z" fill="currentColor"/>
        <circle cx="12" cy="14" r="2" fill="white" opacity="0.8"/>
      </svg>
      <svg class="fish fish-right" viewBox="0 0 64 32" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M16 16C16 16 28 4 48 4C54 4 60 8 62 16C60 24 54 28 48 28C28 28 16 16 16 16Z" fill="currentColor"/>
        <path d="M8 16L2 8V24L8 16Z" fill="currentColor"/>
        <circle cx="52" cy="14" r="2" fill="white" opacity="0.8"/>
      </svg>
      <svg class="fish fish-small" viewBox="0 0 64 32" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M48 16C48 16 36 4 16 4C10 4 4 8 2 16C4 24 10 28 16 28C36 28 48 16 48 16Z" fill="currentColor"/>
        <path d="M56 16L62 8V24L56 16Z" fill="currentColor"/>
        <circle cx="12" cy="14" r="2" fill="white" opacity="0.8"/>
      </svg>
    </div>

    <!-- 水波纹底部 -->
    <div class="wave-container">
      <svg class="wave wave-back" viewBox="0 0 2880 60" preserveAspectRatio="none">
        <path d="M0,30 C240,60 480,0 720,30 C960,60 1200,0 1440,30 C1680,60 1920,0 2160,30 C2400,60 2640,0 2880,30 L2880,60 L0,60 Z" fill="rgba(255,255,255,0.15)"/>
      </svg>
      <svg class="wave wave-front" viewBox="0 0 2880 60" preserveAspectRatio="none">
        <path d="M0,30 C360,0 720,60 1080,30 C1260,15 1440,45 1620,30 C1800,15 2160,60 2520,30 C2700,15 2790,45 2880,30 L2880,60 L0,60 Z" fill="rgba(255,255,255,0.25)"/>
      </svg>
    </div>

    <!-- 内容区 -->
    <div class="banner-content">
      <div class="brand-row">
        <div class="brand-icon">
          <svg viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 6C14 6 8 10 6 16C8 22 14 26 20 26C26 26 32 22 34 16C32 10 26 6 20 6Z" stroke="white" stroke-width="2" fill="none"/>
            <path d="M36 16L40 10V22L36 16Z" fill="white" opacity="0.8"/>
            <circle cx="14" cy="14" r="1.5" fill="white"/>
            <path d="M8 32L14 28H26L32 32" stroke="white" stroke-width="1.5" stroke-linecap="round" opacity="0.6"/>
            <path d="M12 36L16 33H24L28 36" stroke="white" stroke-width="1" stroke-linecap="round" opacity="0.4"/>
          </svg>
        </div>
        <div class="brand-text">
          <h1 class="title">智钓蓝海<span class="highlight">智能垂钓</span>系统</h1>
          <p class="subtitle">实时监测各垂钓水域环境与钓况数据 · 智能推荐最佳垂钓方案</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.banner {
  position: relative;
  background: linear-gradient(135deg, #1a6fc4 0%, #2196f3 30%, #29b6f6 60%, #4fc3f7 100%);
  border-radius: 14px;
  overflow: hidden;
  padding: 36px 40px 28px;
  min-height: 130px;
}

/* ===== 背景装饰 ===== */
.banner-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

/* 浮动气泡 */
.bubble {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  animation: float 8s ease-in-out infinite;
}

.bubble-1 { width: 120px; height: 120px; top: -30px; right: 10%; animation-delay: 0s; }
.bubble-2 { width: 80px; height: 80px; top: 50%; right: 25%; animation-delay: -2s; opacity: 0.6; }
.bubble-3 { width: 60px; height: 60px; bottom: 10px; left: 8%; animation-delay: -4s; }
.bubble-4 { width: 200px; height: 200px; top: -80px; left: -40px; animation-delay: -1s; background: rgba(255,255,255,0.04); }
.bubble-5 { width: 40px; height: 40px; top: 20%; right: 40%; animation-delay: -3s; background: rgba(255,255,255,0.1); }

@keyframes float {
  0%, 100% { transform: translateY(0) scale(1); }
  50% { transform: translateY(-12px) scale(1.05); }
}

/* SVG 鱼装饰 */
.fish {
  position: absolute;
  color: rgba(255, 255, 255, 0.1);
}

.fish-left {
  width: 70px;
  top: 18%;
  left: 5%;
  animation: swim-left 12s ease-in-out infinite;
}

.fish-right {
  width: 50px;
  bottom: 25%;
  right: 8%;
  animation: swim-right 15s ease-in-out infinite;
  color: rgba(255, 255, 255, 0.08);
}

.fish-small {
  width: 35px;
  top: 60%;
  left: 30%;
  animation: swim-left 18s ease-in-out infinite;
  animation-delay: -5s;
  color: rgba(255, 255, 255, 0.06);
}

@keyframes swim-left {
  0%, 100% { transform: translateX(0) translateY(0); }
  25% { transform: translateX(20px) translateY(-5px); }
  50% { transform: translateX(40px) translateY(3px); }
  75% { transform: translateX(15px) translateY(-3px); }
}

@keyframes swim-right {
  0%, 100% { transform: translateX(0) translateY(0); }
  25% { transform: translateX(-15px) translateY(5px); }
  50% { transform: translateX(-30px) translateY(-3px); }
  75% { transform: translateX(-10px) translateY(4px); }
}

/* 水波 */
.wave-container {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 30px;
  pointer-events: none;
  overflow: hidden;
}

.wave {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 200%;
  height: 100%;
}

.wave-back {
  animation: wave-move 6s ease-in-out infinite;
}

.wave-front {
  animation: wave-move 4s ease-in-out infinite reverse;
}

@keyframes wave-move {
  0%, 100% { transform: translateX(0); }
  50% { transform: translateX(-25%); }
}

/* ===== 内容区 ===== */
.banner-content {
  position: relative;
  z-index: 2;
}

.brand-row {
  display: flex;
  align-items: center;
  gap: 20px;
}

.brand-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 14px;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  flex-shrink: 0;
}

.brand-icon svg {
  width: 34px;
  height: 34px;
}

.brand-text {
  flex: 1;
}

.title {
  font-size: 1.7rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 6px;
  letter-spacing: 2px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  line-height: 1.3;
}

.title .highlight {
  position: relative;
  padding: 0 2px;
}

.title .highlight::after {
  content: "";
  position: absolute;
  bottom: 2px;
  left: 0;
  width: 100%;
  height: 6px;
  background: rgba(255, 255, 255, 0.25);
  border-radius: 3px;
  z-index: -1;
}

.subtitle {
  margin: 0;
  font-size: 0.95rem;
  color: rgba(255, 255, 255, 0.85);
  letter-spacing: 1px;
  line-height: 1.5;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .banner {
    padding: 24px 20px 20px;
    min-height: 100px;
  }

  .brand-icon {
    width: 44px;
    height: 44px;
    border-radius: 10px;
  }

  .brand-icon svg {
    width: 26px;
    height: 26px;
  }

  .title {
    font-size: 1.3rem;
    letter-spacing: 1px;
  }

  .subtitle {
    font-size: 0.82rem;
  }

  .fish-left, .fish-right, .fish-small {
    display: none;
  }

  .bubble-1 { width: 60px; height: 60px; }
  .bubble-4 { width: 100px; height: 100px; }
}

@media (max-width: 480px) {
  .banner {
    padding: 18px 16px 16px;
  }

  .brand-row {
    gap: 12px;
  }

  .title {
    font-size: 1.1rem;
  }

  .subtitle {
    font-size: 0.75rem;
  }
}
</style>
