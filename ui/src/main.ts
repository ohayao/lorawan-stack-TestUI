import { createApp } from 'vue';
import App from './App.vue';
//import stores from './vuex';
import vdashboard, { config, IPopOption } from 'vdashboard';
import router from './router/index';
const app = createApp(App);
// stores.forEach((s) => {
//   app.use(s.store, s.key);
// });
config.pop(<IPopOption>{ title: 'Example', css: { width: '250px', height: '180px' }, showMask: true });
app.use(vdashboard);
app.use(router);
app.mount('#app');
