module.exports = {
    publicPath:"/front/",
    devServer: {
        contentBase: '.',
        proxy: {
            "/api": {target: "http://localhost:9968/"},
        },
        host: "0.0.0.0",
        port: 9977,
        hot: true,
        disableHostCheck: true,
        historyApiFallback: {
            rewrites: [{from: /^\/$/, to: "/Index.html"}]
        }
    }
}