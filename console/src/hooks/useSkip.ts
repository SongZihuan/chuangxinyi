
import { useRouter } from 'vue-router';
import { IsOnlineDomain } from '/@/utils/isonlinedomain';

export default function useSkip() {
	const router = useRouter();

	const onSkip = (homePageData: { link: string }) => {
		if (!homePageData?.link) {
			// ElMessage.warning('抱歉，您没有设置外部链接');
			return;
		}
		// if (IsOnlineDomain(homePageData?.link)) {
		// 	window.location.href = homePageData.link;
		// 	return;
		// }
		const newUrl = router.resolve({
			path: '/externalLinkSkip',
			query: {
				link: homePageData.link,
			},
		});
		window.open(newUrl.href, '_blank');
	};
	return {
		onSkip,
	};
}
