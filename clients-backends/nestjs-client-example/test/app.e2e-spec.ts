import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from './../src/app.module';
import { SamuraiUpdateCreateDTO } from '../src/app.controller';

jest.setTimeout(30000);

describe('AppController (e2e)', () => {
  let app: INestApplication;

  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    await app.init();

    await request(app.getHttpServer()).delete(`/samurais/db`).expect(204);
  });

  const createdSamurais: { id: string; samurai: SamuraiUpdateCreateDTO }[] = [];

  const deletedIds = [];

  it('/samurais (POST & PUT + DELETE + GET verification)', async () => {
    for (let i = 1; i <= 40; i++) {
      // Каждые 3 шага — обновление рандомного самурая
      if (i % 3 === 0 && createdSamurais.length > 0) {
        const randomEntry =
          createdSamurais[Math.floor(Math.random() * createdSamurais.length)];

        const updatedSamurai: SamuraiUpdateCreateDTO = {
          name: 'Updated Samurai ' + randomEntry.id,
          attackPower: -parseInt(randomEntry.id, 10) || 0,
          defensePower: 200 + i,
          health: Math.floor(Math.random() * 100),
          weapon: 'swore',
        };

        await request(app.getHttpServer())
          .put(`/samurais/${randomEntry.id}`)
          .send(updatedSamurai)
          .expect(200);

        randomEntry.samurai = updatedSamurai;

        const getResponse = await request(app.getHttpServer())
          .get(`/samurais/${randomEntry.id}`)
          .expect(200);

        expect(getResponse.body).toMatchObject(updatedSamurai);
      }

      // создаём нового самурая
      const newSamurai: SamuraiUpdateCreateDTO = {
        name: 'Samurai ' + i,
        attackPower: i,
        defensePower: 100 + i,
        health: Math.floor(Math.random() * 100),
        weapon: 'katana',
      };

      const postResponse = await request(app.getHttpServer())
        .post('/samurais')
        .send(newSamurai)
        .expect(201);

      const id = postResponse.body.id;
      expect(id).toBeDefined();

      createdSamurais.push({ id, samurai: newSamurai });

      // Каждые 6 шагов — удаление рандомного самурая
      if (i % 6 === 0 && createdSamurais.length > 0) {
        const index = Math.floor(Math.random() * createdSamurais.length);
        const toDelete = createdSamurais[index];

        await request(app.getHttpServer())
          .delete(`/samurais/${toDelete.id}`)
          .expect(200);

        await request(app.getHttpServer())
          .get(`/samurais/${toDelete.id}`)
          .expect(404);

        deletedIds.push(toDelete.id);
        createdSamurais.splice(index, 1); // удалить из списка, чтобы не использовать дальше
      } else {
        const getResponse = await request(app.getHttpServer())
          .get(`/samurais/${id}`)
          .expect(200);

        expect(getResponse.body).toMatchObject(newSamurai);
      }
    }
  });

  it('/samurais (GET after compaction: deleted still 404, others ok)', async () => {
    await request(app.getHttpServer()).post(`/samurais/compaction`).expect(201);

    for (const entry of createdSamurais) {
      const getResponse = await request(app.getHttpServer())
        .get(`/samurais/${entry.id}`)
        .expect(200);

      expect(getResponse.body).toMatchObject(entry.samurai);
    }

    for (const id of deletedIds) {
      await request(app.getHttpServer()).get(`/samurais/${id}`).expect(404);
    }
  });
});
