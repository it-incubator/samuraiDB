import { FileAdapter } from './file.adapter';

class SamuraiDB {
    private fileAdapter: FileAdapter;
    private index: Map<string, any>;
    /**
     *
     * @param {FileAdapter} fileAdapter
     */
    constructor(fileAdapter: FileAdapter) {
        this.fileAdapter = fileAdapter;
    }

    async init() {
        this.index = await this.fileAdapter.readIndex();
    }

    async set(key: string, data: any) {
        const {offset} = await this.fileAdapter.set(key, data)
        // Обновляем индекс: сохраняем смещение для ключа
        this.index.set(key.toString(), offset);
        await this.fileAdapter.saveIndex(this.index);
    }

    async get(key: string) {
        // Проверяем наличие ключа в индексе
        const offset = this.index.get(key);
        if (offset === undefined) {
            return null; // Если ключа нет в индексе, возвращаем null
        }
        return await this.fileAdapter.get(offset);
    }
}


export default SamuraiDB;