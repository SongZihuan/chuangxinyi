import { customRef } from 'vue';
export function debounceRef(value: any, duration = 1000) {
	let timer: any;
	return customRef((track, trigger) => {
		return {
			get() {
				track();
				return value;
			},
			set(val) {
				clearTimeout(timer);
				timer = setTimeout(() => {
					trigger();
					value = val;
				}, duration);
			},
		};
	});
}
