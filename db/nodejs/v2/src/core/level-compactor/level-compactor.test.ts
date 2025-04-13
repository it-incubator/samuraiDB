import {Compactor} from "./level-compactor";
import * as path from 'path';
import {FileManager} from "../samurai-db/file-manager/file-manager";
import {SSTablesManager} from "../samurai-db/ss-tables-manager";
import * as fs from "node:fs";
import {IntegerIdStratagy} from "../samurai-db/integer-id-stratagy";



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

        // Копируем тестовые файлы из samples в data
        const sampleFiles = fs.readdirSync(SAMPLES_DIR);
        sampleFiles.forEach(file => {
            fs.copyFileSync(
                path.join(SAMPLES_DIR, file),
                path.join(DATA_DIR, file)
            );
        });
    });


    test('simple exanple', async () => {
        const compactor = new Compactor(path.join(__dirname, 'test-data', 'data'));

        const fileManager = new FileManager(path.join(__dirname, 'test-data', 'data'));
        fileManager.init();
        const idStrategy = new IntegerIdStratagy();
        const ssmanager = new SSTablesManager(fileManager, idStrategy);
        await ssmanager.init()


        const resultTables = await compactor.compactTables(ssmanager.ssTablesLevels, 'level1');

    })


    test('exanple ehtn first level already exists', async () => {
        const compactor = new Compactor(path.join(__dirname, 'test-data', 'data'));

        const fileManager = new FileManager(path.join(__dirname, 'test-data', 'data'));
        fileManager.init();
        const idStrategy = new IntegerIdStratagy();
        const ssmanager = new SSTablesManager(fileManager, idStrategy);
        await ssmanager.init()


       // const sstables = ssmanager.ssTables from zero level + sstables from 1 level
       //  const resultTables = await compactor.compactTables(sstables, 'level1');

    })
})