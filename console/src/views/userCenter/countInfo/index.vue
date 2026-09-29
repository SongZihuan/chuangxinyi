<template>
	<div class="container layout-padding">
		<el-card shadow="hover" class="mb15">
			<div class="main-title mb10">公告栏</div>
			<div class="notice-bar">
				<div class="notice-date">
					<div class="day">{{ dayjs().format('DD') }}</div>
					<div class="mouth">{{ dayjs().format('MM') }}月</div>
				</div>
				<div class="notice-main">
					<el-carousel
						height="54px"
						arrow="never"
						v-if="announcementList.length >= 1"
						interval="6000"
						ref="carouselRef"
						:pause-on-hover="false"
						indicator-position="none"
					>
						<el-carousel-item v-for="(item, index) in announcementList" :key="index" @click="showContent(item)">
							<div class="content-box-item">
								<div class="content-box-title">
									<div>{{ item.title }}</div>
								</div>
								<div class="content-box-time">
									<el-tag>{{ dayjs.unix(item.startAt).format('YYYY-MM-DD') }}至{{ dayjs.unix(item.stopAt).format('YYYY-MM-DD') }} </el-tag>
								</div>
							</div>
						</el-carousel-item>
					</el-carousel>
					<div class="content-box-item" v-else>
						<div class="content-box-title">
							<div>暂无公告</div>
						</div>
					</div>
				</div>
				<div class="notice-arrow" v-if="announcementList.length > 1">
					<el-image style="width: 24px; height: 24px" :src="arrowActive" fit="cover" lazy v-if="announcementList.length > 1" @click="next" />
					<el-image style="width: 24px; height: 24px" :src="arrow" fit="cover" lazy v-else />
				</div>
			</div>
		</el-card>
		<Application class="mb20" />
		<!-- 财务信息 -->
		<FinanceInfo :finance-info="financeInfo" class="mb20" />
		<el-card shadow="hover">
			<div class="main">
				<div class="user-info">
					<div class="user-info-left">
						<div class="main-title">账号信息</div>
						<div class="user-head" v-loading="avatarLoading">
							<el-tooltip v-if="authUser.updateHeader" class="box-item" effect="dark" content="更换头像" placement="top-start">
								<div style="width: 48px; height: 48px">
									<AvatarUpload
										ref="avatarUploadRef"
										v-if="userInfo"
										@imgSuccess="imgSuccess"
										:userId="userInfo.user.id"
										:key="avatarUploadKey"
									></AvatarUpload>
								</div>
							</el-tooltip>
							<div v-else style="width: 48px; height: 48px">
								<AvatarUpload
									ref="avatarUploadRef"
									v-if="userInfo"
									:disabledType="!authUser.updateHeader"
									@imgSuccess="imgSuccess"
									:userId="userInfo.user.id"
									:key="avatarUploadKey"
								></AvatarUpload>
							</div>
							<div class="user-head-right">
								<div>登录账号</div>
								<div v-if="userInfo && userInfo.user.phone">
									{{ userInfo.user.nickname || userInfo.user.userName || userInfo.user.phone || '异常用户'
									}}<span v-if="userInfo.user.nickname">({{ userInfo.user.userName || userInfo.user.phone }})</span>-<span class="roleName">{{
										userInfo.role.name
									}}</span>
								</div>
							</div>
						</div>
					</div>
					<div class="user-info-right" style="width: 120px; margin-right: -30px">
						<div class="item-action" v-if="authUser.updateUserName" @click.stop="updateUserName('userName')" style="width: 60px">用户名</div>
						<div class="item-action mr10" v-if="authUser.updateUserName" @click.stop="updateUserName('nickName')" style="width: 60px">昵称</div>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无权限</el-text>
						</div>
					</div>
				</div>
				<div class="info-list">
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="phone" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">绑定手机</div>
							<div class="item-text" v-if="userInfo && userInfo.user.phone">
								您已绑定了手机{{ userInfo.user.phone }} ，用于安全和二次校验能正常收取校验码短信等
							</div>
							<div class="item-text" v-else>您暂未绑定手机号</div>
						</div>
						<div class="item-action" v-if="authUser.updatePhone" @click="updatePhone">修改</div>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无权限</el-text>
						</div>
					</div>
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="email" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">绑定邮箱</div>
							<div class="item-text" v-if="userInfo && userInfo.data.hasEmail">
								您已绑定了邮箱{{ userInfo.user.email }} ，用于安全和二次校验能正常收取校验码等
							</div>
							<div class="item-text" v-else>您暂未绑定邮箱</div>
						</div>
						<template v-if="authUser.updateEmail">
							<div class="item-action" @click="updateEmail">修改</div>
						</template>
						<template v-else>
							<div class="item-action">
								<el-text type="info">抱歉，暂无权限</el-text>
							</div>
						</template>
					</div>
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="weixin" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">绑定微信</div>
							<div class="item-text" v-if="wxStatus === 4">
                您暂未绑定微信
              </div>
							<div class="item-text" v-else-if="wxStatus === 1">
								<el-image
									:style="{ minWidth: `30px`, height: `30px`, borderRadius: '50%' }"
									:src="userInfo.user.wechatHeader"
									:zoom-rate="1.2"
									:preview-src-list="[userInfo.user.wechatHeader]"
									preview-teleported
									v-if="userInfo"
									fit="cover"
									close-on-press-escape
									class="mr10"
								/>
								您已绑定了微信({{ userInfo.user.wechatNickName }}) 且关注了服务号
							</div>
							<div class="item-text" v-else-if="wxStatus === 2">
								<el-image
									:style="{ minWidth: `30px`, height: `30px`, borderRadius: '50%' }"
									:src="userInfo.user.wechatHeader"
									:zoom-rate="1.2"
									:preview-src-list="[userInfo.user.wechatHeader]"
									preview-teleported
									v-if="userInfo"
									fit="cover"
									close-on-press-escape
									class="mr10"
								/>
								您已绑定了微信({{ userInfo.user.wechatNickName }}) 但是未关注服务号(<span class="binding" @click="openWxService">点击绑定</span>)
							</div>
              <div class="item-text" v-else-if="wxStatus === 3">
                <el-image
                    :style="{ minWidth: `30px`, height: `30px`, borderRadius: '50%' }"
                    :src="userInfo.user.wechatHeader"
                    :zoom-rate="1.2"
                    :preview-src-list="[userInfo.user.wechatHeader]"
                    preview-teleported
                    v-if="userInfo"
                    fit="cover"
                    close-on-press-escape
                    class="mr10"
                />
                您未绑定了微信({{ userInfo.user.wechatNickName }}) 但是关注了服务号
              </div>
						</div>
						<div class="item-action" v-if="authUser.updateWechat">
							<el-text type="danger" v-if="userInfo && wxStatus !== 4" @click="delWeixin">解绑</el-text>
							<el-text type="primary" v-else @click="updateWeixin">绑定</el-text>
						</div>
						<div class="item-action" v-else>
							<el-text type="info">抱歉，暂无权限</el-text>
						</div>
					</div>
					<!--  微信机器人        -->
					<div class="info-item">
						<SvgIcon name="my-robot" color="#61687c" :size="22"></SvgIcon>
						<div class="item-main">
							<div class="item-title">绑定微信机器人</div>
							<div class="item-text" v-if="userInfo && userInfo.data.hasWxrobot">您已绑定微信机器人</div>
							<div class="item-text" v-else>您暂未绑定微信机器人</div>
						</div>
						<div class="item-action" v-if="authUser.updateWxrobot">
							<el-text type="danger" v-if="userInfo && userInfo.data.hasWxrobot" @click="deleteWxRobot">解绑</el-text>
							<el-text type="primary" v-else @click="updateWxRobot">绑定</el-text>
						</div>
						<div class="item-action" v-else>
							<el-text type="info">抱歉，暂无权限</el-text>
						</div>
					</div>
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="user" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">个人认证</div>
							<div class="item-text" v-if="userInfo && userInfo.data.hasVerified">
								您已上传身份信息 {{ userInfo && userInfo.data.hasUserOriginal ? '已上传身份证' : '您暂未上传身份证' }}
								{{ userInfo && userInfo.data.hasUserFaceCheck ? '已扫脸' : '您暂未扫脸' }}
							</div>
							<div class="item-text" v-else>您暂未进行实名认证</div>
						</div>
						<div class="item-action" v-if="authUser.updateRealName" @click="updateUser">修改</div>
						<div class="item-action" v-else>
							<el-text type="info">抱歉，暂无权限</el-text>
						</div>
					</div>
					<div class="info-item">
						<SvgIcon name="my-peoples" color="#61687c" :size="22"></SvgIcon>
						<div class="item-main">
							<div class="item-title">企业认证</div>
							<div class="item-text" v-if="userInfo && userInfo.data.isCompany">
								您已上传企业信息 {{ userInfo && userInfo.data.hasCompanyOriginal ? '已上传身份证' : '您暂未企业营业执照 法人身份证照' }}
								{{ userInfo && userInfo.data.hasLegalPersonFaceCheck ? '已扫脸' : '您暂未扫脸' }}
							</div>
							<div class="item-text" v-else>您暂未进行实名认证</div>
						</div>
						<template v-if="authUser.getinfo">
							<div class="item-action" @click="updateCompany">修改</div>
						</template>
						<template v-else>
							<div class="item-action">
								<el-text type="info">抱歉，暂无权限</el-text>
							</div>
						</template>
					</div>
					<div class="info-item">
						<SvgIcon name="my-tickets" color="#61687c" :size="22"></SvgIcon>
						<div class="item-main">
							<div class="item-title">发票抬头</div>
							<div class="item-text">
								{{ userInfo && userInfo.title && userInfo.title.taxID && userInfo.title.name ? '已上传发票抬头' : '您暂未上传发票抬头' }}
							</div>
						</div>
						<div class="item-action" v-if="authUser.updateTitle" @click="updateInvoice">
							{{ userInfo && userInfo.title ? '修改' : '填写' }}
						</div>
						<div class="item-action" v-else>
							<el-text type="info">抱歉，暂无权限</el-text>
						</div>
					</div>

					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="resgister" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">注册时间</div>
							<div class="item-text" v-if="userInfo">{{ dayjs.unix(userInfo.user.createAt).format('YYYY-MM-DD HH:mm:ss') }}</div>
						</div>
						<div class="item-action"></div>
					</div>
				</div>
			</div>
		</el-card>

		<!-- 父账号信息 -->
		<FatherInfo />
		<!-- 地址信息 -->
		<UserInfo />
		<!-- 邀请人信息 -->
		<Inviter />

		<el-card shadow="hover" class="card">
			<div class="main">
				<div class="main-title">邀请注册分销</div>
				<div class="info-list">
					<div class="info-item">
						邀请他人在注册的时候输入你的手机号（{{ userInfo && userInfo.user && userInfo.user.phone }}），即可在对方消费时获得返现。
					</div>
					<div class="info-item">
						<SvgIcon name="my-protect" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">分销平台收入</div>
							<div class="item-text">邀请用户在本平台的商品消费中，根据商品规定，返现一部分给您。</div>
							<div class="item-text">分销收入 = 邀请人消费金额 x 分销比例（大约为3%，参见商品详情）</div>
						</div>
					</div>
					<div class="info-item">
						<SvgIcon name="my-protect" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">分销管理功能</div>
							<div class="item-text">实名认证后，分销收益中可提现的部分可通过第三方渠道或人工方式提现。</div>
						</div>
					</div>
				</div>
			</div>
		</el-card>

		<el-card shadow="hover" class="card">
			<div class="main">
				<div class="main-title">安全服务</div>
				<div class="info-list">
					<!--          开启单点登录-->
					<div class="info-item">
						<SvgIcon name="my-safe" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">开启单点登录</div>
							<div class="item-text">开启后，您的账号将无法在多个设备上同时登录</div>
						</div>
						<div class="item-action" v-if="authUser.updateSignOne">
							<el-switch :model-value="signin" inline-prompt active-text="开启" inactive-text="关闭" @change="handelUpdate" />
						</div>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无权限</el-text>
						</div>
					</div>
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="totp" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">绑定二次验证方式(2FA)</div>
							<div class="item-text" v-if="userInfo && userInfo.data.has2FA">您已绑定了绑定二次验证方式(2FA)</div>
							<div class="item-text" v-else>您暂未绑定二次验证方式(2FA)</div>
						</div>
						<div class="item-action">
							<template v-if="userInfo && userInfo.data.has2FA">
								<template v-if="authUser.delete2FA">
									<el-text type="danger" @click="delTotp">解绑</el-text>
								</template>
								<template v-else>
									<el-text type="info">抱歉，暂无权限</el-text>
								</template>
							</template>
							<template v-else>
								<template v-if="authUser.bind2FA">
									<el-text type="primary" @click="bindTotp">绑定</el-text>
								</template>
								<template v-else>
									<el-text type="info">抱歉，暂无权限</el-text>
								</template>
							</template>
						</div>
					</div>
					<!--          登录控制-->
					<div class="info-item">
						<SvgIcon name="my-protect" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<div class="item-main">
							<div class="item-title">登录控制</div>
							<div class="item-text">开启后，您的账号将无法在多个设备上同时登录</div>
						</div>
						<div class="item-action" v-if="authUser.UpdateLoginController" @click="onUpdateController">修改</div>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无权限</el-text>
						</div>
					</div>
					<div class="info-item">
						<SvgIcon name="my-protect" :size="22" class="mr5" color="#61687c"></SvgIcon>
						<config theme="outline" size="24" fill="#333" />
						<div class="item-main">
							<div class="item-title">登录令牌有效时间</div>
							<div class="item-text">登录令牌最短有效时间为15分钟, 当前有效时长为{{ tokenForm.tokenExpiration }}秒</div>
						</div>
						<div class="item-action" v-if="authUser.updatePassword" @click="onToken">修改</div>
						<div class="item-action" v-else>
							<el-text class="mx-1" type="info" size="mini">抱歉，暂无权限</el-text>
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
					<div class="info-item">
						<el-image style="width: 24px; height: 24px" :src="logoff" fit="cover" lazy />
						<div class="item-main">
							<div class="item-title">账号</div>
							<div class="item-text">接收验证码，完成账号注销</div>
						</div>
						<div class="item-action">
							<template v-if="authUser.deleteUser">
								<el-text class="mx-1" type="danger" size="mini" @click="handleLogOff">注销</el-text>
							</template>
							<el-text class="mx-1" v-else type="info" size="mini">抱歉，暂无权限</el-text>
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
		<!-- 微信绑定 -->
		<Weixin ref="weixinRef" @refresh="getUserInfo" />
		<!-- 微信机器人绑定 -->
		<WxRobot ref="wxRobotRef" :wx-robot-form="wxRobotForm" @refresh="getUserInfo" />
		<!-- 2fa绑定 -->
		<TotpCom ref="totpComRef" :totp-form="totpForm" @refresh="getUserInfo"></TotpCom>
		<!-- 个人信息修改 -->
		<UploadUserInfo
			ref="userRef"
			:hasUserInfo="userInfo && userInfo.data.hasVerified"
			:hasUserOriginal="userInfo && userInfo.data.hasUserOriginal"
			:hasFace="userInfo && userInfo.data.hasUserFaceCheck"
			:user-form="userForm"
			@refresh="getUserInfo"
		></UploadUserInfo>
		<!-- 企业信息修改 -->
		<UploadCompany
			ref="companyRef"
			:hasUserInfo="userInfo && userInfo.data.isCompany"
			:hasCompanyOriginal="userInfo && userInfo.data.hasCompanyOriginal"
			:hasFace="userInfo && userInfo.data.hasLegalPersonFaceCheck"
			:company-form="companyForm"
			@refresh="getUserInfo"
		></UploadCompany>
		<!-- 令牌有效时间 -->
		<ExpirationToken ref="expirationTokenRef" :userNameForm="tokenForm" @refresh="getUserInfo" />
		<!-- 修改密码 -->
		<Password ref="passwordRef"></Password>
		<!-- 手机验证注销弹窗 -->
		<LogOffMobile ref="logOffMobileRef" @reset="reset" />
		<!-- 服务号绑定弹窗   -->
		<weixin-service ref="weixinServiceRef" :wx-status="wxStatus" @refresh="getUserInfo"></weixin-service>
		<!-- 发票抬头 -->
		<invoice ref="invoiceRef" :invoice-form="invoiceForm" @refresh="getUserInfo"></invoice>
		<!-- 公告 -->
		<ShowHtml ref="showHtmlRef" :title="currentInfo.title" :contentHtml="currentInfo.content" :subtitle="currentInfo.subtitle" />
		<!--    登录控制-->
		<updateLoginCtr ref="updateLoginCtrRef" @refresh="getUserInfo"></updateLoginCtr>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent, watch } from 'vue';
