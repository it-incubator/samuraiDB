import * as fs from "fs";
import * as path from "path";
import * as readline from "readline";

export type MetaDataType = {
    minId: string
    maxId: string
    count: number
}

export class SSTable {
    private dataFilePath: string;
    private indexFilePath: string;
    public metaData: MetaDataType;
    private index: Map<string, number>;
    private status: 'initialized' | 'not-initialized' = 'not-initialized'
    public readonly layerNumber: number;

    constructor(rootDirectory: string, fileName: string, layerNumber: number) {
        this.layerNumber = layerNumber;

        const levelFolderPath = path.join(rootDirectory, `level${layerNumber}`);
        fs.mkdirSync(levelFolderPath, { recursive: true });

        this.dataFilePath = path.join(levelFolderPath, `${fileName}.sst`);
        this.indexFilePath = path.join(levelFolderPath, `${fileName}.idx`);

        this.index = new Map();

    }

    async init() {
        if (this.status === 'not-initialized') {
            await this.loadIndex();
            await this.loadMetadata();
            this.status = 'initialized';
        } else {
            throw new Error('already initialized');
        }
    }

    private loadIndex(): void {
        if (!fs.existsSync(this.indexFilePath)) return;

        const indexData = fs.readFileSync(this.indexFilePath, "utf-8");
        indexData.split("\n").forEach(line => {
            if (!line.trim()) return;
            const [key, pos] = line.split(":");
            this.index.set(key, Number(pos));
        });
    }

    private async loadMetadata(): Promise<void> {
        const stream = fs.createReadStream(this.dataFilePath, { start: 0, encoding: "utf-8" });
        const rl = readline.createInterface({ input: stream });

        const promise = new Promise<MetaDataType>((resolve) => {
            rl.once("line", (line) => {
                rl.close();
                stream.destroy(); // Close the stream properly
                const [storedKey, storedValue] = line.split(/:(.+)/);
                resolve(JSON.parse(storedValue));
            });

            rl.once("error", () => resolve(null)); // Handle potential stream errors
            stream.once("error", () => resolve(null)); // Handle stream errors
        });

        this.metaData = await promise;
    }

    write(data: { key: string; value: string }[]): Promise<void> {
        const metadata: MetaDataType = {
            minId: data[0].key.toString(),
            maxId: data.at(-1).key.toString(),
            count:  data.length
        }

        this.metaData = metadata;

        return new Promise((resolve, reject) => {
            const dataStream = fs.createWriteStream(this.dataFilePath, { flags: "w" });
            const indexStream = fs.createWriteStream(this.indexFilePath, { flags: "w" });


            let metadataLine = `meta:${JSON.stringify(metadata)}\n`;
            dataStream.write(metadataLine);

            let position = Buffer.byteLength(metadataLine, "utf-8"); // Отслеживаем смещение вручную

            data.forEach(({ key, value }) => {
                const line = `${key}:${value}\n`;

                dataStream.write(line);
                indexStream.write(`${key}:${position}\n`); // Записываем смещение в индекс-файл

                this.index.set(key.toString(), position); // Записываем смещение в память

                position += Buffer.byteLength(line, "utf-8"); // Обновляем смещение
            });

            dataStream.end();
            indexStream.end();

            let completed = 0;
            const checkCompletion = () => {
                if (++completed === 2) resolve();
            };

            dataStream.on("finish", checkCompletion);
            indexStream.on("finish", checkCompletion);
            dataStream.on("error", reject);
            indexStream.on("error", reject);
        });
    }

    read(key: string): Promise<string | null> {
        if (!this.index.has(key)) return Promise.resolve(null);

        const position = this.index.get(key)!;
        const stream = fs.createReadStream(this.dataFilePath, {
            start: position,
            encoding: "utf-8"
        });

        return new Promise((resolve, reject) => {
            let buffer = "";

            stream.on("data", (chunk) => {
                buffer += chunk;
                const newlineIndex = buffer.indexOf("\n");

                if (newlineIndex !== -1) {
                    stream.destroy(); // Остановить поток как только нашли первую строку

                    const line = buffer.slice(0, newlineIndex);
                    console.log('👀 LINE:', line);

                    const [storedKey, storedValue] = line.split(/:(.+)/);
                    console.log('🎹:', storedValue);

                    try {
                        resolve(storedValue ? JSON.parse(storedValue) : null);
                    } catch {
                        resolve(null);
                    }
                }
            });

            stream.on("error", () => resolve(null));
            stream.on("end", () => resolve(null)); // если \n не встретился
        });
    }

    delete(): void {
        this.index = new Map();
        // todo: move to filemanager
        if (fs.existsSync(this.dataFilePath)) fs.unlinkSync(this.dataFilePath);
        if (fs.existsSync(this.indexFilePath)) fs.unlinkSync(this.indexFilePath);
    }
}