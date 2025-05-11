import {IIdManager} from "./i-id-manager";

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

    reset() {
        this.maxId = 0;
    }

    convertToIdFormat(key: any): number {
        if (typeof key === 'number' && Number.isInteger(key)) {
            return key;
        }

        const parsed = Number(key);
        if (Number.isInteger(parsed)) {
            return parsed;
        }

        throw new Error(`Invalid key format: ${key}`);
    }

}