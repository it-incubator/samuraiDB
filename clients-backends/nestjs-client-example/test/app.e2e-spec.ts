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

  it('/ (POST)', async () => {
    const createdIds = [];

    for (let i = 1; i < 41; i++) {
      if (i % 3 === 0) {
        const randomId =
          createdIds[Math.floor(Math.random() * createdIds.length)];

        await request(app.getHttpServer())
          .put('/samurais/' + randomId)
          .send({
            name: 'Updated Samurai ' + randomId,
            attackPower: -randomId,
            defensePower: 200 + i,
            health: Math.floor(Math.random() * 100),
            weapon: 'swore',
          } as SamuraiUpdateCreateDTO)
          .expect(200);
      }
      const response = await request(app.getHttpServer())
        .post('/samurais')
        .send({
          name: 'Samurai ' + i,
          attackPower: i,
          defensePower: 100 + i,
          health: Math.floor(Math.random() * 100),
          weapon: 'katana',
        } as SamuraiUpdateCreateDTO)
        .expect(201);

      createdIds.push(response.body.id);
    }
  });
});
