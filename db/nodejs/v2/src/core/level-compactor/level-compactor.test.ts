import {Compactor} from "./level-compactor";
import * as path from 'path';
import {FileManager} from "../samurai-db/file-manager/file-manager";
import {SSTablesManager} from "../samurai-db/ss-tables-manager";
import * as fs from "node:fs";
import {IntegerIdStratagy} from "../samurai-db/integer-id-stratagy";

function copyRecursiveSync(src: string, dest: string) {
    const stats = fs.statSync(src);
    if (stats.isDirectory()) {
        if (!fs.existsSync(dest)) {
            fs.mkdirSync(dest, { recursive: true });
        }
        const entries = fs.readdirSync(src);
        entries.forEach(entry => {
            const srcPath = path.join(src, entry);
            const destPath = path.join(dest, entry);
            copyRecursiveSync(srcPath, destPath);
        });
    } else {
        fs.copyFileSync(src, dest);
    }
}


describe('SSTableFile.ts', () => {

    const TEST_DATA_DIR = path.join(__dirname, 'test-data');
    const SAMPLES_DIR = path.join(TEST_DATA_DIR, 'samples');
    const DATA_DIR = path.join(TEST_DATA_DIR, 'data');

    beforeEach(() => {
        // Очищаем директорию с данными перед каждым тестом
        if (fs.existsSync(DATA_DIR)) {
            fs.rmSync(DATA_DIR, { recursive: true });
        }
        fs.mkdirSync(DATA_DIR, { recursive: true });

        // Рекурсивно копируем все из samples в data
        copyRecursiveSync(SAMPLES_DIR, DATA_DIR);
    });


    test('simple exanple', async () => {


        const fileManager = new FileManager(path.join(__dirname, 'test-data', 'data'));
        fileManager.init();
        const idStrategy = new IntegerIdStratagy();
        const ssmanager = new SSTablesManager(fileManager, idStrategy);
        await ssmanager.init()

        const compactor = new Compactor(fileManager, ssmanager);


        const resultTables = await compactor.compactTables(0);

    })


    test.skip('exanple ehtn first level already exists', async () => {

        const fileManager = new FileManager(path.join(__dirname, 'test-data', 'data'));
        fileManager.init();
        const idStrategy = new IntegerIdStratagy();
        const ssmanager = new SSTablesManager(fileManager, idStrategy);
        await ssmanager.init()

        const compactor = new Compactor(fileManager, ssmanager);


       // const sstables = ssmanager.ssTables from zero level + sstables from 1 level
       //  const resultTables = await compactor.compactTables(sstables, 'level1');

    })
})