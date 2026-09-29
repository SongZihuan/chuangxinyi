import { createApp } from 'vue';
import pinia from '/@/stores/index';
import App from '/@/App.vue';
import router from '/@/router';
import { directive } from '/@/directive';
import other from '/@/utils/other';
import ElementPlus from 'element-plus';
import './input.css';
import '/@/theme/index.scss';
import Vue3VideoPlayer from '@cloudgeek/vue3-video-player';
// import vConsole from 'vconsole';
// const vconsole = new vConsole();
import AudioPlayer from '@liripeng/vue-audio-player';
import 'virtual:svg-icons-register';

const app = createApp(App);

directive(app);
other.elSvg(app);

app.use(pinia).use(router).use(ElementPlus).use(Vue3VideoPlayer).use(AudioPlayer).mount('#app');