import UserInfo from '/@/views/userCenter/userInfo/index.vue';
import dayjs from 'dayjs';
import phone from '/@/assets/userCenter/phone.png';
import email from '/@/assets/userCenter/email.png';
import totp from '/@/assets/userCenter/totp.png';
import resgister from '/@/assets/userCenter/resgister.png';
import user from '/@/assets/userCenter/user.png';
import logoff from '/@/assets/userCenter/logoff.png';
import login from '/@/assets/userCenter/login.png';
import weixin from '/@/assets/userCenter/weixin.png';
import arrow from '/@/assets/userCenter/arrow.png';
import arrowActive from '/@/assets/userCenter/arrowActive.png';
import { Local, Session } from '/@/utils/storage';
import { useUserApi } from '/@/api/user/user';
import ShowHtml from '/@/components/ShowHtml/index.vue';
import LogOffMobile from '/@/views/userCenter/countInfo/component/logOffMobile.vue';
import Mobile from '/@/views/userCenter/countInfo/component/mobile.vue';
import Email from '/@/views/userCenter/countInfo/component/email.vue';
import TotpCom from '/@/views/userCenter/countInfo/component/totpCom.vue';
import UploadUserInfo from '/@/views/userCenter/countInfo/component/uploadUserInfo.vue';
import UploadCompany from '/@/views/userCenter/countInfo/component/uploadCompany.vue';
import Password from '/@/views/userCenter/countInfo/component/password.vue';
import Weixin from '/@/views/userCenter/countInfo/component/weixin.vue';
import invoice from '/@/views/userCenter/countInfo/component/invoice.vue';
import FatherInfo from '/@/views/userCenter/countInfo/component/fatherInfo.vue';
import Inviter from '/@/views/invite/inviter.vue';

