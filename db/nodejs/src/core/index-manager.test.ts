import { IndexManager } from './index-manager'; // предположим, что файл с кодом называется index.manager.js
import { FileAdapter } from './file.adapter';
import { join } from 'path';
import { promises as fs } from 'fs';

describe('IndexManager', () => {
  const testDir =  join(__dirname, '..', '..', '__tests__', 'db')

  let fileAdapter;
  let indexManager;

  beforeAll(async () => {
    await fs.mkdir(testDir, { recursive: true });
  });

  beforeEach(async () => {
    fileAdapter = new FileAdapter(testDir);
    indexManager = new IndexManager(fileAdapter);
    await indexManager.init();
  });

  afterEach(async () => {
    await fs.unlink(join(testDir, 'index.txt')).catch(() => {});
  });

  afterAll(async () => {
    await fs.rmdir(testDir, { recursive: true });
  });

  it('should initialize index from file', async () => {
    const indexMap = new Map([['key1', 10]]);
    await fileAdapter.saveIndex(indexMap);

    await indexManager.init();

    const offset = await indexManager.getOffset('key1');
    expect(offset).toBe(10);
  });

  it('should set and save offset correctly', async () => {
    await indexManager.setOffset('key2', 20);

    const offset = await indexManager.getOffset('key2');
    expect(offset).toBe(20);

    // Reinitialize and check persistence
    await indexManager.init();
    const persistedOffset = await indexManager.getOffset('key2');
    expect(persistedOffset).toBe(20);
  });

  it('should return undefined for non-existing key', async () => {
    const offset = await indexManager.getOffset('unknownKey');
    expect(offset).toBeUndefined();
  });
});