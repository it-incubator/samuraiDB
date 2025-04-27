import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from './../src/app.module';
import { SamuraiUpdateCreateDTO } from '../src/app.controller';

jest.setTimeout(30000);

describe('AppController (e2e)', () => {
  let app: INestApplication;

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    await app.init();
  });

  const createdSamurais: { id: string; samurai: SamuraiUpdateCreateDTO }[] = [];

  it('/samurais (POST & PUT + GET verification)', async () => {
    for (let i = 1; i <= 40; i++) {
      if (i % 3 === 0 && createdSamurais.length > 0) {
        // Выбираем случайный ранее созданный объект чтобы обновить
        const randomEntry =
          createdSamurais[Math.floor(Math.random() * createdSamurais.length)];

        const updatedSamurai: SamuraiUpdateCreateDTO = {
          name: 'Updated Samurai ' + randomEntry.id,
          attackPower: -parseInt(randomEntry.id, 10) || 0,
          defensePower: 200 + i,
          health: Math.floor(Math.random() * 100),
          weapon: 'swore',
        };

        // Обновляем
        await request(app.getHttpServer())
          .put(`/samurais/${randomEntry.id}`)
          .send(updatedSamurai)
          .expect(200);

        // Обновляем локально в памяти
        randomEntry.samurai = updatedSamurai;

        // Проверяем, что обновление применилось
        const getResponse = await request(app.getHttpServer())
          .get(`/samurais/${randomEntry.id}`)
          .expect(200);

        expect(getResponse.body).toMatchObject(updatedSamurai);
      }

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

      // Проверяем, что POST-самурай был сохранён
      const getResponse = await request(app.getHttpServer())
        .get(`/samurais/${id}`)
        .expect(200);

      expect(getResponse.body).toMatchObject(newSamurai);
    }
  });

  it('/samurais (GET all and verify) after compaction', async () => {
    await request(app.getHttpServer()).post(`/samurais/compaction`).expect(201);

    for (const entry of createdSamurais) {
      const getResponse = await request(app.getHttpServer())
        .get(`/samurais/${entry.id}`)
        .expect(200);

      expect(getResponse.body).toMatchObject(entry.samurai);
    }
  });
});
