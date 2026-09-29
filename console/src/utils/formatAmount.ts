/**
 * 将金额格式化分位转为元
 * @param amount
 */
export function formatAmount(amount: number | string): string {
	if(!amount) return '0.00';
	const formattedAmount = (Number(amount) / 100).toFixed(2);
	return formattedAmount;
}

/**
 * 将金额格式化分位转为元
 */
export function formatAmountToCent(amount: number | string): number {
	const formattedAmount: number = Number(amount) * 100 | 0;
	return formattedAmount;
}

/**
 * 根据传入状态转换为对应的发票名称
 * 1-个人普票 2-企业普票 3-企业专票
 */
export function formatInvoiceType(type: number): string {
	switch (type) {
		case 1:
			return '个人普票';
		case 2:
			return '企业普票';
		case 3:
			return '企业专票';
		case 5:
			return '减少开票额度';
		case 6:
			return '增加开票额度';
		default:
			return '未知';
	}
}
