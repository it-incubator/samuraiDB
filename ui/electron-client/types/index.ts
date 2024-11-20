
export interface Payload {
  type: 'SET' | 'GET';
  payload: unknown;
  uuid: string;
}

export interface TcpClientAPI {
  connectToServer: (host: string, port: number) => Promise<void>;
  onData: (callback: (payload: string) => void) => void;
  onStatus: (callback: (payload: string) => void) => void;
  sendData: <T extends Payload>(data: T) => Promise<void>;
}