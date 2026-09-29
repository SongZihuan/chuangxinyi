<template>
	<div class="container mx-auto h-[100%] mb-[20px]">
		<div class="pt-[20px] text-center text-[#000] text-[32px] font-bold">产品服务套餐价格</div>
		<div class="pt-[10px] pb-[30px] text-center text-[#000] text-[14px]">提供初级版、高级版、旗舰版、企业版多种选择；年付季付均有折扣。</div>
		<div class="flex justify-center items-center">
			<div class="w-[276px] bg-white border-[1px] border-[#FBFCFF] mx-3 relative left-0 top-0 shadow-xl">
				<div class="box-bg p-[16px] bg-[url('../../src/assets/images/box-bg.png')] bg-no-repeat bg-[100%_100%]">
					<div
						class="flex items-center justify-end absolute text-[#fff] text-right pr-[12px] text-[14px] top-0 right-0 w-[100px] h-[24px] bg-[url('../../src/assets/images/tips.png')] bg-no-repeat bg-auto"
						style="background-size: 100%"
					>
						永久免费
					</div>
					<div class="text-[#000] text-[20px] font-bold">初级版</div>
					<div class="text-[#000] text-[14px] mt-1">注册即可拥有</div>
					<div class="flex flex-col justify-start items-start mt-[20px] h-[650px]">
						<div class="text-[#000] text-[14px] my-1">&#9989; 免费版任意共10条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 短链（免费版）</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 活码（免费版）</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 外链（免费版）</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 抖音卡片（免费版）</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 微信卡片（免费版）</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 访问数据</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 开放API</div>

						<div class="text-[#000] text-[14px] my-1">&#9940; 人工客服</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 客户群服务</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 自定义域名</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 服务质量保障</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 合同签约</div>
					</div>
				</div>
				<div class="w-[100%] h-[61px] bg-[#ECF2FF] mt-[20px] flex justify-center items-center text-[#4D79FF] cursor-pointer" @click="toBuy(8)">
					立即购买，8元/月
				</div>
			</div>
			<div class="w-[276px] bg-white border-[1px] border-[#FBFCFF] mx-3 relative left-0 top-0 shadow-xl">
				<div class="box-bg p-[16px] bg-[url('../../src/assets/images/box-bg.png')] bg-no-repeat bg-[100%_100%]">
					<div
						class="flex items-center justify-end absolute text-[#fff] text-right pr-[12px] text-[14px] top-0 right-0 w-[100px] h-[24px] bg-[url('../../src/assets/images/tips.png')] bg-no-repeat bg-auto"
						style="background-size: 100%"
					>
						限时活动
					</div>
					<div class="text-[#000] text-[20px] font-bold">高级版</div>
					<div class="text-[#000] text-[14px] mt-1">季付9折，年付8.5折</div>
					<div class="flex flex-col justify-start items-start mt-[20px] h-[650px]">
						<div class="text-[#000] text-[14px] my-1">&#9989; 普通短链40条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 静默短链10条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; HTTP短链10条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 普通活码15条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 短链活码5条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 活码文件100个</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 普通外链10条</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 短链外链0条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 普通抖音卡片30条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 静默抖音卡片10条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 短链接抖音卡片10条</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 静默短链接抖音卡片0条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 普通微信卡片10条</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 普通静默微信卡片0条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 一体化微信卡片20条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 一体化静默微信卡片20条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 访问数据</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 开放API</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 5x8小时 人工客服</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 客户群服务</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 自定义域名</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 服务质量保障</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 合同签约</div>
					</div>
				</div>
				<div class="w-[100%] h-[61px] bg-[#ECF2FF] mt-[20px] flex justify-center items-center text-[#4D79FF] cursor-pointer" @click="toBuy(68)">
					立即购买，68元/月
				</div>
			</div>
			<div class="w-[276px] bg-white border-[1px] border-[#FBFCFF] mx-3 relative left-0 top-0 shadow-xl">
				<div class="box-bg p-[16px] bg-[url('../../src/assets/images/box-bg.png')] bg-no-repeat bg-[100%_100%]">
					<div
						class="flex items-center justify-end absolute text-[#fff] text-right pr-[12px] text-[14px] top-0 right-0 w-[100px] h-[24px] bg-[url('../../src/assets/images/tips.png')] bg-no-repeat bg-auto"
						style="background-size: 100%"
					>
						限时活动
					</div>
					<div class="text-[#000] text-[20px] font-bold">旗舰版</div>
					<div class="text-[#000] text-[14px] mt-1">季付9折，年付8.5折</div>
					<div class="flex flex-col justify-start items-start mt-[20px] h-[650px]">
						<div class="text-[#000] text-[14px] my-1">&#9989; 普通短链300条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 静默短链100条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; HTTP短链100条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 普通活码60条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 短链活码20条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 活码文件500个</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 普通外链60条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 短链外链20条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 普通抖音卡片250条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 静默抖音卡片100条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 短链接抖音卡片100条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 静默短链接抖音卡片50条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 普通微信卡片80条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 普通静默微信卡片20条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 一体化微信卡片200条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 一体化静默微信卡片200条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 访问数据</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 开放API</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 7x12小时 人工客服</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 客户群服务</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 自定义域名</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 服务质量保障</div>
						<div class="text-[#000] text-[14px] my-1">&#9940; 合同签约</div>
					</div>
				</div>
				<div class="w-[100%] h-[61px] bg-[#ECF2FF] mt-[20px] flex justify-center items-center text-[#4D79FF] cursor-pointer" @click="toBuy(118)">
					立即购买，118元/月
				</div>
			</div>
			<div class="w-[276px] bg-white border-[1px] border-[#FBFCFF] mx-3 relative left-0 top-0 shadow-xl">
				<div class="box-bg p-[16px] bg-[url('../../src/assets/images/box-bg.png')] bg-no-repeat bg-[100%_100%]">
					<div
						class="flex items-center justify-end absolute text-[#fff] text-right pr-[12px] text-[14px] top-0 right-0 w-[100px] h-[24px] bg-[url('../../src/assets/images/tips.png')] bg-no-repeat bg-auto"
						style="background-size: 100%"
					>
						限时活动
					</div>
					<div class="text-[#000] text-[20px] font-bold">企业版</div>
					<div class="text-[#000] text-[14px] mt-1">季付8折，年付7.5折</div>
					<div class="flex flex-col justify-start items-start mt-[20px] h-[650px]">
						<div class="text-[#000] text-[14px] my-1">&#9989; 短链10000条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 活码6000条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 活码文件20000个</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 外链6000条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 抖音卡片10000条</div>

						<div class="text-[#000] text-[14px] my-1">&#9989; 微信卡片10000条</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 访问数据</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 开放API</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 7x24小时 人工客服</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 客户群服务</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 自定义域名1个</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 定制化</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 服务质量保障（赔付）</div>
						<div class="text-[#000] text-[14px] my-1">&#9989; 合同签约</div>
					</div>
				</div>
				<div class="w-[100%] h-[61px] bg-[#ECF2FF] mt-[20px] flex justify-center items-center text-[#4D79FF] cursor-pointer" @click="toBuy(2880)">
					联系客服，2,880元/月起
				</div>
			</div>
		</div>
		<rechargeDialog ref="rechargeDialogRef" />
	</div>
</template>

<script setup lang="ts">
import rechargeDialog from './component/rechargeDialog.vue';
import { ref } from 'vue';
const rechargeDialogRef = ref();
const toBuy = (price: number) => {
	rechargeDialogRef.value.openDialog(price);
};
</script>

<style lang="scss" scoped>
.box-bg {
	background-size: 100%;
}
</style>
