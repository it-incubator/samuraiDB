import { app, BrowserWindow, ipcMain } from 'electron';
import net from 'node:net'; // Импорт модуля net
import { createRequire } from 'node:module';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const require = createRequire(import.meta.url);
const __dirname = path.dirname(fileURLToPath(import.meta.url));

// The built directory structure
//
// ├─┬─┬ dist
// │ │ └── index.html
// │ │
// │ ├─┬ dist-electron
// │ │ ├── main.js
// │ │ └── preload.mjs
// │
process.env.APP_ROOT = path.join(__dirname, '..');

// 🚧 Use ['ENV_NAME'] avoid vite:define plugin - Vite@2.x
export const VITE_DEV_SERVER_URL = process.env['VITE_DEV_SERVER_URL'];
export const MAIN_DIST = path.join(process.env.APP_ROOT, 'dist-electron');
export const RENDERER_DIST = path.join(process.env.APP_ROOT, 'dist');

process.env.VITE_PUBLIC = VITE_DEV_SERVER_URL ? path.join(process.env.APP_ROOT, 'public') : RENDERER_DIST;

let win: BrowserWindow | null;

function createWindow() {
  win = new BrowserWindow({
    icon: path.join(process.env.VITE_PUBLIC, 'electron-vite.svg'),
    webPreferences: {
      preload: path.join(__dirname, 'preload.mjs'),
    },
  });

  // Test active push message to Renderer-process.
  win.webContents.on('did-finish-load', () => {
    win?.webContents.send('main-process-message', new Date().toLocaleString());
  });

  if (VITE_DEV_SERVER_URL) {
    win.loadURL(VITE_DEV_SERVER_URL);
  } else {
    win.loadFile(path.join(RENDERER_DIST, 'index.html'));
  }
}

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
    win = null;
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

app.whenReady().then(createWindow);

// Добавим код для работы с net и IPC

let client: net.Socket | null = null;

ipcMain.handle('connectToServer', async (event, host: string, port: number) => {
  client = new net.Socket();

  client.connect(port, host, () => {
    console.log('Connected to server');
    event.sender.send('server-connection-status', 'Connected');
  });

  client.on('data', (data) => {
    console.log('Received from server:', data.toString());
    event.sender.send('server-data', data.toString());
  });

  client.on('close', () => {
    console.log('Connection closed');
    event.sender.send('server-connection-status', 'Disconnected');
  });

  client.on('error', (error) => {
    console.error('Connection error:', error);
    event.sender.send('server-connection-status', 'Error');
  });
});

// Вы можете добавить другие IPC обработчики для отправки данных на сервер если потребуется

ipcMain.handle('sendData', async (event, data) => {
  if (client) {
    client.write(JSON.stringify(data) + '\n', (err) => {
      if (err) {
        console.error('Error sending data:', err);
        event.sender.send('server-connection-status', 'Error');
      } else {
        console.log('Data sent:', data);
        event.sender.send('server-connection-status', 'Data Sent');
      }
    });
  } else {
    console.error('No active connection to server');
    event.sender.send('server-connection-status', 'No Connection');
  }
});