import {createServer} from 'net';
import {join} from 'path';
import {MemTable} from "../core/mem-table/mem-table";
import {RedBlackTree} from "../core/mem-table/IMemTableStructure/red-black-tree/red-black-tree";
import {SamuraiDb} from "../core/samurai-db/samurai-db";
import {FileManager} from "../core/samurai-db/file-manager/file-manager";
import {SSTablesManager} from "../core/samurai-db/ss-tables-manager";
import {IntegerIdStratagy} from "../core/samurai-db/integer-id-stratagy";

const dir = join(__dirname, '..', '..', 'db');

const redBlackTree = new RedBlackTree<number, any>();
const memTable = new MemTable<number, any>(redBlackTree);
const fileManager = new FileManager('data');
const idStrategy = new IntegerIdStratagy();
const ssTablesManager = new SSTablesManager(fileManager, idStrategy);
const db = new SamuraiDb<number, any>(memTable, fileManager, idStrategy, ssTablesManager);

(async () => {
  await db.init();
})();

const server = createServer(async (socket) => {
  console.log('Client connected');

  socket.on('data', async (data) => {
    let requestAction = JSON.parse(data.toString().split('\n')[0]);

    console.log('Received from client:', data.toString());

    switch (requestAction.type) {
      case 'SET': {
        const {id} = requestAction.payload;

        const entity = await db.set(id ?? null, { ...requestAction.payload });

        let response = {
          payload: entity,
          requestId: requestAction.requestId,
        };
        console.log(JSON.stringify(response));
        socket.write(JSON.stringify(response));
        break;
      }
      case 'GET': {
        const data = await db.get(requestAction.payload.id);
        let response = {
          payload: data,
          requestId: requestAction.requestId,
        };
        console.log('response: ', JSON.stringify(response));
        socket.write(JSON.stringify(response));
        break;
      }
      case 'DELETE': {
        const data = await db.delete(requestAction.payload.id);
        let response = {
          requestId: requestAction.requestId,
        };
        console.log('response: ', JSON.stringify(response));
        socket.write(JSON.stringify(response));
        break;
      }
      default: {
        console.error(`Unknown request type: ${requestAction.type}`);
        socket.write('Unknown request type');
        break;
      }
    }
  });

  socket.on('end', () => {
    console.log('Client disconnected');
  });

  socket.on('error', () => {
    console.log('Client error');
  });
});

server.listen(4001, () => {
  console.log('Server listening on port 4001');
});