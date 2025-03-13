export enum RetryStrategy {
  EXPONENTIAL = 'exponential',
  FIXED = 'fixed',
}

export interface ModuleOptions {
  host: string;
  port: number;
  maxRetries?: number;
  interval?: number;
  retryStrategy?: RetryStrategy;
}
