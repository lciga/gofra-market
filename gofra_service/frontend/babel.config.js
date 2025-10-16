module.exports = function (api) {
	api.cache(true)

	return {
		presets: ['@quasar/babel-preset-app'],
	}
}
