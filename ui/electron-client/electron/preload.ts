import { ipcRenderer, contextBridge } from 'electron'
import {  TcpClientAPI } from '../types';

// --------- Expose some API to the Renderer process ---------
contextBridge.exposeInMainWorld('ipcRenderer', {
  on(...args: Parameters<typeof ipcRenderer.on>) {
    const [channel, listener] = args
    return ipcRenderer.on(channel, (event, ...args) => listener(event, ...args))
  },
  off(...args: Parameters<typeof ipcRenderer.off>) {
    const [channel, ...omit] = args
    return ipcRenderer.off(channel, ...omit)
  },
  send(...args: Parameters<typeof ipcRenderer.send>) {
    const [channel, ...omit] = args
    return ipcRenderer.send(channel, ...omit)
  },
  invoke(...args: Parameters<typeof ipcRenderer.invoke>) {
    const [channel, ...omit] = args
    return ipcRenderer.invoke(channel, ...omit)
  },

  // You can expose other APTs you need here.
  // ...
})



contextBridge.exposeInMainWorld('tcpClient', {

  connectToServer: (host: string, port: number) => ipcRenderer.invoke('connectToServer', host, port),
  onData: (callback: (payload: string) => void) => ipcRenderer.on('server-data', (_event, data) => callback(data)),
  onStatus: (callback: (payload: string) => void) => ipcRenderer.on('server-connection-status', (_event, data) => callback(data)),
  sendData: <T extends TcpClientAPI>(data: T) => ipcRenderer.invoke('sendData', data)
})