import * as fs from "fs";
import * as readline from "readline";
import * as path from "path";
import {SSTable} from "../sstable/sstable";

type LineEntry = { key: string; value: string; reader: readline.Interface; fileIndex: number; done: boolean };

export class Compactor {
    constructor(private directory: string) {}

    async compactTables(sstables: SSTable[], levelFolder: string): Promise<SSTable[]> {
        const outDir = path.join(this.directory, levelFolder);
        if (!fs.existsSync(outDir)) fs.mkdirSync(outDir, { recursive: true });

        const readers: readline.Interface[] = [];
        const currentLines: (LineEntry | null)[] = [];

        const sstableSize = sstables.length > 0
            ? fs.statSync(sstables[0]["dataFilePath"]).size
            : 1024;
        const maxChunkBytes = sstableSize * 2;

        // --- Init readers and read first real line (skip metadata)
        for (let i = 0; i < sstables.length; i++) {
            await sstables[i].init();
            const stream = fs.createReadStream(sstables[i]['dataFilePath'], { encoding: "utf-8" });
            const reader = readline.createInterface({ input: stream });
            readers[i] = reader;

            const firstLine = await new Promise<string | null>(resolve => {
                let skipped = false;
                reader.on("line", (line) => {
                    if (!skipped) {
                        skipped = true;
                        return;
                    }
                    reader.removeAllListeners("line");
                    resolve(line);
                });
                reader.once("close", () => resolve(null));
            });

            if (firstLine) {
                const [key, value] = firstLine.split(/:(.+)/);
                currentLines[i] = { key, value, reader, fileIndex: i, done: false };
            } else {
                currentLines[i] = null;
            }
        }

        // --- Helpers
        const advance = async (i: number): Promise<void> => {
            const reader = readers[i];
            const line = await new Promise<string | null>(resolve => {
                reader.once("line", resolve);
                reader.once("close", () => resolve(null));
            });

            if (line) {
                const [key, value] = line.split(/:(.+)/);
                currentLines[i] = { key, value, reader, fileIndex: i, done: false };
            } else {
                currentLines[i] = null;
            }
        };

        const pickNext = (): { line: LineEntry; index: number } | null => {
            const active = currentLines
                .map((line, index) => ({ line, index }))
                .filter(e => e.line !== null) as { line: LineEntry, index: number }[];

            if (active.length === 0) return null;

            active.sort((a, b) => {
                const ak = Number(a.line.key);
                const bk = Number(b.line.key);
                return ak !== bk ? ak - bk : b.index - a.index; // newest wins
            });

            return active[0];
        };

        // --- Merge loop
        const writtenTables: SSTable[] = [];
        const seen = new Set<string>();
        let chunk: { key: string, value: string }[] = [];
        let chunkBytes = 0;
        let chunkIndex = 0;

        const writeChunk = async () => {
            if (chunk.length === 0) return;
            const keys = chunk.map(e => Number(e.key));
            const meta = {
                minId: String(Math.min(...keys)),
                maxId: String(Math.max(...keys)),
            };
            const fileName = `compacted-${chunkIndex++}`;
            const newTable = new SSTable(outDir, fileName);
            await newTable.write(meta, chunk);
            writtenTables.push(newTable);
            chunk = [];
            chunkBytes = 0;
        };

        while (true) {
            const next = pickNext();
            if (!next) break;

            const { line, index } = next;

            if (!seen.has(line.key)) {
                const entry = { key: line.key, value: line.value };
                const size = Buffer.byteLength(`${entry.key}:${entry.value}\n`, "utf-8");

                if (chunkBytes + size > maxChunkBytes) {
                    await writeChunk();
                }

                chunk.push(entry);
                chunkBytes += size;
                seen.add(line.key);
            }

            await advance(index);
        }

        await writeChunk();
        readers.forEach(r => r.close());
        sstables.forEach(t => t.delete());

        return writtenTables;
    }
}