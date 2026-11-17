
<script setup lang="ts">
// 无需额外的响应式数据
</script>

/* 替换原有的 style 部分 */
<style scoped>
.header {
  background: var(--bg-card) !important; /* 使用奶白底色 */
  padding: 0 24px;
  border-bottom: 1px solid var(--border-color);
  box-shadow: none; /* 移除硬阴影 */
  height: 64px;
  display: flex;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  height: 40px;
  width: 40px;
  border-radius: 50%; /* logo 也可以圆润一点 */
}

.site-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-main); /* 深炭灰替代原来的蓝色 */
  letter-spacing: 1px;
}

/* 去除菜单默认的边框 */
:deep(.ant-menu-horizontal) {
  border-bottom: none !important;
  background: transparent;
}

/* 菜单项 hover 效果放慢、变柔和 */
:deep(.ant-menu-item) {
  transition: all 0.3s ease-in-out !important;
  border-radius: 8px !important;
}

:deep(.ant-menu-item:hover) {
  background-color: rgba(184, 169, 184, 0.1) !important;
  color: var(--color-primary) !important;
}
</style>
