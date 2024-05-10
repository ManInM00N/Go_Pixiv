
module.exports={
    devServer: {
        proxy: {
            '/apis': {
                target: 'http://127.0.0.1:7234',
                changeOrigin: true,
                pathRewrite: {
                    '^/apis': ''
                }
            }
        }
    }
}