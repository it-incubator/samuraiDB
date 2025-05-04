import {SSTable} from "../sstable/sstable";
import {FileManager} from "./file-manager/file-manager";
import {MemTable} from "../mem-table/mem-table";
import {IIdManager} from "./i-id-manager";

export class SSTablesManager {
    ssTablesLevels: Map<number, SSTable[]>;

    get sortedSSTablesLevels() {
        const sortedEntries = [...this.ssTablesLevels.entries()].sort((a, b) => a[0] - b[0]); // asc


        return sortedEntries;
    }

    constructor(private fileManager: FileManager, private idManager: IIdManager<any>) {
        this.ssTablesLevels = new Map<number, SSTable[]>()
        this.ssTablesLevels.set(0, [])
    }

    async deleteTable(ssTable: SSTable) {
        ssTable.delete();

        const ssTables = this.ssTablesLevels.get(ssTable.layerNumber);

        const index = ssTables.findIndex(t => t === ssTable);

        if (index !== -1) {
            ssTables.splice(index, 1); // mutate original array
        }
    }

    async init() {
        const ssTablesNamesNumbersInLayers = this.fileManager.getSSTablesNumbers();

        for (const [levelNumber, ssTablesNumbers] of ssTablesNamesNumbersInLayers) {
            for (const ssTableName of ssTablesNumbers) {
                let ssTable = new SSTable(this.fileManager.getDataFolderPath(), ssTableName.toString(), levelNumber);
                await ssTable.init()
                this.idManager.setMax(ssTable.metaData.maxId);
                this.registerSsTable(ssTable)
            }
        }
    }

    public  registerSsTable(ssTable: SSTable) {
        if (!this.ssTablesLevels.has(ssTable.layerNumber)) {
            this.ssTablesLevels.set(ssTable.layerNumber, [])
        }
        this.ssTablesLevels.get(ssTable.layerNumber).push(ssTable);
    }

    public async flushMemtableToSSTableInZeroLevel(memTable: MemTable<any, any>) {
        const nextSSTableNumber = this.fileManager.getNextSSTableNumberForZeroLevel();
        const newSSTable = new SSTable(this.fileManager.getDataFolderPath(), nextSSTableNumber.toString(), 0);
        await memTable.flush(newSSTable);
        this.fileManager.registerSSTableToZeroLevel(nextSSTableNumber);
        this.ssTablesLevels.get(0).push(newSSTable);
      //  await newSSTable.init();
    }

    drop() {
        for (const ssTablesLevel of this.ssTablesLevels) {
            ssTablesLevel[1].map(ssT => {
                ssT.delete();
                this.fileManager.drop();
            })
        }
    }
}