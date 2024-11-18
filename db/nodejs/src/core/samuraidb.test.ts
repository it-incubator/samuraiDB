import SamuraiDB from './samuraidb';
import { promises as fs } from 'node:fs';
import { FileAdapter } from './file.adapter';
import { join } from 'node:path';

describe('SamuraiDB', () => {
  const dir = join(__dirname, '..', '..', '__tests__', 'db')
  let db: SamuraiDB;
  beforeEach(async () => {
    const fileAdapter = new FileAdapter(dir, );
    db = new SamuraiDB(fileAdapter);

    await db.init();

    fs.readFile = jest.fn();

    if (jest.isMockFunction(fs.appendFile)) {
      (fs.appendFile as jest.Mock).mockRestore();
    }
  });

  test('set should append data to the file', async () => {
    const key = 'testKey';
    const data = { name: 'John Doe', age: 30 };
    const append = jest.spyOn(fs, 'appendFile').mockImplementation(async () => {
    });

    await db.set(key, data);

    expect(append).toHaveBeenCalledWith(expect.any(String), `${key},${JSON.stringify(data)}\n`);
  });

  test('get should return the latest value for a given key', async () => {
    const key = 'testKey';
    const data1 = { name: 'John Doe', age: 30 };
    const data2 = { name: 'Jane Doe', age: 25 };

    await db.set(key, data1);
    await db.set(key, data2);

    const result = await db.get(key);

    expect(result).toEqual(data2);
  });

  test('get should return null if the key is not found', async () => {
    const key = 'nonExistentKey';
    const anotherKey = 'anotherKey';
    const data = '{"value":"test"}\n';

    await db.set(anotherKey, data);

    const result = await db.get(key);

    expect(result).toBeNull();
  });
});