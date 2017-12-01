// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import VueApollo from 'vue-apollo'

import App from './App'
import {apolloClient} from './apollo.coffee'

Vue.use(VueApollo)

Vue.config.productionTip = false

const apolloProvider = new VueApollo({
  defaultClient:apolloClient,
})

/* eslint-disable no-new */
new Vue({
  apolloProvider,
  el: '#app',
  template: '<App/>',
  components: { App }
})
