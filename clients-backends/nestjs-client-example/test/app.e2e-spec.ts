import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from './../src/app.module';
import { SamuraiUpdateCreateDTO } from '../src/app.controller';

describe('AppController (e2e)', () => {
  let app: INestApplication;

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    await app.init();
  });

  it('/samurais (POST & PUT + GET verification)', async () => {
    const createdIds: string[] = [];

    for (let i = 1; i <= 40; i++) {
      if (i % 3 === 0 && createdIds.length > 0) {
        // Выбираем случайный ранее созданный ID
        const randomId =
          createdIds[Math.floor(Math.random() * createdIds.length)];

        const updatedSamurai: SamuraiUpdateCreateDTO = {
          name: 'Updated Samurai ' + randomId,
          attackPower: -parseInt(randomId, 10) || 0,
          defensePower: 200 + i,
          health: Math.floor(Math.random() * 100),
          weapon: 'swore',
        };

        // Обновляем
        await request(app.getHttpServer())
          .put(`/samurais/${randomId}`)
          .send(updatedSamurai)
          .expect(200);

        // Проверяем, что обновление применилось
        const getResponse = await request(app.getHttpServer())
          .get(`/samurais/${randomId}`)
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
      createdIds.push(id);

      // Проверяем, что POST-самурай был сохранён
      const getResponse = await request(app.getHttpServer())
        .get(`/samurais/${id}`)
        .expect(200);

      expect(getResponse.body).toMatchObject(newSamurai);
    }
  });
});
