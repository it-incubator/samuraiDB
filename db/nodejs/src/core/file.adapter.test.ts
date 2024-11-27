import { FileAdapter } from './file.adapter'; // предположим, что файл с твоим кодом называется fileAdapter.js
import { promises as fs } from 'node:fs';
import { join } from 'node:path';

describe('FileAdapter', () => {
  const testDir = join(__dirname, '..', '..', '__tests__', 'db')
  const testFilename = join(testDir, 'samuraidb.txt');
  const testIndexFilename = join(testDir, 'index.txt');

  beforeAll(async () => {
    // Создаем тестовую директорию для файлового адаптера
    await fs.mkdir(testDir, { recursive: true });
  });

  afterEach(async () => {
    // Удаляем созданные тестовые файлы после каждого теста
    await fs.unlink(testFilename).catch(() => {});
    await fs.unlink(testIndexFilename).catch(() => {});
  });

  afterAll(async () => {
    // Удаляем тестовую директорию после всех тестов
    await fs.rm(testDir, { recursive: true }).catch(() => {});
  });

  it('should set and get data correctly', async () => {
    const fileAdapter = new FileAdapter(testDir);
    const key = 'testKey';
    const data = { name: 'Samurai' };

    const { offset } = await fileAdapter.set(key, data);
    const retrievedData = await fileAdapter.get(offset);

    expect(retrievedData).toEqual(data);
  });

  it('should save and read index correctly', async () => {
    const fileAdapter = new FileAdapter(testDir);
    const indexMap = new Map([['testKey', 0], ['anotherKey', 42]]);

    await fileAdapter.saveIndex(indexMap);
    const readIndex = await fileAdapter.readIndex();

    expect(readIndex).toEqual(indexMap);
  });

  it('should handle non-existent index file gracefully', async () => {
    const fileAdapter = new FileAdapter(testDir);

    const readIndex = await fileAdapter.readIndex();

    expect(readIndex).toEqual(new Map());
  });
});