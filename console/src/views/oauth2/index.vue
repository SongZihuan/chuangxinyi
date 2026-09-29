<template>
	<div></div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Session } from '/@/utils/storage';
import { useoauth2Api } from '/@/api/oauth2';
import { NextLoading } from '/@/utils/loading';
import { encodeSearchParams } from '/@/utils/query'
import {ElLoading} from "element-plus";

const route = useRoute();
const router = useRouter();
const redirect_uri = ref('');
const redirect_domain = ref('')
const domain_uid = ref('')
const params = ref('')
const next_redirect_uri = ref('')

const init = async () => {
	if (route.query?.redirect_uri && route.query?.domain) {
    redirect_uri.value = route.query.redirect_uri as string;
    domain_uid.value = route.query.domain as string

    if (route.query?.params) {
      params.value = route.query?.params as string
    } else {
      params.value = JSON.stringify({})
    }

    if (route.query?.redirect) {
      next_redirect_uri.value = route.query?.redirect as string
    } else {
      next_redirect_uri.value = ''
    }

    try {
      redirect_domain.value = (new URL(redirect_uri.value)).hostname
    } catch {
      return false
    }

    return true
	} else {
    await router.push({
      path: "/404"
    })
    return false
  }
};

const login = async () => {
	if (Session.get('token')) {
		let res = await useoauth2Api()
			.oauth2({ domainUID: domain_uid.value }, async () => {
        await router.push({
          path: '/login',
          query: {
            type: "fromOauth2",
            params: params.value,
            next_redirect_uri: next_redirect_uri.value,
            redirect_uri: redirect_uri.value,
            domain: domain_uid.value,
            forceLogin: 1,
          }
        }); // 去登录页
      }, true)
			.then(async (res: any): Promise<boolean> => {
				if (res.code === "SUCCESS") {
          let queryString = encodeSearchParams({
            token: res.data.token,
            subToken: res.data.subToken,
            params: params.value,
            redirect: next_redirect_uri.value,
          })

					window.location.href = `${redirect_uri.value}?${queryString}`;
          return true
				} else if (res.subCode === "NOT_OPEN_WEBSITE") {
          await router.push({
            path: '/oauth2/open',
            query: {
              type: "fromOauth2",
              params: params.value,
              next_redirect_uri: next_redirect_uri.value,
              redirect_uri: redirect_uri.value,
              domain: domain_uid.value,
              isLoginToken: 1,
            }
          }); // 去登录页
          return true
        } else {
          return false
        }
			});
    if (!res) {
      await router.push({
        path: '/login',
        query: {
          type: "fromOauth2",
          params: params.value,
          next_redirect_uri: next_redirect_uri.value,
          redirect_uri: redirect_uri.value,
          domain: domain_uid.value,
          forceLogin: 1,
          isLoginToken: 1,
        }
      }); // 去登录页
    }
	} else {
    await router.push({
      path: '/login',
      query: {
        type: "fromOauth2",
        params: params.value,
        next_redirect_uri: next_redirect_uri.value,
        redirect_uri: redirect_uri.value,
        domain: domain_uid.value,
        forceLogin: 1,
        isLoginToken: 1,
      }
    }); // 去登录页
  }
};

onMounted(async ()=>{
  const loadingInstance = ElLoading.service({ fullscreen: true, text: '正在登入，请稍后...' });
  if (await init()) {
    await login()
  }
  NextLoading.done();
  loadingInstance.close()
})

</script>

<style lang="scss" scoped>
:deep(.el-divider--horizontal) {
	margin: 15px 0px;
}
@media screen and (max-width: 769px) {
	.coantainer .content-box {
		width: 340px !important;
	}
}
.coantainer {
	width: 100%;
	height: 100%;
	background-color: rgba(251, 252, 255);
	.link-wrapper {
		padding: 15px;
		padding-top: 15vh;
		display: flex;
		flex-direction: column;
		align-items: center;
		.logo-warp {
			width: 100%;
			display: flex;
			justify-content: center;
			.logo {
				width: 100px;
				height: auto;
				margin-top: 30px;
			}
		}
		.content-box {
			padding: 15px;
			border: 1px solid #dce3e8;
			background: #fff;
			border-radius: 4px;
			width: 540px;
			margin-top: 20px;
			.content-title {
				font-size: 18px;
				color: #40485b;
				margin-bottom: 10px;
				text-align: center;
				font-weight: bold;
			}
			.content-item {
				padding: 8px 0px;
				color: #40485b;
				margin-left: 15px;
				line-height: 20px;
			}
			.btn {
				display: flex;
				flex-direction: row-reverse;
			}
		}
	}
}
</style>
