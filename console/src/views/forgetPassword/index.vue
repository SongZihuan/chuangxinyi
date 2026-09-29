<template>
	<div class="login-container flex">
    <div class="login-left">
      <div class="login-left-img">
        <img :src="loginMain" />
      </div>
    </div>
		<div class="login-right flex">
			<div class="login-right-warp flex-margin">
				<span class="login-right-warp-one"></span>
				<span class="login-right-warp-two"></span>
				<div class="login-right-warp-mian">
					<div class="logo-warp"><img :src="logo" class="logo" /></div>
					<div class="login-right-warp-main-title">找回密码</div>
					<div class="login-right-warp-main-form" v-if="isEntry">
						<div class="login-right-warp-main-form_btns">
							<div class="personal">
								<el-button round v-waves type="primary" @click="clickEntrance('personal')">
									<el-icon> <ele-UserFilled /> </el-icon>使用人找回</el-button
								>
							</div>
							<div class="enterprise">
								<el-button round v-waves type="success" @click="clickEntrance('business')">
									<el-icon> <ele-Avatar /> </el-icon>法人找回</el-button
								>
							</div>
						</div>

            <div class="register">
              <div>想起密码?</div>
              <div @click="goLogin">立即登录</div>
            </div>

					</div>
					<div class="login-right-warp-main-form" v-else>
						<el-steps :active="stepActive" align-center>
							<el-step title="验证码检验" />
							<el-step title="信息填写" />
							<el-step title="个人认证" />
						</el-steps>
						<template v-if="stepActive == 1">
							<personal v-if="tabsActiveName == 'personal'" @resetBack="onResetBack" />
							<business v-if="tabsActiveName == 'business'" @resetBack="onResetBack" />
						</template>
						<template v-if="stepActive == 2">
							<personal-form
								v-if="tabsActiveName == 'personal'"
								:ruleFormPersonal="ruleFormPersonal"
								@resetBackForm="resetBackFormPersonal"
								@go-back="goBack"
							/>
							<business-form
								v-if="tabsActiveName == 'business'"
								:ruleFormBusinese="ruleFormBusinese"
								@resetBackForm="resetBackFormBusiness"
								@go-back="goBack"
							/>
						</template>
						<template v-if="stepActive == 3">
							<scan-qr-codes
								v-if="tabsActiveName == 'personal'"
								:current-user="tabsActiveName"
								:scan-qr-data="ruleFormPersonal"
								@go-back="goBack"
								@success="scanPersonalSuccess"
							></scan-qr-codes>
							<scan-qr-codes
								v-if="tabsActiveName == 'business'"
								:current-user="tabsActiveName"
								:scan-qr-data="ruleFormBusinese"
								@go-back="goBack"
								@success="scanBusinessSuccess"
							></scan-qr-codes>
						</template>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts" name="forgetpasswordIndex">
import { onMounted, ref } from 'vue';
import { NextLoading } from '/@/utils/loading';
import {Local, Session} from '/@/utils/storage';
import { ElLoading, ElMessage } from 'element-plus';
import { useUserInfo } from '/@/stores/userInfo';
import Personal from '/@/views/forgetPassword/component/personal.vue';
import Business from '/@/views/forgetPassword/component/business.vue';
import PersonalForm from '/@/views/forgetPassword/component/personalForm.vue';
import BusinessForm from '/@/views/forgetPassword/component/businessForm.vue';
import ScanQrCodes from '/@/views/forgetPassword/component/scanQrCodes.vue';
import { useLoginApi } from '/@/api/login';
import { businessInfoFormTypes, personInfoFormTypes } from '/@/api/register/types';
import { useLoginSignIn } from '/@/hooks/useLoginSignIn';
import useFile from '/@/hooks/useFile';
import {useRoute, useRouter} from "vue-router";
import {useoauth2Api} from "/@/api/oauth2";
import {encodeSearchParams} from "/@/utils/query";

const logo = useFile().getFile('logo');
const loginMain = useFile().getFile('login-bg');
const stores = useUserInfo();
const tabsActiveName = ref('business');
const { onSignIn } = useLoginSignIn();

const route = useRoute()
const router = useRouter()

// 是否是入口进入
const isEntry = ref<boolean>(true);
const isRest = ref<boolean>(false);
// 步骤条
const stepActive = ref<number>(1);
// 个人认证
const ruleFormPersonal = ref<personInfoFormTypes>({
	phoneToken: '',
	id: '',
	userName: '',
	userIdCard: '',
	isCompany: false,
	companyName: '',
	companyID: '',
	legalPersonName: '',
	faceToken: '',
});
// 企业认证
const ruleFormBusinese = ref<businessInfoFormTypes>({
	phoneToken: '',
	id: '',
	companyID: '',
	companyName: '',
	faceToken: '',
	legalPersonName: '',
	legalPersonID: '',
});
// 入口点击事件
const clickEntrance = (type: string) => {
	tabsActiveName.value = type;
	isEntry.value = false;
	stepActive.value = 1;
};
// 上一步
const goBack = (targetNumber: number) => {
	stepActive.value = targetNumber;
};

