<template>
  <div class="open">
    <div class="container layout-padding warp">
      <el-card class="box-card">
        <div class="card-block mb30">
          <div class="card-block-title">服务名称</div>
          <div class="card-block-content cursor-pointer" v-if="dataParams.title">
            <div class="title">{{dataParams.title}}</div>
          </div>
        </div>
        <div class="card-block mb15">
          <div class="card-block-title">服务开通说明</div>
          <div class="card-block-content description" v-if="dataParams.description">
            {{dataParams.description}}
          </div>
        </div>
        <div class="card-block mb15" v-for="(item, index) in dataParams.keyMap" :key="index">
          <div class="card-block-title">{{ item.label }}</div>
          <div class="card-block-content description" v-if="item.value">
            {{item.value}}
          </div>
        </div>
      </el-card>
      <el-card class="mt10">
        <div class="card-block">
          <div class="card-block-title">服务协议</div>
          <div class="card-block-content">
            <el-checkbox v-model="checked">我已阅读并同意</el-checkbox>
            <div class="agreement" @click="openDialog">《服务协议》</div>
          </div>
        </div>
      </el-card>
      <div  class="payment">
        <div class="payment-box">
          <div class="left"></div>
          <el-button type="primary" class="right" @click="onOpen">立即开通</el-button>
        </div>
      </div>
    </div>
    <agreementDialog ref="agreementDialogRef" />
  </div>
</template>
<script setup lang="ts">
import {useRoute, useRouter} from "vue-router";
import {encodeSearchParams} from "/@/utils/query";
import {useoauth2Api} from "/@/api/oauth2";
import {ElMessage} from "element-plus"
import {Session,Local} from "/@/utils/storage";
import {defineAsyncComponent, onMounted, ref} from "vue";
const agreementDialog = defineAsyncComponent(() => import('/@/views/oauth2/component/agreementDialog.vue')); //管理员协议组件

const route = useRoute()
const router = useRouter()
const agreementDialogRef = ref()

const openDialog = () => {
  agreementDialogRef.value.openDialog(dataParams.value.agreement)
}

let domain_uid = route.query?.domain as string
let redirect_uri = route.query?.redirect_uri as string
let next_redirect_uri = route.query?.next_redirect_uri as string
let params = route.query?.params as string
let isLoginToken = Number(route.query?.isLoginToken) == 1

const dataParams = ref<any>({
  title: '短链平台',
  description: '',
  keyMap: [],
  agreement: "",
})

const getDomainInfo = async () => {
  await useoauth2Api()
      .oauth2Domain({ domainUID: domain_uid })
      .then((res: any) => {
        if (res.code === "SUCCESS") {
          dataParams.value.title = res.data.name
          dataParams.value.description = res.data.describe
          dataParams.value.keyMap = res.data.keyMap
          dataParams.value.agreement = res.data.agreement
        }
      });
};

onMounted(async ()=>{
  await getDomainInfo()
})


const checkParams = async ()=> {
  if (!domain_uid || !redirect_uri) {
    await router.replace('/')
  }
  const paramsObj = JSON.parse(params)
  if (paramsObj.title) {
    dataParams.value.title = paramsObj.title
  }
}
checkParams()
const checked = ref<boolean>(false)
const onOpen = async () => {
  if(!checked.value) {
    ElMessage.error("请先阅读并同意服务协议")
    return
  }
  let res = await open()
  if (!res) {
    return
  }

  await successOpen()
}

const open = async (): Promise<boolean> => {
  return await useoauth2Api().oauth2Open({ webID: domain_uid }).then((res: any): boolean => {
    return res.code === "SUCCESS";
  })
}

const successOpen = async () => {
  let res = await useoauth2Api()
      .oauth2({ domainUID: domain_uid }, async () => {
        Session.clear()
        Local.clear()
        await router.push({
          path: '/login',
          query: {
            type: route.query?.type || "fromOauth2",
            params: params,
            next_redirect_uri: next_redirect_uri,
            redirect_uri: redirect_uri,
            domain: domain_uid,
            forceLogin: 1,
            isLoginToken: isLoginToken ? 1 : 0,
          }
        }); // 去登录页
      }, isLoginToken)
      .then(async (res: any): Promise<boolean> => {
        if (res.code === "SUCCESS") {
          let queryString = encodeSearchParams({
            token: res.data.token,
            subToken: res.data.subToken,
            params: params,
            redirect: next_redirect_uri,
          })

          // fromApplication不需要新open窗口
          window.location.href = `${redirect_uri}?${queryString}`
          return true
        }  else if (res.subCode === "NOT_OPEN_WEBSITE") {
          ElMessage.error("开通失败")
          return true  // 不返回false，不跳转
        } else {
          return false
        }
      });
  if (!res) {
    await router.replace("/")
  }
}
</script>
<style scoped lang="scss">
.open {
  padding: 20px;
}
.box-card{
  width: 100%;
  .title {
    font-size: 18px;
    color: #373d41;
    border-bottom: 1px solid #ebeef5;
    font-weight: 600;
  }
  .purchase-box {
    padding: 20px 0;
  }
}
.card-block {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  .card-block-title {
    width: 150px;
    text-align: left;
    font-size: 14px;
    color: #373d41;
    font-weight: 600;
    padding-right: 40px;
  }

  .description {
    font-size: 12px;
    color: #787e80;
  }
  .card-block-content {
    width: 100%;
    display: flex;
    justify-content: flex-start;
    align-items: center;
    .agreement {
      color: #409eff;
      cursor: pointer;
    }
    .title {
      text-align: center;
      font-size: 14px;
      color: #373d41;
      font-weight: 600;
      padding: 5px 30px;
      background: rgba(0,112,204,.15);
    }
    .el-select {
      width: 240px;
      .el-input {
        width: 100%;
        .el-input__inner {
          width: 100%;
        }
      }
    }
  }
}
.payment {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 100px;
  background-color: #fff;
  box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.1);
  .payment-box {
    height: 100%;
    display: flex;
    justify-content: flex-end;
    align-items: center;
    padding: 0 20px;
    .left {
      font-size: 12px;
      color: #373d41;
      padding: 0 30px;
      .money {
        font-size: 24px;
        color: #ff8a00;
      }
    }
    .right {
      height: 40px;
      line-height: 40px;
      text-align: center;
      border-radius: 0;
      padding: 0 16px;
      cursor: pointer;
    }
  }
}
.table-money {
  color: #ff8a00;
}
</style>