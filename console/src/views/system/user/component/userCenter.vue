<template>
	<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="869px" @close="closeDialog" destroy-on-close>
		<el-card shadow="hover">
			<div class="main">
				<div class="user-info">
					<div class="user-info-left">
						<div class="main-title">账号信息</div>
						<div class="user-head">
							<div v-if="state.dialog.isShowDialog" style="width: 60px; height: 60px">
								<showImage :user="state.ruleForm.id" />
							</div>
							<div class="user-head-right">
								<div>登录账号</div>
								<div v-if="userInfo && userInfo.phone">
									{{ userInfo.nickname || userInfo.userName || userInfo.phone || '异常用户' }}
								</div>
							</div>
						</div>
					</div>
					<div class="user-info-right">
						<div class="item-action" @click.stop="updateUserName('userName')" style="width: 60px">用户名</div>
					</div>
				</div>
				<div class="info-list">
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="phone" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">手机号</div>
							<div class="item-text" v-if="userInfo && userInfo.phone">{{ userInfo.phone }}</div>
							<div class="item-text" v-else>您暂未绑定手机号</div>
						</div>
						<div class="item-action" v-if="userInfo && userInfo.phone" @click="updatePhone">修改</div>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">暂未开通</el-text>
						</div>
					</div>
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="email" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">邮箱</div>
							<div class="item-text" v-if="userInfo && userInfo.email">{{ userInfo.email }}</div>
							<div class="item-text" v-else>您暂未绑定邮箱</div>
						</div>
						<div class="item-action" @click="updateEmail">修改</div>
						<!--						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">暂未开通</el-text>
						</div>-->
					</div>
					<div class="info-item">
						<SvgIcon name="my-weixin" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">微信</div>
							<div class="item-text" v-if="userInfo && userInfo.wxNickName">您已绑定了微信({{ userInfo.wxNickName }})</div>
							<div class="item-text" v-else>您暂未绑定微信</div>
						</div>
						<el-text type="danger" v-if="userInfo && userInfo.wxNickName" @click="delWeixin">解绑</el-text>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">暂未开通</el-text>
						</div>
					</div>
					<!--          查看父级信息-->
					<div class="info-item">
						<SvgIcon name="my-people" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">父级信息</div>
							<div class="item-text" v-if="userInfo && userInfo.father !== 0">您已绑定了父级信息</div>
							<div class="item-text" v-else>您暂未绑定父级信息</div>
						</div>
						<el-text type="primary" v-if="userInfo && userInfo.father !== 0" @click="fatherInfo">查看</el-text>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无父级信息</el-text>
						</div>
					</div>
					<!--          查看邀请人信息-->
					<div class="info-item">
						<SvgIcon name="my-people" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">邀请人信息</div>
							<div class="item-text" v-if="userInfo && userInfo.invite !== 0">您已绑定了邀请人信息</div>
							<div class="item-text" v-else>您暂未绑定邀请人信息</div>
						</div>
						<el-text type="primary" v-if="userInfo && userInfo.invite !== 0" @click="inviteInfo">查看</el-text>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无邀请人</el-text>
						</div>
					</div>
				</div>
			</div>
		</el-card>
		<el-card shadow="hover" class="card">
			<div class="main">
				<div class="main-title">安全服务</div>
				<div class="info-list">
					<!-- 设置用户状态-->
					<div class="info-item">
						<SvgIcon name="my-people" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">用户状态</div>
							<div class="item-text" v-if="userInfo && userInfo.status">
								<el-text class="mx-1" :type="userColorEnumTypes[userInfo.status]" size="mini">{{ UserStatusEnumTypes[userInfo.status] }}</el-text>
							</div>
						</div>
						<div class="item-action" v-if="userInfo && userInfo.status" @click="updateStatus">修改</div>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无权限</el-text>
						</div>
					</div>
					<!--          登录控制-->
					<div class="info-item">
						<SvgIcon name="my-people" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">登录控制</div>
							<div class="item-text">
								<el-text class="mx-1" size="mini">设置用户登录方式</el-text>
							</div>
						</div>
						<div class="item-action" @click="loginControl">设置</div>
					</div>
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="totp" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">绑定二次验证方式(2FA)</div>
							<div class="item-text" v-if="userInfo && userInfo.has2FA">您已绑定了绑定二次验证方式(2FA)</div>
							<div class="item-text" v-else>您暂未绑定二次验证方式(2FA)</div>
						</div>
						<div class="item-action">
							<el-text v-if="userInfo && userInfo.has2FA" type="danger" @click="delTotp">解绑</el-text>
							<el-text v-else type="info" size="mini">暂未开通</el-text>
						</div>
					</div>
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="login" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">登录密码</div>
							<div class="item-text">安全性高的密码可以使账号更安全。建议您定期更换密码，设置一个包含字母，符号或数字中至少两项且超过6位的密码</div>
						</div>
						<div class="item-action" v-if="authUser.updatePassword" @click="clickPassword">修改</div>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无权限</el-text>
						</div>
					</div>
				</div>
			</div>
		</el-card>
		<!-- 更新用户名 -->
		<user-name ref="userNameRef" :user-name-form="userNameForm" @refresh="getUserInfo" />
		<!-- 更新手机号绑定 -->
		<Mobile ref="mobileRef" :phone-form="phoneForm" @refresh="getUserInfo" />
		<!-- 邮箱绑定 -->
		<Email ref="emailRef" :email-form="emailForm" @refresh="getUserInfo" />
		<!-- 令牌有效时间 -->
		<ExpirationToken ref="expirationTokenRef" :userNameForm="tokenForm" @refresh="getUserInfo" />
		<!-- 修改密码 -->
		<Password ref="passwordRef"></Password>
		<!-- 手机验证注销弹窗 -->
		<LogOffMobile ref="logOffMobileRef" @reset="reset" />
		<!-- 更新用户状态 -->
		<UpdateStatus ref="updateStatusRef" @refresh="getUserInfo" />
		<!-- 查看父级信息 -->
		<FatherInfo v-if="userInfo && userInfo !== 0" ref="fatherInfoRef" />
		<!-- 邀请人信息 -->
		<Inviter v-if="userInfo && userInfo !== 0" ref="inviterRef" />
		<!-- 登录控制 -->
		<updateLoginCtr ref="updateLoginCtrRef" />
	</el-dialog>
