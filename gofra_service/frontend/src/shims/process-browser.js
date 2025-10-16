// Minimal, side-effect free browser shim for `process`.
// This avoids referencing `process` while defining the shim itself
// so bundlers like Webpack/Vite can safely inject it.
const defaultEnv = {
	NODE_ENV: 'production',
	API_URL: 'http://localhost:8080/api',
}

const globalObject = typeof globalThis !== 'undefined' ? globalThis : window
const existingProcess = globalObject && typeof globalObject.process === 'object'
	? globalObject.process
	: null
const existingEnv = existingProcess && typeof existingProcess.env === 'object'
	? existingProcess.env
	: null

const env = Object.assign({}, defaultEnv, existingEnv || {}, {
	// Allow runtime override via window.__API_URL__ if provided by hosting env
	API_URL:
		typeof globalObject !== 'undefined' &&
		globalObject.__API_URL__ &&
		typeof globalObject.__API_URL__ === 'string'
			? globalObject.__API_URL__
			: (existingEnv && existingEnv.API_URL) || defaultEnv.API_URL,
})

const shim = { env }

if (!existingProcess) {
	globalObject.process = shim
} else if (!existingProcess.env) {
	existingProcess.env = env
}

module.exports = shim
module.exports.default = shim
