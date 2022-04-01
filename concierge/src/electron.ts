const { app, BrowserWindow, ipcMain } = require('electron');

import { ProgressMessage } from './models/progressMessage';
import { AvailableInstallation, PreInstallation } from './models/installation';
const tanzuInstall = require('./backend/tanzu-install.ts');
const platform = require('./os-platform.ts')

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow: Electron.BrowserWindow
let preInstallation: PreInstallation

app.on('ready', initialize);
app.on('activate', activateWindow);
app.on('window-all-closed', quitUnlessDarwin);

function initialize() {
  preInstallation = tanzuInstall.preinstall();
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
      mainWindow.webContents.send('app:pre-install-tanzu', preInstallation)
    })

  // SHIMON FOR NOW:
  mainWindow.webContents.openDevTools();
}

ipcMain.on('app:install-tanzu', async (event, arg) => {
  console.log('Received install-tanzu message; chosenInstallation=' + JSON.stringify(arg))
  const chosenInstallation = arg as AvailableInstallation
  const progressMessenger = { report: (msg: ProgressMessage) => mainWindow.webContents.send('app:install-progress', msg)};
  tanzuInstall.install({...preInstallation, chosenInstallation}, progressMessenger)
});
