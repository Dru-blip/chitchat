import {
  Injectable,
  InternalServerErrorException,
  Logger,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { OnEvent } from '@nestjs/event-emitter';
import nodemailer from 'nodemailer';

@Injectable()
export class MailService {
  private readonly logger = new Logger(MailService.name);
  private readonly transporter: nodemailer.Transporter;
  private readonly senderEmail: string;

  constructor(private readonly configService: ConfigService) {
    const host = this.configService.get<string>('SMTP_HOST')!;
    const port = this.configService.get<number>('SMTP_PORT')!;
    const secure = this.configService.get<boolean>('SMTP_SECURE')!;
    const user = this.configService.get<string>('SMTP_USER')!;
    const pass = this.configService.get<string>('SMTP_PASS')!;
    this.senderEmail = this.configService.get<string>('SMTP_FROM')!;

    this.transporter = nodemailer.createTransport({
      host,
      port,
      secure,
      auth: { user, pass },
    });
  }

  @OnEvent('otp.send', { async: true })
  async sendOtpEmail({ to, otp }: { to: string; otp: string }): Promise<void> {
    try {
      await this.transporter.sendMail({
        from: this.senderEmail,
        to,
        subject: 'OTP',
        text: `Your OTP is ${otp}. It expires in 5 minutes.`,
        html: `<p>Your OTP is <b>${otp}</b>.</p><p>It expires in 5 minutes.</p>`,
      });
    } catch (error) {
      this.logger.error(`Failed to send OTP email to ${to}`, error);
      throw new InternalServerErrorException('Unable to send OTP email');
    }
  }
}
