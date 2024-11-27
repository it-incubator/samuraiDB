import { join } from 'node:path';
import { mkdirSync, promises as fs } from 'node:fs';

export class FileAdapter {
  readonly filename: string;
  readonly indexFileName: string;

  constructor(dir: string) {

    this.filename = join(dir, 'samuraidb.txt');
    this.indexFileName = join(dir, 'index.txt');

    // Создаем директорию, если она не существует
    mkdirSync(dir, { recursive: true });
  }

  async set(key: string, data: any) {
    // Сериализуем данные в JSON и создаем строку в формате "ключ,значение"
    const entry = `${key},${JSON.stringify(data)}\n`;

    // Открываем файл для получения текущего смещения
    const fileHandle = await fs.open(this.filename, 'a');
    const offset = (await fileHandle.stat()).size;

    // Добавляем запись в файл
    await fs.appendFile(this.filename, entry);

    // Закрываем файл
    await fileHandle.close();

    return { offset };
  }

  async get(offset: number) {
    // Проверяем наличие ключа в индексе
    if (offset === undefined) {
      throw new Error('Offset must be passed'); // Если ключа нет в индексе, возвращаем null
    }

    // Открываем файл и переходим к нужному смещению
    const fileHandle = await fs.open(this.filename, 'r');

    // Читаем строку с ключом и значением
    const buffer = Buffer.alloc(1024); // Создаем буфер для чтения строки
    await fileHandle.read(buffer, 0, 1024, offset);
    const line = buffer.toString('utf-8').trim();

    // Закрываем файл
    await fileHandle.close();

    const [storedKey, storedValue] = line.split(/,(.+)/);

    return JSON.parse(storedValue);
  }

  async saveIndex(indexMap: Map<string, any>) {
    const serializedMap = JSON.stringify(Array.from(indexMap));
    await fs.writeFile(this.indexFileName, serializedMap, 'utf-8');
  }

  async readIndex() {
    try {
      const fileContent: string | undefined = await fs.readFile(this.indexFileName, 'utf-8');

      if (!fileContent) {
        // File is empty, return an empty Map
        return new Map();
      }

      return new Map(JSON.parse(fileContent));
    } catch (error: any) {
      if (error?.code === 'ENOENT') {
        // File not found, return an empty Map
        return new Map();
      }
      // Re-throw the error if it's not a "file not found" error
      throw error;
    }
  }
}