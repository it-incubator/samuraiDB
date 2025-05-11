import * as fs from "fs";
import * as readline from "readline";
import * as path from "path";
import {SSTable} from "../sstable/sstable";
import {SSTablesManager} from "../samurai-db/ss-tables-manager";
import {FileManager} from "../samurai-db/file-manager/file-manager";
import {IS_DELETED_PROP_NAME} from "../samurai-db/samurai-db";

interface LineEntry {
    key: string;
    value: string;
    fileIndex: number;
}

interface FileReader {
    reader: readline.Interface;
    iterator: AsyncIterator<string>;
    currentLine: LineEntry | null;
    done: boolean;
}

export class Compactor {
    constructor(
        private fileManager: FileManager,
        //private directory: string,
        private ssTablesManager: SSTablesManager) {
    }


    async compactTables(levelForProcessing: number): Promise<SSTable[]> {
        const resultLevel = levelForProcessing + 1;
        const sstables: SSTable[] = this.ssTablesManager.ssTablesLevels.get(levelForProcessing);
        if (sstables.length === 0) return [];

        // const outDir = path.join(this.directory, `level${resultLevel}`);
        // if (!fs.existsSync(outDir)) {
        //     fs.mkdirSync(outDir, { recursive: true });
        // }

        // Инициализация читателей файлов
        const fileReaders: FileReader[] = await this.initializeReaders(sstables);

        // Подготовка к записи результата
        const writtenTables: SSTable[] = [];
        let currentChunk: LineEntry[] = [];
        let chunkSize = 0;
        const maxChunkSize = 1024 * 1024; // 1MB или настраиваемый размер
        let fileCounter = 1;

        try {
            while (true) {
                // Находим минимальный ключ среди текущих строк
                const minKeyEntry = this.findMinKeyEntry(fileReaders);
                if (!minKeyEntry) break; // Все файлы обработаны

                // Находим все записи с текущим ключом и выбираем новейшую
                const currentKey = minKeyEntry.currentLine!.key;
                const latestEntry = await this.findLatestEntryForKey(fileReaders, currentKey);

                const valueAsObject = JSON.parse(latestEntry.value);
                if  (valueAsObject[IS_DELETED_PROP_NAME]) {
                    // Продвигаем вперед читатели, которые имели текущий ключ
                    await this.advanceReadersWithKey(fileReaders, currentKey);
                    continue;
                }

                // Добавляем запись в текущий чанк
                currentChunk.push(latestEntry);
                chunkSize += Buffer.byteLength(JSON.stringify(latestEntry), 'utf8');

                // Если чанк достиг максимального размера, записываем его
                if (currentChunk.length >= sstables[0].metaData.count * 2) {
                    const newTable = await this.writeChunkToSSTable(
                        currentChunk,
                        `${fileCounter++}`,
                        resultLevel
                    );
                    writtenTables.push(newTable);
                    currentChunk = [];
                    chunkSize = 0;
                }

                // Продвигаем вперед читатели, которые имели текущий ключ
                await this.advanceReadersWithKey(fileReaders, currentKey);
            }

            // Записываем оставшиеся данные
            if (currentChunk.length > 0) {
                const newTable = await this.writeChunkToSSTable(
                    currentChunk,
                    `${fileCounter}`,
                    resultLevel
                );
                writtenTables.push(newTable);
            }

            return writtenTables;
        } finally {
            // Закрываем все читатели
            fileReaders.forEach(fr => fr.reader.close());
            await Promise.all([...sstables].map(t => this.ssTablesManager.deleteTable(t)));
            // Удаляем исходные файлы

        }
    }

    private async initializeReaders(sstables: SSTable[]): Promise<FileReader[]> {
        const fileReaders: FileReader[] = [];

        for (let i = 0; i < sstables.length; i++) {
            const stream = fs.createReadStream(sstables[i]['dataFilePath']);
            const reader = readline.createInterface({input: stream});
            const iterator = reader[Symbol.asyncIterator]();

            // Пропускаем метаданные (первую строку)
            await iterator.next();

            const fileReader: FileReader = {
                reader,
                iterator,
                currentLine: null,
                done: false
            };

            // Читаем первую реальную строку
            await this.advanceReader(fileReader, i);
            fileReaders.push(fileReader);
        }

        return fileReaders;
    }

    private async advanceReader(fileReader: FileReader, fileIndex: number): Promise<void> {
        const result = await fileReader.iterator.next();

        if (result.done) {
            fileReader.done = true;
            fileReader.currentLine = null;
            return;
        }

        const [key, value] = result.value.split(/:(.+)/);
        fileReader.currentLine = {key, value, fileIndex};
    }

    private findMinKeyEntry(fileReaders: FileReader[]): FileReader | null {
        return fileReaders
            .filter(fr => !fr.done && fr.currentLine)
            .reduce((min, current) => {
                if (!min) return current;
                return BigInt(current.currentLine!.key) < BigInt(min.currentLine!.key) ? current : min;
            }, null as FileReader | null);
    }

    private async findLatestEntryForKey(fileReaders: FileReader[], key: string): Promise<LineEntry> {
        let latestEntry: LineEntry | null = null;

        for (const fr of fileReaders) {
            if (fr.done || !fr.currentLine || fr.currentLine.key !== key) continue;

            if (!latestEntry || fr.currentLine.fileIndex > latestEntry.fileIndex) {
                latestEntry = fr.currentLine;
            }
        }

        return latestEntry!;
    }

    private async advanceReadersWithKey(fileReaders: FileReader[], key: string): Promise<void> {
        const promises = fileReaders.map(async (fr, index) => {
            if (!fr.done && fr.currentLine && fr.currentLine.key === key) {
                await this.advanceReader(fr, index);
            }
        });

        await Promise.all(promises);
    }

    private async writeChunkToSSTable(
        chunk: LineEntry[],
        fileName: string,
        level
    ): Promise<SSTable> {
        const formatted = chunk.map(entry => ({
            key: entry.key,
            value: entry.value
        }));

        const newTable = new SSTable(this.fileManager.getDataFolderPath(), fileName, level);
        await newTable.write(formatted);

        this.ssTablesManager.registerSsTable(newTable)

        return newTable;
    }
}