'use strict'
const tanzuDarwin = require('./tanzu-install-darwin.ts')
const tanzuWin32 = require('./tanzu-install-win32.ts')

if (process.platform === 'darwin') {
    module.exports = tanzuDarwin
} else if (process.platform === 'win32') {
    module.exports = tanzuWin32
}
