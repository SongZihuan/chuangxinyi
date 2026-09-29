<template>
	<el-card shadow="never">
		<el-descriptions direction="vertical" title="控制台应用列表"></el-descriptions>
		<div class="app-warp">
			<el-row :gutter="10">
				<el-col :xs="24" :sm="3" :md="3" :lg="3" :xl="3" v-for="(item, index) in appList" :key="index" class="mb20" @click="goto(item)">
					<div class="app-item">
						<div class="app-item__icon">
							<SvgIcon :name="item.icon" :size="28"></SvgIcon>
						</div>
						<el-text>{{ item.name }}</el-text>
					</div>
				</el-col>
			</el-row>
		</div>
	</el-card>
</template>
<script setup lang="ts">
import { ref } from 'vue';
import { useApplicationApi } from '/@/api/system/application';
import { useoauth2Api } from '/@/api/oauth2';
import { Local, Session } from '/@/utils/storage';
import { encodeSearchParams } from '/@/utils/query';
import { useRouter } from 'vue-router';

const router = useRouter();

const applicationApi = useApplicationApi();
const appList = ref<any[]>([]);

const getList = () => {
	applicationApi.applicationAll().then((res: any) => {
		if (res.code === 'SUCCESS') {
			appList.value = res.data.application as any[];
		}
	});
};
const goto = async (item: any) => {
	const webID = item.webUid as string;
	let url = new URL(item.url as string);

	let params = url.searchParams.get('params');

	if (!params) {
		params = JSON.stringify({ name: item.name });
	}

	let next_redirect_uri = url.searchParams.get('redirect');
	if (!next_redirect_uri) {
		next_redirect_uri = '/home';
	}

	let redirect_uri = `${url.protocol}//${url.host}${url.pathname}`;
	if (Session.get('token')) {
		let res = await useoauth2Api()
			.oauth2(
				{ domainUID: webID },
				async () => {
					Session.clear();
					Local.clear();
					await router.push({
						path: '/login',
						query: {
							type: 'fromApplication',
							params: params,
							next_redirect_uri: next_redirect_uri,
							redirect_uri: redirect_uri,
							domain: webID,
							forceLogin: 1,
							isLoginToken: 0,
						},
					}); // 去登录页
				},
				false
			)
			.then(async (res: any): Promise<boolean> => {
				if (res.code === 'SUCCESS') {
					let queryString = encodeSearchParams({
						token: res.data.token,
						subToken: res.data.subToken,
						params: params,
						redirect: next_redirect_uri,
					});

					window.open(`${redirect_uri}?${queryString}`);
					return true;
				} else if (res.subCode === 'NOT_OPEN_WEBSITE') {
					window.open(
						router.resolve({
							path: '/oauth2/open',
							query: {
								type: 'fromApplication',
								params: params,
								next_redirect_uri: next_redirect_uri,
								redirect_uri: redirect_uri,
								domain: webID,
								isLoginToken: 0,
							},
						}).href
					); // 新窗口打开
					return true;
				} else {
					return false;
				}
			});
		if (!res) {
			await router.replace('/');
		}
	} else {
		await router.push({
			path: '/login',
			query: {
				type: 'fromApplication',
				params: params,
				next_redirect_uri: next_redirect_uri,
				redirect_uri: redirect_uri,
				domain: webID,
				forceLogin: 1,
				isLoginToken: 0,
			},
		}); // 去登录页
	}
};

getList();
</script>
<style scoped lang="scss">
.app-warp {
	.app-item {
		padding: 10px;
		display: flex;
		justify-content: center;
		align-items: center;
		flex-direction: column;
		cursor: pointer;
		background: #f5f5ff;
		height: 80px;
		span {
			margin-top: 4px;
		}

		&:hover {
			font-weight: bold;
			font-size: 16px;
		}
	}
}
</style>