</template>

<script setup lang="ts">
import { ref, defineAsyncComponent, reactive } from 'vue';
import totp from '/@/assets/userCenter/totp.png';
import login from '/@/assets/userCenter/login.png';
import { Local, Session } from '/@/utils/storage';
import LogOffMobile from '/@/views/userCenter/countInfo/component/logOffMobile.vue';
import Mobile from '/@/views/system/user/component/user/mobile.vue';
import Email from '/@/views/system/user/component/user/email.vue';
import Password from '/@/views/system/user/component/user/password.vue';
import UpdateStatus from '/@/views/system/user/component/user/updateStatus.vue';
import showImage from '/@/components/showImage/index.vue';
import email from '/@/assets/userCenter/email.png';
import { ElMessage, ElMessageBox } from 'element-plus';
import { checkWxRobotTypes } from '/@/api/base/types';
import { updateUserNameTypes } from '/@/api/register/types';
import phone from '/@/assets/userCenter/phone.png';
import { totpConfigTypes } from '/@/api/Totp/types';
import useSubAuth from '/@/hooks/useSubAuth';
import UserName from '/@/views/system/user/component/user/userName.vue';
import { userApi } from '/@/api/system/user';
import { userEmailTypes, userPhoneTypes } from '/@/views/system/user/types';
import { userColorTypes, UserStatusEnum, userTypes } from '/@/data/enum';
import FatherInfo from '/@/views/system/user/component/fatherInfo.vue';
import Inviter from '/@/views/system/user/component/inviter.vue';
import UpdateLoginCtr from '/@/views/userCenter/countInfo/component/updateLoginCtr.vue';

const ExpirationToken = defineAsyncComponent(() => import('/@/views/userCenter/countInfo/component/expirationToken.vue'));
// 引入组件
const logOffMobileRef = ref<any>(false);

