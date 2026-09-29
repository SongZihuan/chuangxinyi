<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="layout-padding-auto">
			<div class="search-header mb15">
				<el-form :inline="true" :model="state.tableData.param">
					<!--          根据token搜索,根据token查询-->
					<el-row>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
							<el-form-item style="width: 80%">
								<el-input text v-model="state.tableData.param.token" placeholder="请输入token" />
							</el-form-item>
						</el-col>
						<el-col :xs="24" :sm="12" :md="12" :lg="12" :xl="12">
							<el-form-item style="width: 80%">
								<el-input text v-model="state.tableData.param.requestid" placeholder="请输入请求ID" />
							</el-form-item>
						</el-col>
						<!--          sql编辑器-->
						<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
							<el-form-item>
								<Codemirror
									v-model:value="state.tableData.param.cond"
									:options="cmOptions"
									border
									:width="sqlEditStyle.width"
									:height="100"
									placeholder="请输入查询sql例如(id = 1)"
								/>
							</el-form-item>
						</el-col>

						<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
							<el-form-item>
								<el-button type="primary" class="ml10" @click="handleSearch('token')">
									<el-icon>
										<ele-Search />
									</el-icon>
									根据token查询
								</el-button>
							</el-form-item>
							<el-form-item>
								<el-button type="primary" class="ml10" @click="handleSearch('request_id')">
									<el-icon>
										<ele-Search />
									</el-icon>
									根据请求ID查询
								</el-button>
							</el-form-item>
							<!--              根据sql查询-->
							<el-form-item>
								<el-button type="primary" class="ml10" @click="handleSearch('sql')">
									<el-icon>
										<ele-Search />
									</el-icon>
									根据sql查询
								</el-button>
							</el-form-item>
							<el-form-item>
								<el-button type="primary" class="ml10" @click="handleSearch('')">
									<el-icon>
										<ele-Search />
									</el-icon>
									重置
								</el-button>
							</el-form-item>
						</el-col>
					</el-row>
				</el-form>
			</div>
			<el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" border>
				<el-table-column label="ID(id)" prop="id" width="80" show-overflow-tooltip align="center">
					<template #default="scope">
						{{ scope.row.id || '-' }}
					</template>
				</el-table-column>
				<el-table-column label="请求ID(request_id)" prop="request_id" width="200" align="center">
					<template #default="scope">
						{{ scope.row.request_id || '-' }}
					</template>
				</el-table-column>
				<el-table-column label="请求ID前缀(request_id_prefix)" prop="request_id_prefix" width="220" align="center">
					<template #default="scope">
						{{ scope.row.request_id_prefix || '-' }}
					</template>
				</el-table-column>
				<el-table-column label="服务器名称(server_name)" prop="server_name" width="120" align="center">
					<template #default="scope">
						{{ scope.row.server_name || '-' }}
					</template>
				</el-table-column>
				<el-table-column label="用户ID(user_id)" prop="user_id" width="80" align="center">
					<template #default="scope">
						<span>{{ scope.row.user_id || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="用户UID(user_uid)" prop="user_uid" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.user_uid || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="角色名称(role_name)" prop="role_name" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.role_name || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="角色标识(role_sign)" prop="role_sign" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.role_sign || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="网站名称(web_name)" prop="web_name" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.web_name || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="网站域名(web_domain)" prop="web_domain" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.web_domain || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="请求网站ID(requests_web_id)" prop="requests_web_id" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.requests_web_id || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="请求网站名称(requests_web_name)" prop="requests_web_name" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.requests_web_name || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="请求网站域名(requests_web_domain)" prop="requests_web_domain" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.requests_web_domain || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="IP1(ip1)" prop="ip1" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.ip1 || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="IP2(ip2)" prop="ip2" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.ip2 || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="端口1(port1)" prop="port1" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.port1 || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="端口2(port2)" prop="port2" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.port2 || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="地理位置编码(geo_code)" prop="geo_code" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.geo_code || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="地理位置(geo)" prop="geo" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.geo || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="协议(scheme)" prop="scheme" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.scheme || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="请求方法(method)" prop="method" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.method || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="主机(host)" prop="host" width="120" align="center">
					<template #default="scope">
						<span>{{ scope.row.host || '-' }}</span>
					</template>
				</el-table-column>
				<el-table-column label="路径(path)" prop="path" width="120" align="center" />
				<el-table-column label="查询(query)" prop="query" width="100" align="center">
					<template #default="scope">
						<el-button type="text" size="small" @click="showJson(scope.row.query, '查询')"> 查询</el-button>
					</template>
				</el-table-column>
				<el-table-column label="请求体(requests_body)" prop="requests_body" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-button type="text" size="small" @click="showJson(scope.row.requests_body, '请求体')"> 请求体</el-button>
					</template>
				</el-table-column>
				<el-table-column label="响应体(response_body)" prop="response_body" show-overflow-tooltip width="100" align="center">
					<template #default="scope">
						<el-button type="text" size="small" @click="showJson(scope.row.response_body, '响应体')"> 响应体</el-button>
					</template>
				</el-table-column>
				<el-table-column label="响应体错误(response_body_error)" prop="response_body_error" show-overflow-tooltip width="100" align="center" />
				<el-table-column label="请求头(requests_header)" prop="requests_header" width="100" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-button type="text" size="small" @click="showJson(scope.row.requests_header, '请求头')"> 请求头 </el-button>
					</template>
				</el-table-column>
				<el-table-column label="响应头(response_header)" prop="response_header" width="100" show-overflow-tooltip align="center">
					<template #default="scope">
						<el-button type="text" size="small" @click="showJson(scope.row.response_header, '响应头')"> 响应头 </el-button>
					</template>
				</el-table-column>
				<el-table-column label="状态码(status_code)" prop="status_code" width="100" align="center">
					<template #default="scope">
						<el-tag v-if="scope.row.status_code >= 200 && scope.row.status_code < 300" type="success">
							{{ scope.row.status_code }}
						</el-tag>
						<el-tag v-else type="warning">{{ scope.row.status_code }}</el-tag>
					</template>
				</el-table-column>
				<el-table-column label="错误(panic_error)" prop="panic_error" width="120" align="center">
					<template #default="scope">
						<span class="error">{{ scope.row.panic_error }}</span>
					</template>
				</el-table-column>
				<el-table-column label="消息(message)" prop="message" width="120" align="center">
					<template #default="scope">
						<el-button type="text" size="small" @click="showJson(scope.row.message, '消息')"> 消息</el-button>
					</template>
				</el-table-column>
				<el-table-column label="耗时(use_time)" prop="use_time" width="120" align="center">
					<!--          显示毫秒-->
					<template #default="scope">
						<span>{{ scope.row.use_time }}ms</span>
					</template>
				</el-table-column>
				<el-table-column label="创建时间(create_at)" prop="create_at" width="120" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.create_at).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="开始时间(start_at)" prop="start_at" width="120" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.start_at).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="结束时间(end_at)" prop="end_at" width="120" align="center">
					<template #default="scope">
						<span>{{ dayjs.unix(scope.row.end_at).format('YYYY-MM-DD HH:mm:ss') }}</span>
					</template>
				</el-table-column>
				<el-table-column label="操作" width="200" fixed="right" align="center">
					<template #default="scope">
						<el-button type="text" size="small" @click="copyText(scope.row.user_token)">复制token</el-button>
						<el-button type="text" size="small" @click="copyText(scope.row.request_id)">复制请求ID</el-button>
					</template>
				</el-table-column>
			</el-table>
			<el-pagination hide-on-single-page
				@size-change="onHandleSizeChange"
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				:page-sizes="[10, 20, 30]"
				v-model:current-page="state.tableData.param.page"
				background
				v-model:page-size="state.tableData.param.pagesize"
				layout="total, sizes, prev, pager, next, jumper"
				:total="state.tableData.total"
			>
			</el-pagination>
		</el-card>
		<!-- json展示弹窗 -->
		<el-dialog :title="dialog.title" v-model="dialog.isShowDialog" width="769px" height="500px" center destroy-on-close>
			<jsonEditor :jsonData="jsonData" v-if="isJson(currentInfo)" />
			<el-card shadow="never" v-else>
				<template #header>
					<div class="card-header">其他数据格式</div>
				</template>
				<span>{{ currentInfo }}</span>
			</el-card>
		</el-dialog>
	</div>
