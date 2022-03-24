const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');
const defaultInclude = path.resolve(__dirname, 'src')

module.exports = [
    {
    mode: 'development',
    entry: './src/electron.ts',
    target: 'electron-main',
    module: {
      rules: [
          {
            test: /\.ts$/,
            include: /src/,
            use: [{ loader: 'ts-loader' }]
          }
      ]
    },
    output: {
      path: __dirname + '/dist',
      filename: 'electron.js'
    }
  },
  {
    mode: 'development',
    entry: './src/frontend/react.tsx',
    target: 'electron-renderer',
    devtool: 'source-map',
    module: {
        rules: [
            {
              test: /\.ts(x?)$/,
              include: /src/,
              use: [{ loader: 'ts-loader' }]
            },
            {
                test: /\.css$/,
                use: [{ loader: 'style-loader' }, { loader: 'css-loader' }, { loader: 'postcss-loader' }],
                include: defaultInclude
            },
        ]
    },
    output: {
      path: __dirname + '/dist',
      filename: 'frontend.js'
    },
    plugins: [
      new HtmlWebpackPlugin({
        template: './src/frontend/index.html'
      })
    ],
      resolve: {
          extensions: ['.ts', '.tsx', '.js', '.css']
      },
  }
];
