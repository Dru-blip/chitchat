import { betterAuth } from 'better-auth';
import { prismaAdapter } from 'better-auth/adapters/prisma';
import { PrismaClient } from './prisma/generated/prisma/client';

export const auth = betterAuth({
  database: prismaAdapter({} as PrismaClient, {
    provider: 'sqlite',
  }),
});