</template>

<script setup lang="ts">
import { reactive, onMounted, watch, defineAsyncComponent, ref, computed } from 'vue';
import dayjs from 'dayjs';
import { useAccessrecordApi } from '/@/api/accessrecord';
import Codemirror from 'codemirror-editor-vue3';
// placeholder
import 'codemirror/lib/codemirror.css';
import 'codemirror/addon/display/placeholder.js';
import 'codemirror/mode/sql/sql.js';
import 'codemirror/theme/dracula.css';
import commonFunction from '/@/utils/commonFunction';

const jsonEditor = defineAsyncComponent(() => import('/@/components/json-editor/index.vue'));
const state = reactive<any>({
	tableData: {
		data: [],
		total: 0,
		loading: false,
		param: {
			page: 1,
			pagesize: 20,
			timetype: '',
			starttime: 0,
			endtime: 0,
			range: [],
			token: '',
			requestid: '',
			cond: '',
		},
	},
});
const cmOptions = {
	mode: 'text/x-sparksql', // Language mode
	theme: 'default', // Theme
	//代码折叠
	lineWrapping: true,
	foldGutter: true,
	gutters: ['CodeMirror-linenumbers', 'CodeMirror-foldgutter'],
};
const dialog = reactive({
	title: 'JSON',
	isShowDialog: false,
});
const jsonData = ref();
// 判断的是否是JSON字符串
const isJson = (str: any) => {
	if (typeof str == 'string') {
		try {
			var obj = JSON.parse(str);
			// 等于这个条件说明就是JSON字符串 会返回true
			if (typeof obj == 'object' && obj) {
				return true;
			} else {
				//不是就返回false
				return false;
			}
		} catch (e) {
			return false;
		}
	}
	return false;
};
let searchType = ref('');
const handleSearch = (type) => {
	state.tableData.param.page = 1;
	searchType.value = type;
	if (type == 'token' && state.tableData.param.token != '') {
		getTokenTableData();
	} else if (type == 'request_id' && state.tableData.param.requestid != '') {
		getRequestIdTableData();
	} else if (type == 'sql' && state.tableData.param.cond != '') {
		getSqllTableData();
	} else {
		resetSearch();
	}
};
const resetSearch = () => {
	state.tableData.param.token = '';
	state.tableData.param.page = 1;
	state.tableData.param.requestid = '';
	state.tableData.param.cond = '';
	getTableData();
};
// 复制token
const currentInfo = ref();
const showJson = (row: any, title: string) => {
	currentInfo.value = row;
	if (isJson(row)) {
		jsonData.value = JSON.parse(row);
		dialog.title = title + 'JSON数据';
	}
	dialog.title = title + '数据';
	dialog.isShowDialog = true;
};
// sql编辑器长度
const { isMobile, copyText } = commonFunction();
const currentMobile = computed(() => {
	return isMobile();
});
const sqlEditStyle = computed(() => {
	let style = {};
	if (currentMobile.value) {
		style = {
			width: '100%',
		};
	} else {
		// 监听页面宽度设置sql编辑器长度
		style = {
			width: document.body.clientWidth - 300,
		};
	}
	return style;
});

