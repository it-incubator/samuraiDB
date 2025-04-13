import * as fs from "fs";
import * as path from "path";

export class FileManager {
    private directory: string;
    //private sstablesFilesNumbers: Set<number>;
    private sstablesFilesNumbers: Map<number, Set<number>>;

    constructor(directory: string = "data") {
        this.directory = directory;
        this.sstablesFilesNumbers = new Map<number, Set<number>>()
        this.sstablesFilesNumbers.set(0, new Set());
        this.ensureDirectoryExists();

    }

    init() {
        this.scanExistingSSTables();
    }

    getDataFolderPath() {
        return this.directory;
    }

    public getSSTablesNumbers() {
      return this.sstablesFilesNumbers;
    }

    private ensureDirectoryExists(): void {
        if (!fs.existsSync(this.directory)) {
            fs.mkdirSync(this.directory, { recursive: true });
        }
    }

    private scanExistingSSTables(): void {
        const levels = fs.readdirSync(this.directory, { withFileTypes: true })
            .filter(dirent => dirent.isDirectory())
            .map(dirent => dirent.name);

        for (const levelName of levels) {
            const levelMatch = levelName.match(/^level(\d+)$/);
            if (!levelMatch) continue;

            const levelNumber = parseInt(levelMatch[1], 10);
            const levelPath = path.join(this.directory, levelName);
            const files = fs.readdirSync(levelPath);

            for (const file of files) {
                const match = file.match(/^(\d+)\.sst$/);
                if (!match) continue;

                const fileNumber = parseInt(match[1], 10);

                if (!this.sstablesFilesNumbers.has(levelNumber)) {
                    this.sstablesFilesNumbers.set(levelNumber, new Set());
                }

                this.sstablesFilesNumbers.get(levelNumber)!.add(fileNumber);
            }
        }
    }


    public getNextSSTableNumberForZeroLevel(): number {
        return this.sstablesFilesNumbers.get(0).size + 1;
    }

    public registerSSTableToZeroLevel(number: number): void {
        this.sstablesFilesNumbers.get(0).add(number);
    }
}