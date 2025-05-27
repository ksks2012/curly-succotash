<template>
  <div class="app-container">
    <!-- Navigation Bar -->
    <nav class="bg-blue-600 p-4 shadow-md">
      <div class="container mx-auto flex justify-between items-center">
        <h1 class="text-white text-2xl font-bold">{{ $t('appTitle') }}</h1>
        <div class="flex items-center space-x-4">
          <router-link to="/" class="text-white hover:text-blue-200">{{ $t('home') }}</router-link>
          <router-link to="/about" class="text-white hover:text-blue-200">{{ $t('about') }}</router-link>
          <!-- Language Switcher -->
          <select v-model="$i18n.locale" class="bg-blue-500 text-white p-2 rounded">
            <option value="en">English</option>
            <option value="zh">中文</option>
          </select>
        </div>
      </div>
    </nav>
    <!-- Main Content -->
    <div class="container mx-auto p-4">
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
};
</script>

<style>
.app-container {
  min-height: 100vh;
  background-color: #f7fafc;
}
</style>