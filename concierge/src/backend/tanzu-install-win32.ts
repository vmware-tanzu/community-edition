'use strict'

// returns {path, version, edition}
function checkExistingInstallationWin32() {
    return {}
}

function installWin32(progressMessenger) {
    console.log('Installing Win32...')
    progressMessenger({error: true, message: 'Installation on Win32 not yet supported'})
    return false
}

module.exports.install = installWin32
module.exports.checkExistingInstallation = checkExistingInstallationWin32
