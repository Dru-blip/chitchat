import type { KeyPairType } from "@privacyresearch/libsignal-protocol-typescript";
import { CryptoKeyStore } from "./keys";
import { Store } from "./store";

export const identityStore = new CryptoKeyStore("identity_keys");
export const signedPrekeyStore = new CryptoKeyStore("signed_prekeys");
export const oneTimePrekeyStore = new CryptoKeyStore("one_time_prekeys");

const counterStore = new Store<number>("id_counters");

export async function nextPersistentId(
  scope: string = "global",
): Promise<number> {
  const current = (await counterStore.get(scope)) ?? 0;
  const next = current + 1;
  await counterStore.set(scope, next);
  return next;
}

export async function getSerializedIdentityKeys(): Promise<{
  pubKey: string;
  privKey: string;
}> {
  const identityKeyPair = (await identityStore.get("identity")!) as KeyPairType;
  return {
    pubKey: new Uint8Array(identityKeyPair.pubKey).toBase64(),
    privKey: new Uint8Array(identityKeyPair.privKey).toBase64(),
  };
}