watch(
	() => state.tableData.param.range,
	() => {
		if (state.tableData.param.range.length > 0) {
			state.tableData.param.starttime = dayjs(state.tableData.param.range[0]).unix();
			state.tableData.param.endtime = dayjs(state.tableData.param.range[1]).unix();
		}
	},
	{
		deep: true,
	}
);
// 初始化表格数据
const getTableData = () => {
	state.tableData.loading = true;
	useAccessrecordApi()
		.accessrecordList(state.tableData.param)
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.record;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
// 根据sql获取访问记录
const getSqllTableData = () => {
	state.tableData.loading = true;
	useAccessrecordApi()
		.accessrecordSqlGet({ cond: state.tableData.param.cond, ...state.tableData.param })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.record;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
// 根据token获取访问记录
const getTokenTableData = () => {
	state.tableData.loading = true;
	useAccessrecordApi()
		.accessrecordTokenGet({ token: state.tableData.param.token, ...state.tableData.param })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				state.tableData.data = res.data.record;
				state.tableData.total = res.data.count;
				state.tableData.loading = false;
			}
		});
};
// 根据请求ID获取访问记录
const getRequestIdTableData = () => {
	state.tableData.loading = true;
	useAccessrecordApi()
		.accessrecordIdGet({ requestID: state.tableData.param.requestid, ...state.tableData.param })
		.then((res: any) => {
			if (res.code === "SUCCESS") {
				if (res.data.find) {
					state.tableData.data = [res.data.record];
				} else {
					state.tableData.data = [];
				}

				state.tableData.total = 1;
				state.tableData.loading = false;
			}
		});
};
// 分页改变
const onHandleSizeChange = (val: number) => {
	state.tableData.param.pagesize = val;
	getTableData();
};
// 分页改变
const onHandleCurrentChange = (val: number) => {
	state.tableData.param.page = val;
	getTableData();
};
// 页面加载时
onMounted(() => {
	getTableData();
});
</script>

<style scoped lang="scss"></style>
