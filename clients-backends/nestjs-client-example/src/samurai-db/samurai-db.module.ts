import { Module } from '@nestjs/common';
import { SamuraiDBDriver } from './samurai-db-driver';
import { ConfigurableModuleClass } from './database.module-definition';
import { SamuraiDBConnect } from './infrastructure/samurai-db-connect';
import { ConnectionService } from './infrastructure/connection.service';

@Module({
  imports: [],
  controllers: [],
  providers: [
    SamuraiDBDriver,
    { provide: ConnectionService, useClass: SamuraiDBConnect },
  ],
  exports: [SamuraiDBDriver],
})
export class SamuraiDbModule extends ConfigurableModuleClass {}
