export function isWeiXin(): boolean {
	return /(micromessenger)/i.test(navigator.userAgent);

}
