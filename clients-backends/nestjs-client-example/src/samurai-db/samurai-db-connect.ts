import { Injectable } from '@nestjs/common';
import * as net from 'net';

@Injectable()
export class SamuraiDBConnect {
  protected tcpClient: net.Socket;
  private readonly host: string;
  private readonly port: number;
  private retryInterval: number;
  private readonly maxRetries: number;
  private attempt: number;
  protected status: 'CONNECTING' | 'CONNECTED' = 'CONNECTING';

  constructor(
    host: string,
    port: number,
    maxRetries = 5,
    initialRetryInterval = 1000,
  ) {
    this.host = host;
    this.port = port;
    this.maxRetries = maxRetries;
    this.retryInterval = initialRetryInterval;
    this.attempt = 0; // Инициализация попыток с 0
    this.connect();
  }

  private connect(): void {
    this.tcpClient = net.createConnection(
      { host: this.host, port: this.port },
      () => {
        console.log('Connected to server');
        this.status = 'CONNECTED';
        this.retryInterval = 1000; // Сброс интервала при успешном подключении
        this.attempt = 0; // Сброс счетчика попыток при успешном подключении
      },
    );

    this.tcpClient.on('error', (err) => {
      console.error('Connection error:', err.message);
      this.status = 'CONNECTING';
      this.attempt++; // Увеличение счетчика попыток

      if (this.attempt <= this.maxRetries) {
        console.log(
          `Attempt ${this.attempt} failed. Retrying in ${this.retryInterval / 1000}s...`,
        );
        setTimeout(() => this.connect(), this.retryInterval);
        this.retryInterval *= 2; // Увеличение времени задержки
      } else {
        console.error('Max retries reached. Please check the server.');
      }
    });
  }
}
