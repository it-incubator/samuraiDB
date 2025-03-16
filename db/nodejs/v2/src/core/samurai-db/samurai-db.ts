import {ISamuraiDB} from "./i-samurai-db";
import {IMemTable} from "../mem-table/i-mem-table";
import {SSTable} from "../sstable/sstable";
import * as fs from "fs";
import {FileManager} from "./file-manager/file-manager";


export interface IIdManager<TKey> {
    getNext(): TKey
    setMax(key: TKey):void
}

export class IntegerIdStratagy implements IIdManager<number> {
    maxId = 0;

    getNext(): number {
        ++this.maxId;
        return this.maxId;
    }
    setMax(maxId: number) {
        if (this.maxId < maxId) {
            this.maxId = maxId;
        }
    }

}


export class SSTablesManager {
    ssTables: SSTable[] = []

    constructor(private fileManager: FileManager, private idManager: IIdManager<any>) {
    }

    async restore() {
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

export class SamuraiDb<TKey, TValue> implements ISamuraiDB<TKey, TValue> {
    constructor(private memTable: IMemTable<TKey, TValue>, private fileManager: FileManager, private idManager: IIdManager<TKey>, private sSTablesManager: SSTablesManager) {

    }

    public async init() {
        return this.sSTablesManager.restore();
    }

    public async set(key: TKey | null, value: TValue): Promise<void> {
        const correctedKey = key === null ? this.idManager.getNext() : key;
        this.memTable.set(correctedKey, {...value, correctedKey});
        const needFlushToSSTable = this.memTable.isFull();
        if (needFlushToSSTable) {
            await this.sSTablesManager.flushMemtableToSSTable(this.memTable)
        }
    }

    public async get(key: TKey): Promise<TValue> | null {
        let foundItem = this.memTable.get(key);
        if (foundItem) {
            console.log("foundItem: ", foundItem)
            return foundItem;
        }

        for (const ssTable of [...this.sSTablesManager.ssTables].reverse()) {
            foundItem = await ssTable.read(key.toString()) as TValue;
            if (foundItem) {
                return foundItem;
            }
        }

        return null;
    }

    public delete(key: TKey): void {
        this.memTable.delete(key);
    }


}