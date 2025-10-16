module.exports = function () {
    return {
        boot: ['axios'],
        extras: ['material-icons'],
        framework: {
            lang: 'ru',
            iconSet: 'material-icons',
            config: {
                brand: {
                    primary: '#D4B896',
                    secondary: '#8B4513',
                    accent: '#FF6B35',
                    dark: '#2E2E2E',
                    positive: '#21BA45',
                    negative: '#C10015',
                    info: '#31CCEC',
                    warning: '#F2C037'
                }
            },
            plugins: ['Notify', 'Dialog']
        },
        build: {
            vueRouterMode: 'history',
            env: {
                // default API URL for production bundle (can be overridden at runtime)
                API_URL: process.env.API_URL || 'http://localhost:8080/api'
            },
            extendWebpack(cfg) {
                const webpack = require('webpack')
                // Provide a browser-friendly `process` so libraries referencing
                // process.env.* don't throw at runtime.
                const { resolve } = require('path')
                cfg.plugins.push(
                    new webpack.ProvidePlugin({
                        process: resolve(__dirname, 'src/shims/process-browser.js')
                    })
                )
                // Ensure NODE_ENV is defined inside the bundle
                cfg.plugins.push(
                    new webpack.DefinePlugin({
                        'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'production')
                    })
                )
            }
        }
    }
}
