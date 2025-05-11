import {createLineReader, start} from "./file-line-by-line-reader";
import * as path from "path";

const TEST_DIR = path.join(__dirname, "test-data");

describe('file-line yb line reader.ts', () => {
    test('test', async () => {
        let filePath = path.join(TEST_DIR, '1.sst');
        const counter = await start(filePath)
      expect(counter).toEqual(8);


    })

    test('test with empty line at the end and then normal line and then empty', async () => {
        let filePath = path.join(TEST_DIR, '2.sst');
        const counter = await start(filePath)
        expect(counter).toEqual(10);
    })

    test('read line by line', async () => {
        let filePath = path.join(TEST_DIR, '3.sst');

        const lines = []
        const reader = createLineReader(filePath)
        const iterator = reader[Symbol.asyncIterator]();

        await new Promise<void>( (res) => {
               const interval = setInterval(async () => {
                   const lineRes = await iterator.next();
                   if (lineRes.done) {
                       clearInterval(interval)
                       res();
                   } else {
                       lines.push(lineRes.value);
                   }

               }, 100)

        })

        expect(lines.length).toEqual(3);
    })
})