import { useGoogleAuthApi } from '/@/api/Totp/index';
import { ElMessage, ElMessageBox } from 'element-plus';
import { checkEmailTypes, checkPhoneTypes, checkWxRobotTypes } from '/@/api/base/types';
import { updateUserNameTypes, uploadCompanyInfoJsonTypes, uploadUserInfoJsontTypes } from '/@/api/register/types';
import WeixinService from '/@/views/userCenter/countInfo/component/weixinService.vue';
import { totpConfigTypes } from '/@/api/Totp/types';
import { useWeixinApi } from '/@/api/weixin/index';
import { invoiceTitleTypes } from '/@/api/invoice/title/types';
import useSubAuth from '/@/hooks/useSubAuth';
import useOnlineUser from '/@/api/onlineUser';
import UserName from '/@/views/userCenter/countInfo/component/userName.vue';
import WxRobot from '/@/views/userCenter/countInfo/component/wxRobot.vue';
import { useRegisterApi } from '/@/api/register';
import { useUserInfo } from '/@/stores/userInfo';
import { useLoginSignIn } from '/@/hooks/useLoginSignIn';
import FinanceInfo from '/@/views/userCenter/countInfo/component/financeInfo.vue';
import { useSocketListInfo } from '/@/stores/socketListInfo';
import { storeToRefs } from 'pinia';
import { checkPhoneRes } from '/@/views/userCenter/countInfo/types';

