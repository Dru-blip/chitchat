import localforage from "localforage";

export type StoreConfig = {
  dbName?: string;
  description?: string;
};

export class Store<T> {
  private db: ReturnType<typeof localforage.createInstance>;

  constructor(storeName: string, config: StoreConfig = {}) {
    this.db = localforage.createInstance({
      name: config.dbName || "db",
      storeName,
      description: config.description,
      driver: [localforage.INDEXEDDB],
    });
  }

  async get(key: string): Promise<T | null> {
    return this.db.getItem<T>(key);
  }

  async set(key: string, value: T): Promise<T> {
    return this.db.setItem<T>(key, value);
  }

  async remove(key: string): Promise<void> {
    await this.db.removeItem(key);
  }

  async clear(): Promise<void> {
    await this.db.clear();
  }

  async keys(): Promise<string[]> {
    return this.db.keys();
  }

  async has(key: string): Promise<boolean> {
    const val = await this.db.getItem<T>(key);
    return val !== null;
  }

  async values(): Promise<T[]> {
    const vals: T[] = [];
    await this.db.iterate<T, void>((value) => {
      vals.push(value);
    });
    return vals;
  }

  async entries(): Promise<[string, T][]> {
    const pairs: [string, T][] = [];
    await this.db.iterate<T, void>((value, key) => {
      pairs.push([key, value]);
    });
    return pairs;
  }

  async count(): Promise<number> {
    return this.db.length();
  }
}