const goLogin = () => {
  router.push({ path: '/login', query: route.query});
};

// 第一步数据返回
const onResetBack = async (token: string) => {
	// 存储 token 到浏览器缓存
	Session.set('phoneToken', token);
	stores.setxtoken(token);
	isRest.value = false;
	stepActive.value = 2;
	// 模拟数据，对接接口时，记得删除多余代码及对应依赖的引入。用于 `/src/stores/userInfo.ts` 中不同用户登录判断（模拟数据）
};
// 第二步数据接收
const resetBackFormPersonal = (rulesForm: personInfoFormTypes) => {
	ruleFormPersonal.value = rulesForm;
	stepActive.value = 3;
};
const resetBackFormBusiness = (rulesForm: businessInfoFormTypes) => {
	ruleFormBusinese.value = rulesForm;
	stepActive.value = 3;
};
// 个人认证成功
const scanPersonalSuccess = async (faceToken: string) => {
	if (!faceToken) {
		ElMessage.warning('人脸识别失败');
		return;
	}
	if (!Session.get('phoneToken')) {
		ElMessage.warning('请先获取验证码');
		return;
	}
	ruleFormPersonal.value.phoneToken = Session.get('phoneToken');
	ruleFormPersonal.value.faceToken = faceToken;
	const loadingInstance = ElLoading.service({
		text: '正在上传',
		background: 'rgba(0,0,0,.2)',
	});
	await useLoginApi()
		.idcardLogin(ruleFormPersonal.value)
		.then(async (res: any) => {
			loadingInstance.close();
			if (res.code === "SUCCESS") {
				let siginres = await onSignIn(res.data.token, res.data.type, res.data.subType, false)
        if (siginres) {
          if (route.query?.type === 'fromOauth2' || route.query?.type === 'fromApplication') {
            let domain_uid = route.query?.domain;
            let redirect_uri = route.query?.redirect_uri;
            let next_redirect_uri = route.query?.next_redirect_uri;
            let params = route.query?.params;
            let isLoginToken = Number(route.query?.isLoginToken) == 1;

            let res = await useoauth2Api()
                .oauth2(
                    { domainUID: domain_uid },
                    async () => {
                      ElMessage.error('授权登录失败');
                      Session.clear();
                      Local.clear();
                    },
                    isLoginToken
                )
                .then(async (res: any): Promise<boolean> => {
                  if (res.code === 'SUCCESS') {
                    let queryString = encodeSearchParams({
                      token: res.data.token,
                      subToken: res.data.subToken,
                      params: params,
                      redirect: next_redirect_uri,
                    });

                    if (route.query?.type === 'fromApplication') {
                      window.open(`${redirect_uri}?${queryString}`);
                      await router.replace('/home');
                    } else {
                      window.location.href = `${redirect_uri}?${queryString}`;
                    }
                    return true;
                  } else if (res.subCode === 'NOT_OPEN_WEBSITE') {
                    if (route.query?.type === 'fromApplication') {
                      window.open(router.resolve({
                        path: '/oauth2/open',
                        query: {
                          type: route.query?.type,
                          params: params,
                          next_redirect_uri: next_redirect_uri,
                          redirect_uri: redirect_uri,
                          domain: domain_uid,
                          isLoginToken: isLoginToken ? 1 : 0,
                        },
                      }).href);  // 新窗口打开
                    } else {
                      await router.push({
                        path: '/oauth2/open',
                        query: {
                          type: route.query?.type,
                          params: params,
                          next_redirect_uri: next_redirect_uri,
                          redirect_uri: redirect_uri,
                          domain: domain_uid,
                          isLoginToken: isLoginToken ? 1 : 0,
                        },
                      }); // 去登录页
                    }
                    return true;
                  } else {
                    return false;
                  }
                });
            if (!res) {
              ElMessage.error('授权登录失败');
              Session.clear();
              Local.clear();
            }
          } else if (route.query?.redirect) {
            let params = route.query?.params as string;
            let querys = {};
            if (params) {
              querys = JSON.parse(params);
            }

            await router.replace({
              path: <string>route.query?.redirect,
              query: querys,
            });
          } else {
            await router.replace('/home');
          }
        } else {
          ElMessage.error('登录失败');
        }
        NextLoading.done()
			} else {
				ElMessage.error('上传失败');
			}
		});
};
//  企业认证成功
const scanBusinessSuccess = async (faceToken: string) => {
	if (!faceToken) {
		ElMessage.warning('人脸识别失败');
		return;
	}
	if (!Session.get('phoneToken')) {
		ElMessage.warning('请先获取验证码');
		return;
	}
	ruleFormBusinese.value.phoneToken = Session.get('phoneToken');
	ruleFormBusinese.value.faceToken = faceToken;
	const loadingInstance = ElLoading.service({
		text: '正在上传',
		background: 'rgba(0,0,0,.2)',
	});
	await useLoginApi()
		.legalPersonIDCardLogin(ruleFormBusinese.value)
		.then(async (res: any) => {
			loadingInstance.close();
			if (res.code === "SUCCESS") {
				let signres = await onSignIn(res.data.token, res.data.type, res.data.subType);
        if (signres) {
          if (route.query?.type === 'fromOauth2' || route.query?.type === 'fromApplication') {
            let domain_uid = route.query?.domain;
            let redirect_uri = route.query?.redirect_uri;
            let next_redirect_uri = route.query?.next_redirect_uri;
            let params = route.query?.params;
            let isLoginToken = Number(route.query?.isLoginToken) == 1;

            let res = await useoauth2Api()
                .oauth2(
                    { domainUID: domain_uid },
                    async () => {
                      ElMessage.error('授权登录失败');
                      Session.clear();
                      Local.clear();
                    },
                    isLoginToken
                )
                .then(async (res: any): Promise<boolean> => {
                  if (res.code === 'SUCCESS') {
                    let queryString = encodeSearchParams({
                      token: res.data.token,
                      subToken: res.data.subToken,
                      params: params,
                      redirect: next_redirect_uri,
                    });

                    if (route.query?.type === 'fromApplication') {
                      window.open(`${redirect_uri}?${queryString}`);
                      await router.replace('/home');
                    } else {
                      window.location.href = `${redirect_uri}?${queryString}`;
                    }
                    return true;
                  } else if (res.subCode === 'NOT_OPEN_WEBSITE') {
                    if (route.query?.type === 'fromApplication') {
                      window.open(router.resolve({
                        path: '/oauth2/open',
                        query: {
                          type: route.query?.type,
                          params: params,
                          next_redirect_uri: next_redirect_uri,
                          redirect_uri: redirect_uri,
                          domain: domain_uid,
                          isLoginToken: isLoginToken ? 1 : 0,
                        },
                      }).href);  // 新窗口打开
                    } else {
                      await router.push({
                        path: '/oauth2/open',
                        query: {
                          type: route.query?.type,
                          params: params,
                          next_redirect_uri: next_redirect_uri,
                          redirect_uri: redirect_uri,
                          domain: domain_uid,
                          isLoginToken: isLoginToken ? 1 : 0,
                        },
                      }); // 去登录页
                    }
                    return true;
                  } else {
                    return false;
                  }
                });
            if (!res) {
              ElMessage.error('授权登录失败');
              Session.clear();
              Local.clear();
            }
          } else if (route.query?.redirect) {
            let params = route.query?.params as string;
            let querys = {};
            if (params) {
              querys = JSON.parse(params);
            }

            await router.replace({
              path: <string>route.query?.redirect,
              query: querys,
            });
          } else {
            await router.replace('/home');
          }
        } else {
          ElMessage.error("登录失败")
        }
        NextLoading.done()
			} else {
				ElMessage.error('上传失败');
			}
		});
};
// 页面加载时
onMounted(() => {
	NextLoading.done();
});
</script>

