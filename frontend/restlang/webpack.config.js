const path = require("path");
const pluginLiveReload = require('webpack-livereload-plugin');

module.exports = {
    entry: "./src/index.ts",
    mode: "development",
    module: {
        rules: [
            {
                test: /\.ne$/,
                use: [
                    'nearley-loader',
                ],
            },
            {
                test: /\.tsx?$/,
                use: "ts-loader",
                exclude: /node_modules/,
            },
        ],
    },
    resolve: {
        extensions: [".tsx", ".ts", ".js"],
    },
    output: {
        filename: "app.js",
        path: path.resolve(__dirname, ".."),
    },
    plugins: [
        new pluginLiveReload({})
    ]
};