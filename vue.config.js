module.exports = {
    publicPath: "/Werewolf/front/",
    devServer: {
        contentBase: '.',
        proxy: {
            "/Werewolf/api": {
                target: "http://localhost:9968/",
                pathRewrite: {
                    "^/Werewolf": "",
                }
            },
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