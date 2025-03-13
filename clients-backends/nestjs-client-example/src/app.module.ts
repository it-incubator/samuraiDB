import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { SamuraiDbModule } from './samurai-db/samurai-db.module';
import { ConfigModule } from '@nestjs/config';
import { RetryStrategy } from './samurai-db/interfaces/module-options';

@Module({
  imports: [
    ConfigModule.forRoot({ envFilePath: '.env.local' }),
    SamuraiDbModule.register({
      host: 'localhost',
      port: 4001,
      maxRetries: 100,
      interval: 1000,
      retryStrategy: RetryStrategy.FIXED,
    }),
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
