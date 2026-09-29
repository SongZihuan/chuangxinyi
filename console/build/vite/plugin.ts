import type { Plugin } from 'vite';
import { viteMockServe } from 'vite-plugin-mock';

export default function mockPlugin(VITE_USE_MOCK: boolean) {
	return {
		...viteMockServe({
			mockPath: 'mock', //mock文件路径，在根路径下创建一个mock文件
			localEnabled: VITE_USE_MOCK, //mock开关
			prodEnabled: false, //生产环境下为false，这样就不会被打包到生产包中
			ignore: /^\_/, //忽略开始_路径
			injectCode: `
			import { setupProdMockServer } from './mockProdServer';
			setupProdMockServer();
		  `,
		}),
		logger: false,
	} as Plugin;
}
