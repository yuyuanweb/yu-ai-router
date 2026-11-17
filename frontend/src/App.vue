<script setup lang="ts">
import BasicLayout from '@/layouts/BasicLayout.vue'
import { theme } from 'ant-design-vue'

// 定义莫兰迪主题 Token
const morandiTheme = {
  token: {
    colorPrimary: '#B8A9B8', // 雾霾紫作为主色
    colorInfo: '#A9B4C2',    // 雾蓝作为信息色
    colorSuccess: '#A8BFA8', // 浅鼠尾草绿 (柔和版)
    colorWarning: '#D4C4B7', // 浅杏色
    colorError: '#C9B0B0',   // 豆沙色作为错误提示
    colorTextBase: '#4A4744', // 文字颜色
    colorBgBase: '#F5F2EF',   // 背景底色
    colorBgContainer: '#FAF7F4', // 卡片/容器底色
    colorBorder: '#EBE5DF',   // 边框颜色
    borderRadius: 12,         // 全局大圆角
    borderRadiusLG: 16,       // 大组件圆角
    fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
    boxShadow: '0 4px 20px rgba(184, 169, 184, 0.08)', // 极淡的柔和阴影
  },
  components: {
    Button: {
      borderRadius: 10,
      controlHeight: 38,
      fontWeight: 500,
    },
    Input: {
      borderRadius: 10,
      controlHeight: 38,
    },
    Card: {
      borderRadiusLG: 16,
    },
    Menu: {
      itemBorderRadius: 8,
      itemSelectedBg: 'rgba(184, 169, 184, 0.15)', // 选中的背景
      itemSelectedColor: '#8B7896', // 选中的文字颜色
    }
  }
}
</script>

<template>
  <a-config-provider :theme="morandiTheme">
    <BasicLayout />
  </a-config-provider>
</template>

<style>
/* 你的原有全局样式 */
* {
  box-sizing: border-box;
}
body {
  margin: 0;
  padding: 0;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  overflow-x: hidden;
}
html {
  overflow-x: hidden;
}
</style>
