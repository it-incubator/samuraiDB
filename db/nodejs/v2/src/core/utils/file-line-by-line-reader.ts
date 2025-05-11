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

