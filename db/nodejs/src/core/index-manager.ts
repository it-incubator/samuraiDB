import { FileAdapter } from "./file.adapter";

export class IndexManager {
  private index: Map<string, any>

  constructor(protected fileAdapter: FileAdapter) {}

  async init() {
    this.index = await this.fileAdapter.readIndex();
  }

  async setOffset(key: string, offset: number) {
    this.index.set(key, offset);
    await this.fileAdapter.saveIndex(this.index);
  }

  async getOffset(key: string) {
    return this.index.get(key);
  }
}
