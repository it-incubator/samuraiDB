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

}