const userColorEnumTypes: userTypes = userColorTypes;
const UserStatusEnumTypes: userTypes = UserStatusEnum;
const mobileRef = ref();
const emailRef = ref();
const userInfo = ref<any>();
const passwordRef = ref();
const updateStatusRef = ref();
const userNameRef = ref();
const authUser = useSubAuth();
const expirationTokenRef = ref();
const fatherInfoRef = ref();
const inviterRef = ref();
const updateLoginCtrRef = ref();
const emit = defineEmits(['refresh']);
const state = reactive({
	dialog: {
		isShowDialog: false,
		title: '',
		submitTxt: '',
	},
	ruleForm: {} as any,
});
// 单点登录
const signin = ref<boolean>(false);
// 手机号
const phoneForm = ref<userPhoneTypes>({
	phone: '',
	uid: '',
});
const tokenForm = ref({
	tokenExpiration: '',
});
// 用户名
const userNameForm = ref<updateUserNameTypes>({
	userName: '',
	nickname: '',
});
// 邮箱
const emailForm = ref<userEmailTypes>({
	email: '',
	uid: '',
});
// 2fa基本数据
const totpForm = ref<totpConfigTypes>({
	phone: '',
	secret: '',
});
// 微信机器人
const wxRobotForm = ref<checkWxRobotTypes>({
	webhook: '',
	isDelete: false,
});

