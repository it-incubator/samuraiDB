import SamuraiDB from './core/samuraidb';
import { FileAdapter } from './core/file.adapter';

const fileAdapter = new FileAdapter(__dirname + '/samuraidb.txt');
const samuraiDB = new SamuraiDB(fileAdapter);


async function seedDB() {
  for (let i = 0; i < 100; i++) {
    const id = String(i);
    await samuraiDB.set(id, {
      id: id,
      title: 'Some title ' + id,
      date: new Date(2024, 11, 11),
    });
  }
}

async function updateSomeData() {
  for (let i = 0; i < 100; i++) {
    if (i % 2 === 0) {
      const id = String(i);
      await samuraiDB.set(id, {
        id,
        title: 'Updated title ' + id,
        date: new Date(2024, 11, 11),
      });
    }
  }
}

(async () => {
  await seedDB();
  // updateSomeData();
  const foundItem = await samuraiDB.get('98');
  console.log(foundItem);
})();
