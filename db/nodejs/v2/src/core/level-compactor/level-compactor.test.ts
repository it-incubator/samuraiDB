import {Compactor} from "./level-compactor";
import * as path from 'path';
import {FileManager} from "../samurai-db/file-manager/file-manager";
import {SSTablesManager} from "../samurai-db/ss-tables-manager";
import {IntegerIdStratagy} from "../samurai-db/samurai-db";



describe('SSTableFile.ts', () => {
    test('blabal', async () => {
        const compactor = new Compactor(path.join(__dirname, 'test-data', 'data'));

        const fileManager = new FileManager(path.join(__dirname, 'test-data', 'data'));
        const idStrategy = new IntegerIdStratagy();
        const ssmanager = new SSTablesManager(fileManager, idStrategy);
        await ssmanager.init()


        const resultTables = await compactor.compactTables(ssmanager.ssTables, 'level1');

    })
})