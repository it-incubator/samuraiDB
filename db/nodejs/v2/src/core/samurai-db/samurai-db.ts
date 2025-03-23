import {ISamuraiDB} from "./i-samurai-db";
import {IMemTable} from "../mem-table/i-mem-table";
import {FileManager} from "./file-manager/file-manager";
import {SSTablesManager} from "./ss-tables-manager";


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


export class SamuraiDb<TKey, TValue> implements ISamuraiDB<TKey, TValue> {
    constructor(private memTable: IMemTable<TKey, TValue>, private fileManager: FileManager, private idManager: IIdManager<TKey>, private sSTablesManager: SSTablesManager) {

    }

    public async init() {
        return this.sSTablesManager.init();
    }

    public async set(key: TKey | null, value: TValue): Promise<TValue> {
        const correctedKey = key === null ? this.idManager.getNext() : key;
        this.memTable.set(correctedKey, {...value, id: correctedKey});
        const needFlushToSSTable = this.memTable.isFull();
        if (needFlushToSSTable) {
            await this.sSTablesManager.flushMemtableToSSTable(this.memTable)
        }
        return {...value, id: correctedKey};
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