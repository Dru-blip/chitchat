import { betterAuth } from 'better-auth';
import { prismaAdapter } from 'better-auth/adapters/prisma';
import { PrismaClient } from './prisma/generated/prisma/client';
import { emailOTP } from 'better-auth/plugins';
import { emailHarmony } from 'better-auth-harmony';

export const auth = betterAuth({
  database: prismaAdapter({} as PrismaClient, {
    provider: 'sqlite',
  }),
  emailAndPassword: {
    enabled: true,
  },
  plugins: [
    emailHarmony(),
    emailOTP({
      sendVerificationOTP: async () => {},
      storeOTP: 'hashed',
    }),
  ],
});