const AvatarUpload = defineAsyncComponent(() => import('/@/components/avatarUpload/index.vue'));
const ExpirationToken = defineAsyncComponent(() => import('/@/views/userCenter/countInfo/component/expirationToken.vue'));
const Application = defineAsyncComponent(() => import('/@/views/userCenter/countInfo/component/application.vue'));
const updateLoginCtr = defineAsyncComponent(() => import('/@/views/userCenter/countInfo/component/updateLoginCtr.vue'));

// 引入组件
const carouselRef = ref();
const logOffMobileRef = ref<any>(false);
const stores = useUserInfo();
const useWeixinApiCollect = useWeixinApi();
const useGoogleAuthApiCollect = useGoogleAuthApi();
const useUserApiCollect = useUserApi();
const showHtmlRef = ref();
const mobileRef = ref();
const emailRef = ref();
const userInfo = ref();
const totpComRef = ref();
const userRef = ref();
const companyRef = ref();
const passwordRef = ref();
const weixinRef = ref();
const weixinServiceRef = ref();
const invoiceRef = ref();
const userNameRef = ref();
const authUser = useSubAuth();
const wxRobotRef = ref();
const expirationTokenRef = ref();
const updateLoginCtrRef = ref();
const avatarUploadKey = ref(1); //上传头像组件key
const avatarLoading = ref(false); //上传头像loading
const { resetSignIn } = useLoginSignIn();
const userInfoStore = useUserInfo();
const avatarUploadRef = ref();
watch(
	() => userInfoStore.userInfos,
	(newVal) => {
		setUserInfo(newVal);
	},
	{ deep: true }
);
// 发票抬头
const invoiceForm = ref<invoiceTitleTypes>({
	name: '',
	taxID: '',
	bankID: '',
	bank: '',
});
// 登录控制
let updateLoginCtrForm = ref<checkPhoneRes>({
	allowPhone: false,
	allowEmail: false,
	allowPassword: false,
	allowWeChat: false,
	allowSecondFA: false,
});
// 单点登录
const signin = ref<boolean>(false);
// 手机号
const phoneForm = ref<checkPhoneTypes>({
	phone: '',
	code: '',
	type: 'PhoneCheck',
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
const emailForm = ref<checkEmailTypes>({
	email: '',
	code: '',
	type: 'EmailCheck',
});
// 个人认证
const userForm = ref<uploadUserInfoJsontTypes>({
	authenticationMethod: 1,
	userName: '',
	userIDCard: '',
});
// 企业认证
const companyForm = ref<uploadCompanyInfoJsonTypes>({
	authenticationMethod: 1,
	legalPersonName: '',
	legalPersonIDCard: '',
	companyName: '',
	companyID: '',
});
// 微信绑定状态 0未绑定 1已绑定 2服务号绑定 3服务号未绑定
const wxStatus = ref<number>(0);
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
// 财务信息
const financeInfo = ref<any>({
	walletID: '',
	balance: '',
	notBilled: '',
	billed: '',
	hasBilled: '',
	cny: '',
	waitBalance: '',
	withdraw: '',
	waitWithdraw: '',
	notWithdraw: '',
	hasWithdraw: '',
});
const currentInfo = ref({ title: '', content: '', subtitle: '' });
const showContent = (row: any) => {
	currentInfo.value = row;
	currentInfo.value.subtitle = dayjs.unix(row.startAt).format('YYYY-MM-DD') + '至' + dayjs.unix(row.stopAt).format('YYYY-MM-DD');
	showHtmlRef.value.openDialog();
};
const next = () => {
	carouselRef.value.next();
};
const onUpdateController = () => {
	updateLoginCtrRef.value.openDialog(updateLoginCtrForm.value, 'user');
};
//跟新手机号
const updatePhone = () => {
	mobileRef.value.openDialog();
};
// 微信机器人
const updateWxRobot = () => {
	wxRobotRef.value.openDialog('add');
};
const useRegisterCollect = useRegisterApi();
// 删除微信机器人
const deleteWxRobot = () => {
	ElMessageBox.confirm('是否删除微信机器人?', '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(async () => {
			await useRegisterCollect.updateWxRobot({ isDelete: true }).then((res: any) => {
				if (res.code === 'SUCCESS') {
					ElMessage.success('删除微信机器人成功');
					getUserInfo();
				}
			});
		})
		.catch(() => {});
};
//更新用户名
const updateUserName = (type: string) => {
	userNameRef.value.openDialog(type);
};
//更新邮箱
const updateEmail = () => {
	emailRef.value.openDialog();
};
//更新微信
const updateWeixin = () => {
	weixinRef.value.openDialog();
};
// 更新发票抬头
const updateInvoice = () => {
	invoiceRef.value.openDialog();
};
//解绑微信
const delWeixin = () => {
	ElMessageBox.confirm('是否解绑微信?', '提示', {
		confirmButtonText: '确认',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useWeixinApiCollect.bindWeixin({ isDelete: true }).then((res: any) => {
				if (res.code === 'SUCCESS') {
					ElMessage.success('解绑微信成功');
					getUserInfo();
				}
			});
		})
		.catch(() => {});
};
//更新2fa
const bindTotp = () => {
	totpComRef.value.openDialog();
};
// 打开服务号
const openWxService = () => {
	weixinServiceRef.value.openDialog();
};
const imgSuccess = async (base64: string) => {
	avatarLoading.value = true;
	await useUserApiCollect.avatarUpdate({ header: base64, isDelete: false }).then((res: any) => {
		avatarLoading.value = false;
		if (res.code === 'SUCCESS') {
			ElMessage({
				type: 'success',
				message: '头像上传成功!',
			});
		}
	});
	await getUserInfo();
	avatarUploadKey.value++;
};

//删除2fa
const delTotp = () => {
  console.log("AAA", authUser.delete2FAOther)
	if (authUser.delete2FAOther) {
		ElMessageBox.confirm('解绑双因素验证器', '是否解绑2FA验证?', {
			confirmButtonText: '确定',
			cancelButtonText: '取消',
			type: 'warning',
		})
			.then(() => {
				useGoogleAuthApiCollect.rootDelTotp({}).then((res: any) => {
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
	} else {
		ElMessageBox.prompt('请输入双因素验证器验证码', '是否解绑2FA验证?', {
			confirmButtonText: '确定',
			cancelButtonText: '取消',
			type: 'warning',
			inputPattern: /^.+$/,
			inputErrorMessage: '请输入双因素验证器验证码',
		})
			.then(({ value }) => {
				useGoogleAuthApiCollect.delTotp({ code: value }).then((res: any) => {
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
// 单点登录
const handelUpdate = (val: boolean) => {
	const tip = val ? '开启' : '关闭';
	ElMessageBox.confirm('确定要' + tip + '单点登录吗?', '提示', {
		confirmButtonText: '确定',
		cancelButtonText: '取消',
		type: 'warning',
	})
		.then(() => {
			useOnlineUser()
				.updateOnlineUser({ signin: val })
				.then((res: any) => {
					if (res.code === 'SUCCESS') {
						if (val) {
							ElMessage.success('开启成功');
						} else {
							ElMessage.success('关闭成功');
						}
						getUserInfo();
					}
				});
		})
		.catch(() => {});
};
//更新用户信息
const updateUser = () => {
	userRef.value.openDialog(userForm.value);
};
//更新企业'
const updateCompany = () => {
	companyRef.value.openDialog(companyForm.value);
};
//密码
const clickPassword = () => {
	passwordRef.value.openDialog();
};
//更新token
const onToken = () => {
	expirationTokenRef.value.openDialog();
};
//注销账号
const handleLogOff = () => {
	logOffMobileRef.value.openDialog();
};
//重置
const reset = async () => {
	let res = await resetSignIn();
	if (!res) {
		Session.clear();
		Local.clear();
		window.location.reload();
	}
};
const setUserInfo = (data: any) => {
	userInfo.value = data;
	if (userInfo.value?.user) {
		phoneForm.value.phone = userInfo.value.user.phone;
		totpForm.value.phone = userInfo.value.user.phone;
		emailForm.value.email = userInfo.value.user.email;
		signin.value = userInfo.value.user.signin;
	}
	if (userInfo.value?.info) {
		userForm.value.userName = userInfo.value.info.userName;
		userForm.value.userIDCard = userInfo.value.info.userIDCard;
		if (userInfo.value?.info.isCompany) {
			companyForm.value.legalPersonName = userInfo.value.info.legalPersonName;
			companyForm.value.legalPersonIDCard = userInfo.value.info.legalPersonIdCard;
			companyForm.value.companyName = userInfo.value.info.companyName;
			companyForm.value.companyID = userInfo.value.info.companyID;
		}
	}
	if (userInfo.value?.balance) {
		financeInfo.value = userInfo.value.balance;
	}
	if (userInfo.value?.title) {
		invoiceForm.value.name = userInfo.value.title.name;
		invoiceForm.value.taxID = userInfo.value.title.taxID;
		invoiceForm.value.bankID = userInfo.value.title.bandID;
		invoiceForm.value.bank = userInfo.value.title.band;
	}
	if (userInfo.value?.data) {
		if (userInfo.value.data.hasWeChat && userInfo.value.data.hasFuwuhao) {
			wxStatus.value = 1;
		} else if (userInfo.value.data.hasWeChat && !userInfo.value.data.hasFuwuhao) {
			wxStatus.value = 2;
      if (sessionStorage.getItem("fuwuhao-login") === "true") {  // 服务号登录如果解绑了微信必须退出登录
        Session.clear();
        Local.clear();
        window.location.href = '/login'; // 回到登录页
      }
		} else if (!userInfo.value.data.hasWeChat && userInfo.value.data.hasFuwuhao) {
      wxStatus.value = 3;
    } else if (!userInfo.value.data.hasWeChat && !userInfo.value.data.hasFuwuhao) {
      wxStatus.value = 4;
      if (sessionStorage.getItem("fuwuhao-login") === "true") {  // 服务号登录如果解绑了微信必须退出登录
        Session.clear();
        Local.clear();
        window.location.href = '/login'; // 回到登录页
      }
    }
		updateLoginCtrForm.value.allowPhone = userInfo.value.data.allowPhone;
		updateLoginCtrForm.value.allowEmail = userInfo.value.data.allowEmail;
		updateLoginCtrForm.value.allowPassword = userInfo.value.data.allowPassword;
		updateLoginCtrForm.value.allowWeChat = userInfo.value.data.allowWeChat;
		updateLoginCtrForm.value.allowSecondFA = userInfo.value.data.allowSecondFA;
	}
	if (userInfo.value?.user?.phone) {
		if (userInfo.value?.user.userName) {
			userNameForm.value.userName = userInfo.value.user.userName;
		}

		userNameForm.value.nickname = userInfo.value.user.nickname;
		tokenForm.value.tokenExpiration = userInfo.value.user.tokenExpire;
	}
	avatarUploadRef.value && avatarUploadRef.value.onUpdate();
};
const getUserInfo = async () => {
	await useUserApiCollect.userInfo().then((res: any) => {
		if (res.code === 'SUCCESS') {
			userInfo.value = res.data;
			setUserInfo(res.data);
			stores.setUserData(res.data);
			Session.set('userInfo', res.data);
			Session.set('userData', res.data);
			stores.setUserInfos();
		}
	});
};
const socketInfo = useSocketListInfo();
const { announcementList } = storeToRefs(socketInfo);
onMounted(async () => {
	await getUserInfo();
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
