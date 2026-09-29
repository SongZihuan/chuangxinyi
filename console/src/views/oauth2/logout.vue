<template>
	<div></div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Session, Local } from '/@/utils/storage';
import { useoauth2Api } from '/@/api/oauth2';
import { encodeSearchParams } from '/@/utils/query'
import {ElMessageBox} from "element-plus";

const route = useRoute();
const router = useRouter();
const redirect_uri = ref('');
const domain_uid = ref('')
const params = ref('')

const init = async () => {
	if (route.query?.redirect_uri && route.query?.domain) {
    redirect_uri.value = route.query.redirect_uri as string;
    domain_uid.value = route.query.domain as string

    if (route.query?.params) {
      params.value = route.query?.params as string
    } else {
      params.value = JSON.stringify({})
    }

    await getDomainInfo(domain_uid.value)
    return true
	} else {
    await router.push({
      path: "/404"
    })
    return false
  }
};

const domainInfo = ref({ name: '', describe: '', agreement: '', keyMap : [] as any[] });

//获取站点信息
const getDomainInfo = async (url: string) => {
  await useoauth2Api()
      .oauth2Domain({ domainUID: url })
      .then(async (res: any) => {
        if (res.code === "SUCCESS") {
          domainInfo.value = res.data;
        } else {
          await router.push("/404")
        }
      });
};

const logout = async () => {
  const nowToken = Session.get('token');
  const loginToken = Session.get('login-token');
  let sure = false

  if (!nowToken || !loginToken || nowToken === loginToken.token) {
    await ElMessageBox.confirm(
      `是否确认退出${import.meta.env.VITE_PROJECT_NAME}和${domainInfo.value.name}？`,
      'Warning',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
    ).then(() => {
      sure = true
      Session.clear()
      Local.clear()
    }).catch(() => {
      sure = false
    })
  } else {
    await ElMessageBox.confirm(
        `是否确认退出${domainInfo.value.name}的子账号登录？`,
        'Warning',
        {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning',
        }
    ).then(() => {
      sure = true
    }).catch(() => {
      sure = false
    })
  }

  let queryString = encodeSearchParams({
    logout: sure ? 1 : 0,
  })
  window.location.href = `${redirect_uri.value}?${queryString}`;
};

onMounted(async ()=>{
  if (await init()) {
    await logout()
  }
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
