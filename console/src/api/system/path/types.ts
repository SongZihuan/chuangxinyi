export const pathModeDict = [
	{ label: '前缀匹配', value: 1 },
	{ label: '完全匹配', value: 2 },
	{ label: '正则匹配', value: 3 },
];
export const captchaModeDict = [
	{ label: '不需要验证码', value: 1 },
	{ label: '需要验证码', value: 2 },
	{ label: '需要验证码,仅限制静默', value: 3 },
	{ label: '需要验证码,仅限制滑块', value: 4 },
];
export const corsModeDict = [
	{ label: '允许跨域', value: 1 },
	{ label: '允许外站跨域', value: 2 },
	{ label: '限制跨域', value: 3 },
];
export const busyModeDict = [
	{ label: 'IP限制模式', value: 1 },
	{ label: '用户限制模式', value: 2 },
];
export const methodDict = [
	{ label: 'GET', value: 'GET' },
	{ label: 'POST', value: 'POST' },
];
export const adminModeDict = [
	{ label: '非管理员接口', value: 1 },
	{ label: '外站管理员', value: 2 },
	{ label: '用户中心管理员', value: 3 },
];