<style scoped lang="scss">
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

.login-action {
	display: flex;
	flex-direction: row;
	align-items: center;
	position: relative;
	z-index: 2002;

	.admin-agreement {
		color: var(--el-color-primary);
		margin-left: 4px;
		margin-right: 30px;
		text-decoration: underline;
		cursor: pointer;
	}
}

:deep(.el-tabs__active-bar) {
	height: 3px;
}

:deep(input::-webkit-input-placeholder) {
	color: #9ca5ba;
	font-size: 12px;
}

:deep(.el-input__inner) {
	background-color: transparent !important;
}

:deep(.el-tabs__item.is-active) {
	color: #121212;
	font-weight: bolder;
}

:deep(.el-step__icon-inner) {
	font-size: 16px !important;
}

.login-container {
	height: 100%;
	background: var(--el-color-white);

  .login-left {
    flex: 1;
    position: relative;
    display: flex;
    background-color: rgba(211, 239, 255, 1);
    align-items: center;
    justify-content: center;
    .content-box {
      padding: 15px;
      border: 1px solid var(--el-color-primary);
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
    .login-left-logo {
      display: flex;
      align-items: center;
      position: absolute;
      top: 50px;
      left: 80px;
      z-index: 1;
      animation: logoAnimation 0.3s ease;
      img {
        width: 52px;
        height: 52px;
      }
      .login-left-logo-text {
        display: flex;
        flex-direction: column;
        span {
          margin-left: 10px;
          font-size: 28px;
          color: #26a59a;
        }
        .login-left-logo-text-msg {
          font-size: 12px;
          color: #32a99e;
        }
      }
    }
    .login-left-img {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      width: 50%;
      height: 40%;
      img {
        width: 100%;
        height: 100%;
        animation: error-num 0.6s ease;
        object-fit: contain;
      }
    }
    .login-left-waves {
      position: absolute;
      top: 0;
      right: -100px;
    }
  }

	.login-right {
		width: 700px;

		.login-right-warp {
			border: 1px solid var(--el-color-primary-light-3);
			border-radius: 3px;
			width: 500px;
			min-height: 650px;
			position: relative;
			overflow: hidden;
			background-color: var(--el-color-white);
			box-sizing: unset;

			.login-right-warp-one,
			.login-right-warp-two {
				position: absolute;
				display: block;
				width: inherit;
				height: inherit;

				&::before,
				&::after {
					content: '';
					position: absolute;
					z-index: 1;
				}
			}

			.login-right-warp-one {
				&::before {
					filter: hue-rotate(0deg);
					top: 0px;
					left: 0;
					width: 100%;
					height: 3px;
					background: linear-gradient(90deg, transparent, var(--el-color-primary));
					animation: loginLeft 3s linear infinite;
				}

				&::after {
					filter: hue-rotate(60deg);
					top: -100%;
					right: 2px;
					width: 3px;
					height: 100%;
					background: linear-gradient(180deg, transparent, var(--el-color-primary));
					animation: loginTop 3s linear infinite;
					animation-delay: 0.7s;
				}
			}

			.login-right-warp-two {
				&::before {
					filter: hue-rotate(120deg);
					bottom: 2px;
					right: -100%;
					width: 100%;
					height: 3px;
					background: linear-gradient(270deg, transparent, var(--el-color-primary));
					animation: loginRight 3s linear infinite;
					animation-delay: 1.4s;
				}

				&::after {
					filter: hue-rotate(300deg);
					bottom: -100%;
					left: 0px;
					width: 3px;
					height: 100%;
					background: linear-gradient(360deg, transparent, var(--el-color-primary));
					animation: loginBottom 3s linear infinite;
					animation-delay: 2.1s;
				}
			}

			.login-right-warp-mian {
				display: flex;
				flex-direction: column;
				justify-content: center;
				height: 100%;

				.login-right-warp-main-title {
					height: 80px;
					line-height: 80px;
					font-size: 27px;
					text-align: center;
					letter-spacing: 3px;
					animation: logoAnimation 0.3s ease;
					animation-delay: 0.3s;
					color: var(--el-text-color-primary);
				}

				.login-right-warp-main-form {
					flex: 1;
					padding: 0 30px;

					.login-right-warp-main-form_btns {
						display: flex;
						justify-content: center;
            align-items: center;
						flex-direction: column;
						height: 350px;

						.personal,
						.enterprise {
							width: 200px;
							margin: 45px 0;
						}

						:deep(.el-button) {
							width: 100%;
							height: 60px;
							border-radius: 10px;
						}
					}

					.login-content-main-sacn {
						position: absolute;
						top: 0;
						right: 0;
						width: 50px;
						height: 50px;
						overflow: hidden;
						cursor: pointer;
						transition: all ease 0.3s;
						color: var(--el-color-primary);

						&-delta {
							position: absolute;
							width: 35px;
							height: 70px;
							z-index: 2;
							top: 2px;
							right: 21px;
							background: var(--el-color-white);
							transform: rotate(-45deg);
						}

						&:hover {
							opacity: 1;
							transition: all ease 0.3s;
							color: var(--el-color-primary) !important;
						}

						i {
							width: 47px;
							height: 50px;
							display: inline-block;
							font-size: 48px;
							position: absolute;
							right: 1px;
							top: 0px;
						}
					}

					.register {
						display: flex;
						justify-content: center;
						font-size: 14px;
						margin-top: 20px;
						cursor: pointer;
						position: relative;
						z-index: 2004;

						:nth-child(1) {
							color: #9ca5ba;
						}

						:nth-child(2) {
							color: var(--el-color-primary);
							margin-left: 6px;
							text-decoration: underline;
						}
					}
				}

				.page-foot {
					width: calc(100% - 30px);
					margin-left: 15px;
					display: flex;
					padding-bottom: 10px;
					text-align: center;
					color: #9ca5ba;
					display: flex;
					flex-direction: row;
					justify-content: center;
					cursor: pointer;
					position: relative;
					z-index: 2004;

					a {
						margin-right: 8px;
						color: #9ca5ba;
					}

					img {
						width: 20px;
						height: 20px;
					}
				}
			}
		}
	}
}
</style>
