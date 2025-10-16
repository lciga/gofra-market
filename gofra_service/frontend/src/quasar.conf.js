module.exports = function (ctx) {
    return {
        framework: {
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
            vueRouterMode: 'history'
        }
    }
}