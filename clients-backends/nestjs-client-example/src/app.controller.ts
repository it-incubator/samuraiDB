import {
  Body,
  Controller,
  Delete,
  Get,
  HttpCode,
  HttpStatus,
  NotFoundException,
  Param,
  Post,
  Put,
} from '@nestjs/common';
import { AppService } from './app.service';
import { SamuraiDBDriver } from './samurai-db/samurai-db-driver';
import { ApiBody, ApiOperation, ApiParam, ApiProperty } from '@nestjs/swagger';

export class SamuraiUpdateCreateDTO {
  @ApiProperty({
    example: 'Miyamoto Musashi',
    description: 'Name of the Samurai',
  })
  name: string;

  @ApiProperty({ example: 100, description: 'Health points of the Samurai' })
  health: number;

  @ApiProperty({ example: 50, description: 'Attack power of the Samurai' })
  attackPower: number;

  @ApiProperty({ example: 30, description: 'Defense power of the Samurai' })
  defensePower: number;

  @ApiProperty({ example: 'Katana', description: 'Weapon used by the Samurai' })
  weapon: string;
}

export class SamuraiEntity extends SamuraiUpdateCreateDTO {
  @ApiProperty({
    example: '123',
    description: 'Unique identifier for the Samurai',
  })
  id: string;
}

@Controller('samurais')
export class AppController {
  constructor(
    private readonly samuraiDBDriver: SamuraiDBDriver<SamuraiEntity>,
  ) {}

  @Delete('db')
  @HttpCode(HttpStatus.NO_CONTENT)
  async drop(): Promise<void> {
    return this.samuraiDBDriver.drop();
  }

  @ApiOperation({ summary: 'Run Compaction' })
  @ApiParam({ name: 'id', required: true, description: 'Samurai ID' })
  @Post('compaction')
  async runCompaction(): Promise<void> {
    return this.samuraiDBDriver.runCompaction();
  }

  @ApiOperation({ summary: 'Get Samurai by ID' })
  @ApiParam({ name: 'id', required: true, description: 'Samurai ID' })
  @Get(':id')
  async getById(@Param('id') id: string): Promise<SamuraiEntity> {
    const item = await this.samuraiDBDriver.getById(id);
    if (!item) {
      console.log('❤️ not found ' + id);
      throw new NotFoundException();
    }

    return item;
  }

  @ApiOperation({ summary: 'Create new Samurai' })
  @ApiBody({ type: SamuraiUpdateCreateDTO })
  @Post()
  async create(@Body() dto: SamuraiUpdateCreateDTO): Promise<SamuraiEntity> {
    const result = await this.samuraiDBDriver.set<SamuraiEntity>(dto);
    return result;
  }

  @ApiOperation({ summary: 'Update Samurai by ID' })
  @ApiParam({ name: 'id', required: true, description: 'Samurai ID' })
  @ApiBody({ type: SamuraiUpdateCreateDTO })
  @Put(':id')
  async update(
    @Param('id') id: string,
    @Body() dto: SamuraiUpdateCreateDTO,
  ): Promise<SamuraiEntity> {
    return this.samuraiDBDriver.updateById(id, dto);
  }

  @ApiOperation({ summary: 'Delete Samurai by ID' })
  @ApiParam({ name: 'id', required: true, description: 'Samurai ID' })
  @Delete(':id')
  async delete(@Param('id') id: string): Promise<void> {
    return this.samuraiDBDriver.deleteById(id);
  }
}
