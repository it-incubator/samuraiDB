import { Socket } from 'net';

export abstract class ConnectionService {
  public client: Socket;

  protected abstract connect(): void;
}
