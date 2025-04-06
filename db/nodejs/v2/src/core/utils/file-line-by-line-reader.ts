import * as fs from "node:fs";
import * as path from "path";
import * as readline from "node:readline";





export const start = (filePath) => {
   return new Promise((res, rej) => {

       const stream = fs.createReadStream(filePath, { encoding: "utf-8" });
       const reader = readline.createInterface({ input: stream });

       let linesCounter = 0;
       reader.on('line', (line) => {
           linesCounter++;
       })
       reader.on('close', () => {
           res(linesCounter);
       })
   });
}

export const createLineReader = (filePath) => {
    const stream = fs.createReadStream(filePath, { encoding: "utf-8" });
    const reader = readline.createInterface({ input: stream });

    // // Add our own closed flag
    // // @ts-ignore
    // reader._customClosed = false;
    //
    // reader.on('line', (line) => {
    //
    //     reader.pause();
    // })
    // // Listen for close events on both reader and stream
    // reader.on('close', () => {
    //     lineHandler(null);
    // });
    //
    // return  {
    //     nextLine() {
    //
    //     }
    // };

    return reader;
}

export const readLine = (reader) => {
    return new Promise((resolve) => {
        // Check if we've already hit the end of the file using our custom flag
        if (reader._customClosed) {
            resolve(null);
            return;
        }

        // Flag to indicate if this promise has been resolved
        let isResolved = false;

        const onLine = (line) => {
            if (!isResolved) {
                isResolved = true;
                reader.removeListener('line', onLine);
                reader.removeListener('close', onClose);
                resolve(line);
            }
        };

        const onClose = () => {
            if (!isResolved) {
                isResolved = true;
                reader._customClosed = true;
                reader.removeListener('line', onLine);
                reader.removeListener('close', onClose);
                resolve(null);
            }
        };

        reader.once('line', onLine);
        reader.once('close', onClose);
    });
}