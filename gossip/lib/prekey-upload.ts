import type {
  PreKeyPairType,
  SignedPreKeyPairType,
} from "@privacyresearch/libsignal-protocol-typescript";
import {
  oneTimePrekeyStore,
  signedPrekeyStore,
} from "@/lib/local-stores";
import { Store } from "@/lib/local-stores/store";
import { apiFetch } from "@/lib/utils";

type UploadPrekeysOptions = {
  skipSessionCheck?: boolean;
};

type UploadPrekeysResponse = {
  uploaded: boolean;
  skippedUnauthenticated: boolean;
};

const uploadedPrekeyStore = new Store<boolean>("uploaded_prekeys", {
  dbName: "KeyVault",
});

async function hasAuthenticatedSession(): Promise<boolean> {
  const res = await fetch(process.env.NEXT_PUBLIC_API_URL + "auth/me", {
    credentials: "include",
  });

  return res.ok;
}

async function getLatestSignedPrekey(): Promise<SignedPreKeyPairType | null> {
  const entries = (await signedPrekeyStore.entries()) as [
    string,
    SignedPreKeyPairType,
  ][];

  return entries.reduce<SignedPreKeyPairType | null>((latest, [, prekey]) => {
    if (!latest || prekey.keyId > latest.keyId) {
      return prekey;
    }
    return latest;
  }, null);
}

async function getPendingPrekeys(): Promise<PreKeyPairType[]> {
  const entries = (await oneTimePrekeyStore.entries()) as [
    string,
    PreKeyPairType,
  ][];

  const pending: PreKeyPairType[] = [];
  for (const [, prekey] of entries) {
    const uploaded = await uploadedPrekeyStore.get(String(prekey.keyId));
    if (!uploaded) {
      pending.push(prekey);
    }
  }

  return [...pending].sort((a, b) => a.keyId - b.keyId);
}

export async function uploadPrekeysToServer(
  options: UploadPrekeysOptions = {},
): Promise<UploadPrekeysResponse> {
  if (!options.skipSessionCheck) {
    const authenticated = await hasAuthenticatedSession();
    if (!authenticated) {
      return { uploaded: false, skippedUnauthenticated: true };
    }
  }

  const signedPrekey = await getLatestSignedPrekey();
  const pendingPrekeys = await getPendingPrekeys();

  if (!signedPrekey || pendingPrekeys.length === 0) {
    return { uploaded: false, skippedUnauthenticated: false };
  }

  const { error } = await apiFetch<{ message: string }, { message: string }>(
    "keys/",
    {
      method: "POST",
      body: JSON.stringify({
        prekeyIds: pendingPrekeys.map((prekey) => prekey.keyId),
        prekeys: pendingPrekeys.map((prekey) =>
          new Uint8Array(prekey.keyPair.pubKey).toBase64(),
        ),
        signedPreKey: {
          id: signedPrekey.keyId,
          key: new Uint8Array(signedPrekey.keyPair.pubKey).toBase64(),
          signature: new Uint8Array(signedPrekey.signature).toBase64(),
        },
      }),
    },
  );

  if (error) {
    throw new Error(error.message ?? "Failed to upload prekeys");
  }

  await Promise.all(
    pendingPrekeys.map((prekey) =>
      uploadedPrekeyStore.set(String(prekey.keyId), true),
    ),
  );

  return { uploaded: true, skippedUnauthenticated: false };
}
