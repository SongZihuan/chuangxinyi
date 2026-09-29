// versionUpdatePlugin.js
import fs from 'fs';
import path from 'path';
function handleType(oldVersion: string, type: string = 'release') {
	var oldVersionArr = oldVersion.split('.');
	//版本号第一位 如：1.2.3 则为 1
	var firstNum = +oldVersionArr[0];
	//版本号第二位 如：1.2.3 则为 2
	var secondNum = +oldVersionArr[1];
	//版本号第三位 如：1.2.3 则为 3
	var thirdNum = +oldVersionArr[2];
	switch (type) {
		case 'release':
			//release分支的处理逻辑
			++secondNum;
			thirdNum = 0;
			break;

		case 'hotfix':
			//hotfix分支的处理逻辑
			++thirdNum;
			break;

		default:
			// 默认按照最小版本处理
			++thirdNum;
			break;
	}
	return firstNum + '.' + secondNum + '.' + thirdNum;
}
export default () => {
	return {
		name: 'version-update',
		buildStart() {
			if (process.env.NODE_ENV === 'development') return;
			const url = path.join(process.cwd(), 'package.json');
			const json = fs.readFileSync(url, 'utf-8');
			const pkg = JSON.parse(json);
			// hotfix： 最小版本更新  release：穩定版本
			const version = handleType(pkg.version, 'hotfix');
			pkg.version = version;
			fs.writeFileSync('package.json', JSON.stringify(pkg, null, 2));
		},
	};
};
