import {
  KeyHelper,
  type KeyPairType,
} from "@privacyresearch/libsignal-protocol-typescript";
import {
  identityStore,
  nextPersistentId,
  oneTimePrekeyStore,
  signedPrekeyStore,
} from "./local-stores";

export const initializeRegistrationId = async (): Promise<void> => {
  const existing = await identityStore.get("registrationId");
  if (existing) return;

  const registrationId = KeyHelper.generateRegistrationId();
  await identityStore.set("registrationId", registrationId as any);
};

export const initializeIdentityKey = async (): Promise<void> => {
  const existingIdentity = await identityStore.loadPair("identity");
  if (existingIdentity) {
    return;
  }

  const keypair = await KeyHelper.generateIdentityKeyPair();
  await identityStore.savePair("identity", keypair);
};

export const initializeSignedPreKey = async (
  threshold: number = 1,
): Promise<void> => {
  const identityKeyPair = await identityStore.loadPair("identity");
  if (!identityKeyPair) {
    throw new Error("Identity key not found ");
  }

  const existingCount = await signedPrekeyStore.count();

  if (existingCount >= threshold) {
    return;
  }

  const signedKeyId = await nextPersistentId("signed_prekeys");

  const keypair = await KeyHelper.generateSignedPreKey(
    identityKeyPair as KeyPairType,
    signedKeyId,
  );
  await signedPrekeyStore.savePair(`signed_prekey_${signedKeyId}`, keypair);
};

export const initializePreKeys = async (
  threshold = 10,
  target = 100,
): Promise<void> => {
  const existingCount = await oneTimePrekeyStore.count();

  if (existingCount >= threshold) {
    return;
  }

  const needed = target - existingCount;
  for (let i = 0; i < needed; i++) {
    const id = await nextPersistentId("prekeys");
    const keypair = await KeyHelper.generatePreKey(id);
    await oneTimePrekeyStore.savePair(`prekey_${id}`, keypair);
  }
};

export const initializeKeys = async (): Promise<void> => {
  await initializeRegistrationId();
  await initializeIdentityKey();
  await initializeSignedPreKey();
  await initializePreKeys();
};
