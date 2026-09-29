import axios, { CreateAxiosDefaults } from 'axios';
import { ElLoading, ElMessage } from 'element-plus';
import { Local, Session } from '/@/utils/storage';
import qs from 'qs';
import { useLoginSignSub, useLoginSignIn } from '/@/hooks/useLoginSignIn';
import { useRoute } from 'vue-router';
import router from '../router/index';

// 配置新建一个 axios 实例
const _service = axios.create({
	withCredentials: false,
	baseURL: import.meta.env.VITE_API_URL,
	timeout: 50000,
	headers: { 'Content-Type': 'application/json' },
	paramsSerializer: {
		serialize(params) {
			return qs.stringify(params, { allowDots: true });
		},
	},
	isLoginToken: false,
	notJson: false,
	isRetry: false,
	logout: null,

	routePath: '',
	routeQuery: {},
} as CreateAxiosDefaults<any>);

// 添加请求拦截器
_service.interceptors.request.use(
	(config: any) => {
		config.headers!['X-RunMode'] = 'release';

		if (!config.logout) {
			config.logout = () => {
				let nextRoute = {
					path: '/login',
					query: {
						forceLogin: 1,
					},
				};

				if (config.routePath) {
					nextRoute = {
						path: '/login',
						query: {
							redirect: config.routePath,
							params: JSON.stringify(config.routeQuery),
							forceLogin: 1,
						},
					};
				}

				const loadingInstance = ElLoading.service({ fullscreen: true, text: '你的登录过期了' });
				setTimeout(async () => {
					Session.clear(); // 清除浏览器全部临时缓存
					Local.clear();
					await router.push(nextRoute);
					loadingInstance.close();
				}, 1000);
			};
		}

		let token = Session.get('token');
		let logintoken = Session.get('login-token')?.token;
		if (token && logintoken) {
			if (config.isLoginToken as boolean) {
				config.headers!['x-token'] = `${logintoken}`;
			} else {
				config.headers!['x-token'] = `${token}`;
			}
		}

		return config;
	},
	(error) => {
		// 对请求错误做些什么
		return Promise.reject(error);
	}
);

const retry = 'FRONT_RETRY';

interface resp {
	code: string;
	subCode: string;
	msg: string;
	data: any;
}

const NotMessageCode = ['USER_NOT_HOME_PAGE', 'NOT_OPEN_WEBSITE', 'FACE_CHECK_WAIT', 'BAD_PASS_TOKEN', 'DEFRAY_INSUFFICIENT', 'PAY_MUST_VERIFY'];

// 添加响应拦截器
_service.interceptors.response.use(
	async (response: any): Promise<any> => {
		if (response.config.notJson) {
			return response.data;
		}

		const res = response.data as resp;

		if (!res.code || !res.subCode) {
			return Promise.resolve(_service.interceptors.response);
		}

		if (res.code === 'SUCCESS') {
			return Promise.resolve(res);
		} else if (res.code === 'POLICY_DENY' || res.code === 'CORS_DENY' || res.code === 'WEBSITE_DENY' || res.code === 'NOT_TOKEN_DENY') {
			return Promise.resolve(res);
		} else if (res.code === 'DOUBLE_CHECK_DENY') {
			ElMessage.error({ message: '二次身份验证失败，请重试', grouping: true });
			return Promise.resolve(res);
		} else if (res.code === 'ROBOT_DENY') {
			if (res.subCode !== 'CAPTCHA_SECOND_CHECK') {
				ElMessage.error({ message: '人机验证失败，行为疑似机器人', grouping: true });
			}
			return Promise.resolve(res);
		} else if (res.code === 'TOKEN_DENY') {
			if (response.config.isRetry) {
				if (response.config.logout) {
					await response.config.logout(response);
				} else {
					const loadingInstance = ElLoading.service({ fullscreen: true, text: '你的登录过期了' });
					setTimeout(() => {
						Session.clear(); // 清除浏览器全部临时缓存
						Local.clear();
						window.location.href = '/login';
						loadingInstance.close();
					}, 1000);
				}
				return Promise.resolve(_service.interceptors.response);
			}

			const { regetToken } = useLoginSignSub();
			const { isLoginUser, resetSignIn } = useLoginSignIn();

			let res = await regetToken();
			if (!res) {
				if (isLoginUser()) {
					if (response.config.logout) {
						await response.config.logout(response);
					} else {
						const loadingInstance = ElLoading.service({ fullscreen: true, text: '你的登录过期了' });
						setTimeout(() => {
							Session.clear(); // 清除浏览器全部临时缓存
							Local.clear();
							window.location.href = '/login';
							loadingInstance.close();
						}, 1000);
					}
				} else {
					let res2 = await resetSignIn();
					if (!res2) {
						if (response.config.logout) {
							await response.config.logout(response);
						} else {
							const loadingInstance = ElLoading.service({ fullscreen: true, text: '你的登录过期了' });
							setTimeout(() => {
								Session.clear(); // 清除浏览器全部临时缓存
								Local.clear();
								window.location.href = '/login';
								loadingInstance.close();
							}, 1000);
						}
					}
				}
				return Promise.resolve(_service.interceptors.response);
			} else {
				return Promise.resolve({ code: retry });
			}
		} else if (res.code === 'LOGIC_DENY') {
			if (NotMessageCode.includes(res.subCode)) {
				return Promise.resolve(res);
			}

			if (res.msg.length !== 0) {
				ElMessage.error({ message: res.msg, grouping: true });
			}

			return Promise.resolve(res);
		}
		return Promise.resolve(res);
	},
	(error) => {
		if (error.message.indexOf('timeout') !== -1) {
			ElMessage.error({ message: '网络超时', grouping: true });
		} else if (error.message === 'Network Error') {
			ElMessage.error({ message: '网络连接错误', grouping: true });
		} else if (error.response?.data) {
			ElMessage.error({ message: error.response.statusText, grouping: true });
		} else {
			ElMessage.error({ message: '网络错误', grouping: true });
		}
		return Promise.resolve(error);
	}
);

// 导出 axios 实例
export default async function (cfg: any): Promise<resp> {
	const route = useRoute();
	if (route) {
		cfg.routePath = route.path;
		cfg.routeQuery = route.query;
	} else {
		cfg.routePath = '';
		cfg.routeQuery = {};
	}

	cfg.isRetry = false;
	let res: resp = await _service(cfg);

	if (res.code === retry) {
		cfg.isRetry = true;
		return await _service(cfg);
	}
	return res;
}
