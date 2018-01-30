const webpack = require("webpack");
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const UglifyJSPlugin = require('uglifyjs-webpack-plugin');

const indexExtractLess = new ExtractTextPlugin({filename: "index.bundle.css"});
const formExtractLess = new ExtractTextPlugin({filename: "form.bundle.css"});

module.exports = [{
    entry: './frontend/app/app-index.js',
    output: {
        path: __dirname + '/frontend/dist/',
        filename: 'index.bundle.js'
    },
    module: {
        rules: [{
            test: /\.scss/,
            use: indexExtractLess.extract({
                use: [{
                    loader: "css-loader", options: { minimize: true }
                },{
                    loader: "sass-loader"
                }],
                fallback: "style-loader"
            })
        }, {
            test: /\.css/,
            use: indexExtractLess.extract({
                use: [{
                    loader: "css-loader", options: { minimize: true }
                }],
                fallback: "style-loader"
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
        new UglifyJSPlugin({
            extractComments: true
        }),
        indexExtractLess
    ]
}, {
    entry: './frontend/app/app.js',
    output: {
        path: __dirname + '/frontend/dist/',
        filename: 'form.bundle.js'
    },
    module: {
        rules: [{
            test: /\.scss/,
            use: formExtractLess.extract({
                use: [{
                    loader: "css-loader", options: { minimize: true }
                },{
                    loader: "sass-loader"
                }],
                fallback: "style-loader"
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
        new UglifyJSPlugin({
            extractComments: true
        }),
        formExtractLess
    ]
}];