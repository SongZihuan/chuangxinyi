import { useSocketStore, useSocketStoreWithOut } from '/@/stores/webstocket';
import VueNativeSock from 'vue-native-websocket-vue3';
import { createApp } from 'vue';
import App from '/@/App.vue';

export const useSocket = () => {
	const app = createApp(App);
	const socketStore = useSocketStore();
	if (socketStore.isConnected) return;
	// localStorage.setItem('socketStatus', 'false');
	if (localStorage.getItem('socketStatus') == 'true' && socketStore.tokenStatus == true) {
		return;
	}

	socketStore.resetState();
	// 安装VueNativeSock
	app.use(VueNativeSock, import.meta.env.VITE_WS_URL, {
		// 启用pinia集成 | enable pinia integration
		store: socketStore,
		// 数据发送/接收使用使用json
		format: 'json',
		// 开启手动调用 connect() 连接服务器
		connectManually: false,
		// 开启自动重连
		reconnection: true,
		// 尝试重连的次数
		reconnectionAttempts: -1,
		// 重连间隔时间
		reconnectionDelay: 3000,
	});
	// 监听页面刷新
	socketStore.sendTokenStatus = false;
	localStorage.setItem('socketStatus', 'true');
};
