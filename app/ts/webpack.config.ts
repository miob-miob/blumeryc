import TerserPlugin from 'terser-webpack-plugin'
import path from 'path'

export default {
  entry: {
    bundle: {
      import: './src/index.ts',
      filename: '[name]/index.js',
    },
  },
  output: {
    // coz of lambda fn API
    libraryTarget: 'commonjs2',
    path: path.join(__dirname, 'dist'),
    publicPath: '/',
  },
  target: 'node',
  node: {
    // Need this when working with express, otherwise the build fails
    __dirname: false, // if you don't put this is, __dirname
    __filename: false, // and __filename return blank or /
  },
  optimization: {
    minimize: false,
    // minimize: true,
    minimizer: [new TerserPlugin()],
  },
  module: {
    rules: [
      {
        test: /\.ts$/,
        use: 'ts-loader',
      },
    ],
  },
  resolve: {
    extensions: ['.ts', '.js'],
  },
}
