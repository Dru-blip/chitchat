import { Store } from "./store";

import * as libsignal from "@privacyresearch/libsignal-protocol-typescript";

export class CryptoKeyStore<
  T extends libsignal.PreKeyPairType | libsignal.KeyPairType,
> extends Store<T> {
  constructor(storeName: string, dbName = "KeyVault") {
    super(storeName, {
      dbName,
    });
  }

  async savePair(key: string, pair: T) {
    return this.set(key, pair);
  }

  async loadPair(key: string): Promise<T | null> {
    return this.get(key);
  }
}
