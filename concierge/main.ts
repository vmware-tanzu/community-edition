'use strict'
// Import parts of electron to use
const { app, BrowserWindow, ipcMain } = require('electron')
const os = require( 'os' );
const path = require('path')
const url = require('url')

const tanzuInstall = require('./src/backend/tanzu-install')

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let mainWindow

// Keep a reference for dev mode
let dev = false

// Broken:
// if (process.defaultApp || /[\\/]electron-prebuilt[\\/]/.test(process.execPath) || /[\\/]electron[\\/]/.test(process.execPath)) {
//   dev = true
// }

if (process.env.NODE_ENV !== undefined && process.env.NODE_ENV === 'development') {
  dev = true
}

// Temporary fix broken high-dpi scale factor on Windows (125% scaling)
// info: https://github.com/electron/electron/issues/9691
if (process.platform === 'win32') {
  app.commandLine.appendSwitch('high-dpi-support', 'true')
  app.commandLine.appendSwitch('force-device-scale-factor', '1')
}

ipcMain.on('app:install-tanzu', async (event, arg) => {
  console.log('Received install-tanzu message; arg=' + arg)
  tanzuInstall.install(progressMessageObject => {
    mainWindow.webContents.send('app:install-progress', progressMessageObject)
  })
});

function createWindow() {
  // Create the browser window.
  mainWindow = new BrowserWindow({
    width: 1024,
    height: 768,
    show: false,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false
    }
  })

  mainWindow.webContents.on('did-finish-load', () => {
    mainWindow.webContents.send('app:existing-install-tanzu', tanzuInstall.checkExistingInstallation())
  })

  // and load the index.html of the app.
  let indexPath

  if (dev && process.argv.indexOf('--noDevServer') === -1) {
    indexPath = url.format({
      protocol: 'http:',
      host: 'localhost:8080',
      pathname: 'index.html',
      slashes: true
    })
  } else {
    indexPath = url.format({
      protocol: 'file:',
      pathname: path.join(__dirname, 'dist', 'index.html'),
      slashes: true
    })
  }

  mainWindow.loadURL(indexPath)

  // Don't show until we are ready and loaded
  mainWindow.once('ready-to-show', () => {
    mainWindow.show()

    // Open the DevTools automatically if developing
    if (dev) {
      const { default: installExtension, REACT_DEVELOPER_TOOLS } = require('electron-devtools-installer')

      installExtension(REACT_DEVELOPER_TOOLS)
        .catch(err => console.log('Error loading React DevTools: ', err))
      // SHIMON: let's not: mainWindow.webContents.openDevTools()
    }
  })

  // Emitted when the window is closed.
  mainWindow.on('closed', function() {
    // Dereference the window object, usually you would store windows
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    mainWindow = null
  })
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  createWindow()
  console.log('HOME DIRECTORY: ' + os.homedir())
  console.log('app dir: ' + __dirname)
})

// Quit when all windows are closed.
app.on('window-all-closed', () => {
  // On macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('activate', () => {
  // On macOS it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (mainWindow === null) {
    createWindow()
  }
})
