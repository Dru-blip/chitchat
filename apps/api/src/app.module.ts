import { Module } from '@nestjs/common';
import { PrismaModule } from './prisma/prisma.module';
import { AuthModule } from '@thallesp/nestjs-better-auth';
import { PrismaService } from './prisma/prisma.service';
import { ConfigModule } from '@nestjs/config';
import { betterAuth } from 'better-auth';
import { prismaAdapter } from 'better-auth/adapters/prisma';
import { emailOTP } from 'better-auth/plugins';
import { emailHarmony } from 'better-auth-harmony';
import { MailModule } from './mail/mail.module';
import { MailService } from './mail/mail.service';
import { EventEmitter2, EventEmitterModule } from '@nestjs/event-emitter';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    EventEmitterModule.forRoot(),
    PrismaModule,
    MailModule,
    AuthModule.forRootAsync({
      imports: [PrismaModule, MailModule],
      inject: [PrismaService, MailService],
      useFactory: (prisma: PrismaService, eventEmitter: EventEmitter2) => ({
        auth: betterAuth({
          database: prismaAdapter(prisma, {
            provider: 'sqlite',
          }),
          emailAndPassword: {
            enabled: true,
          },
          plugins: [
            emailHarmony(),
            emailOTP({
              sendVerificationOTP: async ({ otp, email, type }) => {
                if (type === 'sign-in') {
                  eventEmitter.emit('otp.send', { to: email, otp });
                }
              },
              storeOTP: 'hashed',
            }),
          ],
        }),
      }),
    }),
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