const updateLoginCtrFrom = ref({
	allowPhone: false,
	allowEmail: false,
	allowPassword: false,
	allowWeChat: false,
	allowSecondFA: false,
});
//跟新手机号
const updatePhone = () => {
	mobileRef.value.openDialog();
};
// 登录控制
const loginControl = () => {
	updateLoginCtrRef.value.openDialog({ uid: userInfo.value.uid, ...updateLoginCtrFrom.value }, 'admin');
};
//更新用户名
const updateUserName = (type: string) => {
	userNameRef.value.openDialog(type, { id: userInfo.value.uid });
};
//更新邮箱
const updateEmail = () => {
	emailRef.value.openDialog();
};
//更新用户状态
const updateStatus = () => {
	updateStatusRef.value.openDialog({ uid: userInfo.value.uid, status: userInfo.value.status });
};
//解绑微信
const delWeixin = () => {
	ElMessageBox.confirm('是否解绑微信?', '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			userApi()
				.unbindWechat({ uid: userInfo.value.uid })
				.then((res: any) => {
					if (res.code === 'SUCCESS') {
						ElMessage.success('解绑微信成功');
            if (sessionStorage.getItem("fuwuhao-login") === "true") {  // 服务号登录如果解绑了微信必须退出登录
              Session.clear();
              Local.clear();
              window.location.href = '/login'; // 回到登录页
            } else {
              getUserInfo();
            }
					}
				});
		})
		.catch(() => {});
};
const openDialog = async (row?: any) => {
	state.ruleForm = JSON.parse(JSON.stringify(row));
	state.dialog.title = '修改用户';
	state.dialog.submitTxt = '修 改';
	await getUserInfo();
	state.dialog.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
	state.dialog.isShowDialog = false;
};
//删除2fa
const delTotp = () => {
  if(authUser.is2FA){
    userApi()
        .unbind2FA({ uid: userInfo.value.uid })
        .then((res: any) => {
          if (res.code === 'SUCCESS') {
            ElMessage({
              type: 'success',
              message: '解绑成功',
            });
            getUserInfo();
          }
        });
  }else{
    ElMessageBox.confirm('请输入双因素验证器验证码', '是否解绑2FA验证?', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
        .then(() => {
          userApi()
              .unbind2FA({ uid: userInfo.value.uid })
              .then((res: any) => {
                if (res.code === 'SUCCESS') {
                  ElMessage({
                    type: 'success',
                    message: '解绑成功',
                  });
                  getUserInfo();
                }
              });
        })
        .catch(() => {});
  }

};
//密码
const clickPassword = () => {
	passwordRef.value.openDialog({ id: userInfo.value.uid });
};
//重置
const reset = () => {
	Session.clear();
	Local.clear();
	// 使用 reload 时，不需要调用 resetRoute() 重置路由
	window.location.reload();
};
const fatherInfo = () => {
	fatherInfoRef.value.openDialog({ id: userInfo.value.father });
};
const inviteInfo = () => {
	inviterRef.value.openDialog({ id: userInfo.value.invite });
};
const getUserInfo = async () => {
	await userApi()
		.getUserData({ uid: state.ruleForm.id })
		.then((res: any) => {
			if (res.code === 'SUCCESS') {
				userInfo.value = res.data;
				phoneForm.value.phone = userInfo.value.phone;
				phoneForm.value.uid = userInfo.value.uid;
				totpForm.value.phone = userInfo.value.phone;
				emailForm.value.email = userInfo.value.email;
				emailForm.value.uid = userInfo.value.uid;
				signin.value = userInfo.value.signin;
				if (userInfo.value?.wxWebHook) {
					wxRobotForm.value.webhook = userInfo.value.wxWebHook;
				}
				updateLoginCtrFrom.value.allowPhone = userInfo.value.allowPhone;
				updateLoginCtrFrom.value.allowEmail = userInfo.value.allowEmail;
				updateLoginCtrFrom.value.allowPassword = userInfo.value.allowPassword;
				updateLoginCtrFrom.value.allowWeChat = userInfo.value.allowWeChat;
				updateLoginCtrFrom.value.allowSecondFA = userInfo.value.allowSecondFA;
				if (userInfo.value?.phone) {
					if (userInfo.value?.userName) {
						userNameForm.value.userName = userInfo.value.userName;
					}
					userNameForm.value.nickname = userInfo.value.nickname;
				}
				emit('refresh');
			}
		});
};
// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style lang="scss" scoped>
.el-input__wrapper {
	margin-bottom: 10px;
}

.main-title {
	font-size: 16px;
	font-weight: 600;
	color: #000000;
	line-height: 24px;
}

:deep(.el-card) {
	overflow-y: auto;
}
.container {
	width: 100%;
	:deep(.el-card__body) {
		overflow: hidden;
	}
}

.card {
	margin-top: 16px;
}

.notice-bar {
	display: flex;
	flex-direction: row;

	.notice-date {
		width: 42px;
		height: 54px;
		background: #f1f2f5;
		border-radius: 4px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;

		.day {
			color: #61687c;
			height: 24px;
			font-weight: bold;
			font-size: 16px;
			line-height: 24px;
		}

		.mouth {
			height: 20px;
			font-size: 12px;
			line-height: 20px;
			color: #61687c;
		}
	}

	.notice-main {
		flex: 1;
		height: 54px;
		margin: 0px 15px;

		.content-box-item {
			height: 54px;
			display: flex;
			flex-direction: column;
			justify-content: space-between;
			cursor: pointer;
		}

		.content-box-title {
			height: 22px;
			font-size: 14px;
			font-weight: 600;
			color: #414960;
			line-height: 22px;
		}
	}

	.notice-arrow {
		width: 24px;
		height: 54px;
		margin-left: 8px;
		display: flex;
		align-items: center;
		justify-content: center;
	}
}

.main {
	.user-info {
		display: flex;
		justify-content: space-between;
		align-items: center;

		.user-info-right {
			display: flex;
			flex-direction: row;
		}

		.user-head {
			padding: 16px 0px 20px 0px;
			display: flex;
			border-bottom: 1px solid #f1f2f5;

			.user-head-right {
				margin-left: 16px;

				:nth-child(1) {
					height: 22px;
					font-size: 14px;
					font-weight: 400;
					color: #61687c;
					line-height: 22px;
				}

				:nth-child(2) {
					height: 22px;
					font-size: 14px;
					font-weight: 400;
					color: #414960;
					line-height: 22px;
					margin-top: 8px;
				}

				.roleName {
					color: var(--el-text-color-secondary);
				}
			}
		}

		.item-action {
			color: var(--el-color-primary);
			line-height: 40px;
			padding-top: 10px;
			padding-left: 8px;
			cursor: pointer;
		}
	}

	.info-list {
		.info-item {
			padding: 16px 0px 20px 0px;
			display: flex;
			flex-direction: row;
			border-bottom: 1px solid #f1f2f5;
			max-height: 200px;

			.item-main {
				flex: 1;
				margin-left: 8px;

				.item-title {
					height: 22px;
					font-size: 14px;
					font-weight: 400;
					color: #61687c;
					line-height: 22px;
				}

				.item-text {
					font-size: 14px;
					font-weight: 400;
					color: #414960;
					min-height: 22px;
					margin-top: 8px;
					display: flex;
					align-items: center;

					span {
						display: flex;
						align-items: center;
					}

					.binding {
						color: var(--el-color-primary);
						cursor: pointer;
					}
				}
			}

			.item-action {
				color: var(--el-color-primary);
				line-height: 40px;
				padding-top: 10px;
				padding-left: 8px;
				cursor: pointer;
			}
		}
	}
}
</style>
