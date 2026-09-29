// 通用函数
import useClipboard from 'vue-clipboard3';
import { ElMessage } from 'element-plus';
import { formatDate } from '/@/utils/formatTime';

export default function () {
	const { toClipboard } = useClipboard();

	// 百分比格式化
	const percentFormat = (row: EmptyArrayType, column: number, cellValue: string) => {
		return cellValue ? `${cellValue}%` : '-';
	};
	// 列表日期时间格式化
	const dateFormatYMD = (row: EmptyArrayType, column: number, cellValue: string) => {
		if (!cellValue) return '-';
		return formatDate(new Date(cellValue), 'YYYY-mm-dd');
	};
	// 列表日期时间格式化
	const dateFormatYMDHMS = (row: EmptyArrayType, column: number, cellValue: string) => {
		if (!cellValue) return '-';
		return formatDate(new Date(cellValue), 'YYYY-mm-dd HH:MM:SS');
	};
	// 列表日期时间格式化
	const dateFormatHMS = (row: EmptyArrayType, column: number, cellValue: string) => {
		if (!cellValue) return '-';
		let time = 0;
		if (typeof row === 'number') time = row;
		if (typeof cellValue === 'number') time = cellValue;
		return formatDate(new Date(time * 1000), 'HH:MM:SS');
	};
	// 小数格式化
	const scaleFormat = (value: string = '0', scale: number = 4) => {
		return Number.parseFloat(value).toFixed(scale);
	};
	// 小数格式化
	const scale2Format = (value: string = '0') => {
		return Number.parseFloat(value).toFixed(2);
	};
	//生成32位随机数(大小写数字)
	const uuid = (num: number) => {
		var str = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
		var result = '';
		for (var i = num; i > 0; --i) result += str[Math.floor(Math.random() * str.length)];

		return result;
	};
	// 点击复制文本
	const copyText = (text: string) => {
		return new Promise((resolve, reject) => {
			try {
				// 复制
				toClipboard(text);
				// 下面可以设置复制成功的提示框等操作
				ElMessage.success('复制成功');
				resolve(text);
			} catch (e) {
				// 复制失败
				ElMessage.error('复制失败');
				reject(e);
			}
		});
	};
	const flatToTree = (flatData: any, idProp: string, parentIdProp: string) => {
		const treeData: any = [];
		const idMap: any = {};
		// 将所有节点存储到 idMap 中，方便快速查找节点
		for (const node of flatData) {
			idMap[node[idProp]] = node;
			node.children = [];
		}
		// 遍历所有节点，将它们添加到树形结构中
		for (const node of flatData) {
			const parent = idMap[node[parentIdProp]];
			if (parent) {
				parent.children.push(node);
			} else {
				treeData.push(node);
			}
		}
		return treeData;
	};
	//判断手机端
	const isMobile = () => {
		var userAgentInfo = navigator.userAgent;
		var mobileAgents = ['Android', 'iPhone', 'SymbianOS', 'Windows Phone', 'iPad', 'iPod'];
		var mobile_flag = false;
		//根据userAgent判断是否是手机
		for (var v = 0; v < mobileAgents.length; v++) {
			if (userAgentInfo.indexOf(mobileAgents[v]) > 0) {
				mobile_flag = true;
				break;
			}
		}
		var screen_width = window.screen.width;
		var screen_height = window.screen.height;

		//根据屏幕分辨率判断是否是手机
		if (screen_width > 325 && screen_height < 750) {
			mobile_flag = true;
		}

		return mobile_flag;
	};
	const isImage = (fileName: string) => {
		if(!fileName) return false;
		const imageExtensions = ['jpg', 'jpeg', 'png', 'gif','svg'];
		if(imageExtensions.includes(fileName.toLowerCase())) return true;
		const extension: string | undefined = fileName.split('.').pop();
		if (!extension) {
			return;
		}
		if (imageExtensions.includes(extension.toLowerCase())) {
			return true;
		}
		return;
	};
	return {
		percentFormat,
		dateFormatYMD,
		dateFormatYMDHMS,
		dateFormatHMS,
		scaleFormat,
		scale2Format,
		copyText,
		uuid,
		flatToTree,
		isMobile,
		isImage
	};
}
