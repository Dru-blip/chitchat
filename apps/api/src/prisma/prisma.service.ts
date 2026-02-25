import {
  Injectable,
  OnModuleInit,
  OnModuleDestroy,
  Logger,
} from '@nestjs/common';
import { PrismaClient } from './generated/prisma/client';
import { PrismaBetterSqlite3 } from '@prisma/adapter-better-sqlite3';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class PrismaService
  extends PrismaClient
  implements OnModuleInit, OnModuleDestroy
{
  private readonly logger = new Logger(PrismaService.name);

  constructor(readonly configService: ConfigService) {
    super({
      adapter: new PrismaBetterSqlite3({
        url: configService.get<string>('DATABASE_URL')!,
      }),
      log:
        process.env.NODE_ENV === 'development'
          ? ['query', 'info', 'warn', 'error']
          : ['error'],
    });
  }

  async onModuleInit() {
    await this.$connect();
    this.logger.log('Prisma connection established');
  }

  async onModuleDestroy() {
    await this.$disconnect();
    this.logger.log('Prisma connection closed');
  }
}
