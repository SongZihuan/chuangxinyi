// https://pinia.vuejs.org/
import { createPinia } from 'pinia';
import { App } from "vue";
// 创建
const pinia = createPinia();

// 导出
export function setupStore(app: App<Element>) {
    app.use(pinia);
}
export default pinia;
