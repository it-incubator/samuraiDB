import {SSTable} from "../sstable/sstable";
import {FileManager} from "./file-manager/file-manager";
import {IMemTable} from "../mem-table/i-mem-table";
import {IIdManager} from "./samurai-db";

export class SSTablesManager {
    ssTables: SSTable[] = []

    constructor(private fileManager: FileManager, private idManager: IIdManager<any>) {
    }

    async init() {
        const ssTablesNamesNumbers = this.fileManager.getSSTablesNumbers();

        for (const ssTableName of ssTablesNamesNumbers) {
            let ssTable = new SSTable(this.fileManager.getDataFolderPath(), ssTableName.toString());
            await ssTable.init()
            this.idManager.setMax(ssTable.metaData.maxId);
            this.ssTables.push(ssTable)
        }
    }

    public async flushMemtableToSSTable(memTable: IMemTable<any, any>) {
        const nextSSTableNumber = this.fileManager.getNextSSTableNumber();
        const newSSTable = new SSTable(this.fileManager.getDataFolderPath(), nextSSTableNumber.toString());
        await memTable.flush(newSSTable);
        this.fileManager.registerSSTable(nextSSTableNumber);
        this.ssTables.push(newSSTable);
    }
}