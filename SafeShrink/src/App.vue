<template>
  <div id="app">
    <SafeShrink />
  </div>
</template>

<script setup>
import SafeShrink from './components/SafeShrink.vue'
</script>

<style>
#app {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  background-color: #111827;
  color: #ffffff;
}
</style>