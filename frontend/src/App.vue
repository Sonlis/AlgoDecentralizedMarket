<template>
  <div id="app">
    <Charts />
  </div>
</template>

<script>
import Charts from './components/Charts.vue'

export default {
  name: 'App',
  components: {
    Charts
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
