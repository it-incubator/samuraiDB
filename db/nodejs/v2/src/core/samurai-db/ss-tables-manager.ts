import {SSTable} from "../sstable/sstable";
import {FileManager} from "./file-manager/file-manager";
import {MemTable} from "../mem-table/mem-table";
import {IIdManager} from "./i-id-manager";

export class SSTablesManager {
    ssTablesLevels: Map<number, SSTable[]>;

    get sortedSSTablesLevels() {
        const sortedEntries = [...this.ssTablesLevels.entries()].sort((a, b) => a[0] - b[0]);


        return sortedEntries;
    }

    constructor(private fileManager: FileManager, private idManager: IIdManager<any>) {
        this.ssTablesLevels = new Map<number, SSTable[]>()
        this.ssTablesLevels.set(0, [])
    }

    async init() {
        const ssTablesNamesNumbersInLayers = this.fileManager.getSSTablesNumbers();

        for (const [levelNumber, ssTablesNumbers] of ssTablesNamesNumbersInLayers) {
            for (const ssTableName of ssTablesNumbers) {
                let ssTable = new SSTable(this.fileManager.getDataFolderPath(), ssTableName.toString(), levelNumber);
                await ssTable.init()
                this.idManager.setMax(ssTable.metaData.maxId);
                if (!this.ssTablesLevels.has(levelNumber)) {
                    this.ssTablesLevels.set(levelNumber, [])
                }

                this.ssTablesLevels.get(levelNumber).push(ssTable);
            }
        }
    }

    public async flushMemtableToSSTableInZeroLevel(memTable: MemTable<any, any>) {
        const nextSSTableNumber = this.fileManager.getNextSSTableNumberForZeroLevel();
        const newSSTable = new SSTable(this.fileManager.getDataFolderPath(), nextSSTableNumber.toString(), 0);
        await memTable.flush(newSSTable);
        this.fileManager.registerSSTableToZeroLevel(nextSSTableNumber);
        this.ssTablesLevels.get(0).push(newSSTable);
    }
}