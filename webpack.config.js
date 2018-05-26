const webpack = require("webpack");
const ExtractTextPlugin = require("extract-text-webpack-plugin");

module.exports = [{
    entry: './frontend/app/app.js',
    output: {
        path: __dirname + '/frontend/',
        publicPath : __dirname + '/frontend/',
        filename: 'bundle.js'
    },
    module: {
        rules: [{
            test: /\.scss$/,
            use: ExtractTextPlugin.extract({
                fallback: 'style-loader',
                use: ['css-loader','sass-loader']
            })
        },{
            test: /\.(jpe|jpg|woff|woff2|eot|ttf|svg)(\?.*$|$)/,
            loader: 'url-loader?limit=100000'
        }]
    },
    plugins: [
        new webpack.ProvidePlugin({
            $: 'jquery',
            jQuery: 'jquery',
            'window.jQuery': 'jquery',
            Popper: ['popper.js', 'default']
        }),
        new ExtractTextPlugin("bundle.css")
    ]
}];