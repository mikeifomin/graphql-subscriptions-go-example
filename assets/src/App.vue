<template lang="pug">
#app
  Plot
  
</template>

<script lang="coffee">
import Plot from './components/Plot'

export default {
  name: 'app'
  components: {
    Plot
  }
}
</script>

<style lang="stylus">
#app 
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
body, html
  padding 0
  margin 0

</style>
