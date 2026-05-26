import {
  StorageType,
  KeyPairType,
  PreKeyPairType,
  SignedPreKeyPairType,
  SessionRecordType,
  Direction,
} from "@privacyresearch/libsignal-protocol-typescript";

export class SignalProtocolStore implements StorageType {
  async getIdentityKeyPair(): Promise<KeyPairType | undefined> {
    return undefined;
  }

  async getLocalRegistrationId(): Promise<number | undefined> {
    return undefined;
  }

  async isTrustedIdentity(
    identifier: string,
    identityKey: ArrayBuffer,
    direction: Direction,
  ): Promise<boolean> {
    return true;
  }

  async saveIdentity(
    encodedAddress: string,
    publicKey: ArrayBuffer,
    nonblockingApproval?: boolean,
  ): Promise<boolean> {
    return false;
  }

  async loadPreKey(
    encodedAddress: string | number,
  ): Promise<KeyPairType | undefined> {
    return undefined;
  }

  async storePreKey(
    keyId: number | string,
    keyPair: KeyPairType,
  ): Promise<void> {
    return;
  }

  async removePreKey(keyId: number | string): Promise<void> {
    return;
  }

  async storeSession(
    encodedAddress: string,
    record: SessionRecordType,
  ): Promise<void> {
    return;
  }

  async loadSession(
    encodedAddress: string,
  ): Promise<SessionRecordType | undefined> {
    return undefined;
  }

  async loadSignedPreKey(
    keyId: number | string,
  ): Promise<KeyPairType | undefined> {
    return undefined;
  }

  async storeSignedPreKey(
    keyId: number | string,
    keyPair: KeyPairType,
  ): Promise<void> {
    return;
  }

  async removeSignedPreKey(keyId: number | string): Promise<void> {
    return;
  }
}
