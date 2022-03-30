'use strict'

import { ExistingInstallation } from '../models/installation';

function preinstallWin32(): ExistingInstallation {
    return undefined
}

module.exports.steps = []
module.exports.preinstall = preinstallWin32
