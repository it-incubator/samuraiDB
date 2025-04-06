import {IMemTable} from "./i-mem-table";
import {IMemTableStructure} from "./IMemTableStructure/i-mem-table-structure";
import {MetaDataType, SSTable} from "../sstable/sstable";

export class MemTable<TKey, TValue> implements IMemTable<TKey, TValue> {
    private structure: IMemTableStructure<TKey, TValue>;

    constructor(structure: IMemTableStructure<TKey, TValue>) {
        this.structure = structure;
    }

    public set(key: TKey, value: TValue): void {
        this.structure.insert(key, value);
    }

    public delete(key: TKey): void {
        this.structure.delete(key);
    }

    public get(key: TKey): TValue | undefined {
        const found = this.structure.find(key);
        return found || undefined;
    }

    public isFull(): boolean {
         return this.structure.getCount() > 5;
    }

    public async flush(ssTable: SSTable): Promise<void> {
        const data = this.structure.getSortedArray()
            .map(item => ({
                key: item.key as string, // todo: need to make stringify??
                value: JSON.stringify(item.value)
            }));


        await ssTable.write(data);

        this.structure.clear();
    }
}