import { Module } from '@nestjs/common';
import { PrismaModule } from './prisma/prisma.module';
import { AuthModule } from '@thallesp/nestjs-better-auth';
import { PrismaService } from './prisma/prisma.service';
import { ConfigModule } from '@nestjs/config';
import { betterAuth } from 'better-auth';
import { prismaAdapter } from 'better-auth/adapters/prisma';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    PrismaModule,
    AuthModule.forRootAsync({
      imports: [PrismaModule],
      inject: [PrismaService],
      useFactory: (prisma: PrismaService) => ({
        auth: betterAuth({
          database: prismaAdapter(prisma, {
            provider: 'sqlite',
          }),
          emailAndPassword: {
            enabled: true,
          },
        }),
      }),
    }),
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
