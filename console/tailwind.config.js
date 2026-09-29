/** @type {import('tailwindcss').Config} */
module.exports = {
	darkMode: 'class',
	important: true,
	corePlugins: {
		preflight: false,
	},
	blocklist: ['container'],
	content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
	theme: {
		extend: {
			height: {},
		},
	},
};
