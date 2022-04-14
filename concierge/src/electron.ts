const { app, BrowserWindow, ipcMain } = require('electron');

import { ProgressMessage } from './models/progressMessage';
import { AvailableInstallation, PreInstallation } from './models/installation';
const tanzuInstall = require('./backend/tanzu-install.ts');
const platform = require('./os-platform.ts')
const utils = require('./utils.ts')
import { Messages } from './models/messages'

/***************************************************************************
 * As of 2022.04.12 we are setting this project down.
 * Here are some notes in case we pick it back up at some point:
 *
 * 1. Rather than have the current set up steps mapped out in advance, it
 * may make more sense to think of the installation in three phases:
 * A. installation prep, consisting of setting up paths and unpacking the
 * chosen archive
 * B. after unpacking the archive, putting together an array of steps (pass
 * to the front end side), which would include installing each plugin individually,
 * any configuration changes, etc. Each of these could be a separate step, and
 * the frontend could request each step one at a time.
 * C. after installation, verify the installation (or cleanup after a cancel)
 *
 * This approach would also have the benefit of inverting the control (so
 * that the renderer thread would drive the process rather than the backend)
 * which would allow the app to be more responsive (and the user to cancel)
 *
 * 2. The full styling was never implemented. For mockups talk with Louis Weitzman
 *
 * 3. The frontend inter-process communication (ipc) code was left in the App.tsx class
 * and the intention is to isolate it further into its own module. There is a
 * registration pattern that allows the ipc code to set the state inside of a
 * component (see the registerXXX methods).
 * There may be a better pattern to achieve that functionality.
 *
 * 4. The backend inter-process communication (ipc) code was left in this module, but
 * it feels small enough not to be an issue.
 *
 * 5. There was a start at documenting what was expected in the props object passed
 * to each component, as well as a verifyProps() method that would enforce that
 * expectation. It would be good to complete those tasks.
/***************************************************************************
*/

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow: Electron.BrowserWindow
let preInstallation: PreInstallation

function progressMessenger(window) {
  if (!window) {
    return { report: (msg: ProgressMessage) => console.log(`MESSAGE: ${JSON.stringify(msg)}`) }
  }
  return { report: (msg: ProgressMessage) => window.webContents.send(`${Messages.RESPONSE_PROGRESS}`, msg)}
}

app.on('ready', initialize);
app.on('activate', activateWindow);
app.on('window-all-closed', quitUnlessDarwin);

function initialize() {
  utils.fixPath()
  preInstallation = tanzuInstall.preinstall(progressMessenger(mainWindow));
  console.log('PREINSTALL RESULT: ' + JSON.stringify(preInstallation));
  createWindow();
}

function activateWindow() {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (mainWindow === null) {
      createWindow()
    }
}

function quitUnlessDarwin() {
  // Quit when all windows are closed... except on macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== platform.osMac) {
    app.quit()
  }
}

function createWindow () {
  // Create the browser window.
  mainWindow = new BrowserWindow({
    width: 1800, // 800,
    height: 900, // 600,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false
    }
  });  // and load the index.html of the app.
  mainWindow.loadFile('index.html');
    mainWindow.webContents.on('did-finish-load', () => {
      mainWindow.webContents.send(`${Messages.RESPONSE_PRE_INSTALL}`, preInstallation)
    })

  // SHIMON FOR NOW:
  mainWindow.webContents.openDevTools();
}

ipcMain.on(`${Messages.REQUEST_INSTALL}`, async (event, arg) => {
  console.log(`Received ${Messages.REQUEST_INSTALL} message; chosenInstallation=JSON.stringify(arg)`)
  const chosenInstallation = arg as AvailableInstallation
  tanzuInstall.install({...preInstallation, chosenInstallation}, progressMessenger(mainWindow))
});

ipcMain.on(`${Messages.REQUEST_PRE_INSTALL}`, async (event) => {
  console.log(`Received ${Messages.REQUEST_PRE_INSTALL} message`)
  preInstallation = tanzuInstall.preinstall(progressMessenger(mainWindow));
  console.log('PREINSTALL RESULT: ' + JSON.stringify(preInstallation));
  mainWindow.webContents.send(`${Messages.RESPONSE_PRE_INSTALL}`, preInstallation)
});

ipcMain.on(`${Messages.REQUEST_PRE_INSTALL}`, event => {
  const pluginList = tanzuInstall.pluginList(progressMessenger(mainWindow))
  mainWindow.webContents.send(`${Messages.RESPONSE_PRE_INSTALL}`, pluginList)
})

ipcMain.on(`${Messages.REQUEST_LAUNCH_KICKSTART}`, event => {
  tanzuInstall.launchKickstart(progressMessenger(mainWindow))
})

ipcMain.on(`${Messages.REQUEST_LAUNCH_UI}`, event => {
  tanzuInstall.launchTanzuUi(progressMessenger(mainWindow))
})

