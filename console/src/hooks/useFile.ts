import { ref } from 'vue';

export default function useFile() {
	const apiUrl = ref();
	apiUrl.value = import.meta.env.VITE_API_URL;
	const getFile = (fid: string) => {
		return `${apiUrl.value}/public/ui/file?fid=${fid}&download=false&date=${new Date().getTime()}`;
	};
	return {
		getFile,
	};
}
