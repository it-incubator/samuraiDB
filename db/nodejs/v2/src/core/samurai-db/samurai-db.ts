import {ISamuraiDB} from "./i-samurai-db";
import {FileManager} from "./file-manager/file-manager";
import {SSTablesManager} from "./ss-tables-manager";
import {MemTable} from "../mem-table/mem-table";
import {IIdManager} from "./i-id-manager";


export class SamuraiDb<TKey, TValue> implements ISamuraiDB<TKey, TValue> {
    constructor(private memTable: MemTable<TKey, TValue>, private fileManager: FileManager, private idManager: IIdManager<TKey>, private sSTablesManager: SSTablesManager) {

    }

    public async init() {
        this.fileManager.init();
        return this.sSTablesManager.init();
    }

    public async set(key: TKey | null, value: TValue): Promise<TValue> {
        const correctedKey = key === null ? this.idManager.getNext() : key;
        this.memTable.set(correctedKey, {...value, id: correctedKey});
        const needFlushToSSTable = this.memTable.isFull();
        if (needFlushToSSTable) {
            await this.sSTablesManager.flushMemtableToSSTableInZeroLevel(this.memTable)
        }
        return {...value, id: correctedKey};
    }

    public async get(key: TKey): Promise<TValue> | null {
        let foundItem = this.memTable.get(key);
        if (foundItem) {
            console.log("foundItem inside memtable: ", foundItem)
            return foundItem;
        }

        for (const [level, ssTables] of this.sSTablesManager.sortedSSTablesLevels) {
                // for 0 level we need start from the latest/newest table to the oldest
               // for other level we don't need make revers so in futture make optimization'
                for (const ssTable of [...ssTables].reverse()) {
                    foundItem = await ssTable.read(key.toString()) as TValue;
                    if (foundItem) {
                        return foundItem;
                    }
                }
        }

        return null;
    }

    public delete(key: TKey): void {
        this.memTable.delete(key);
    }


}