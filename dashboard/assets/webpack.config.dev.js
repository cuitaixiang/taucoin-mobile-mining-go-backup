// Copyright 2019 The go-tau Authors
// This file is part of the go-tau library.
//
// The go-tau library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-tau library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-tau library. If not, see <http://www.gnu.org/licenses/>.

const webpack = require('webpack');
const merge = require('webpack-merge');
const WebpackDashboard = require('webpack-dashboard/plugin');
const common = require('./webpack.config.common.js');

module.exports = merge(common, {
	mode:    'development',
	plugins: [
		new WebpackDashboard(),
		new webpack.HotModuleReplacementPlugin(),
	],
	// devtool:   'eval',
	devtool:   'source-map',
	devServer: {
		port:     8081,
		hot:      true,
		compress: true,
	},
});
