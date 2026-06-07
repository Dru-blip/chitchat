import { Message, ConversationMeta } from "@/types";
import localforage from "localforage";

export class MessageRepository {
  private store = localforage.createInstance({
    name: "db",
    storeName: "messages",
  });

  private metaStore = localforage.createInstance({
    name: "db",
    storeName: "conversationMeta",
  });

  async appendMessage(conversationId: string, message: Message): Promise<void> {
    const key = `msg:${conversationId}:${message.sent_at}:${message.id}`;
    await this.store.setItem(key, message);

    await this.updateConversationMeta(conversationId, message);
  }

  private async updateConversationMeta(
    conversationId: string,
    message: Message,
  ): Promise<void> {
    const metaKey = `meta:${conversationId}`;
    const existing = (await this.metaStore.getItem<ConversationMeta>(
      metaKey,
    )) || {
      conversationId,
      count: 0,
      lastTimestamp: 0,
      lastMessageId: "",
    };

    await this.metaStore.setItem(metaKey, {
      ...existing,
      count: existing.count + 1,
      lastTimestamp: message.sent_at,
      lastMessageId: message.id,
    });
  }

  async getMessages(
    conversationId: string,
    options: {
      limit?: number;
      offset?: number;
      reverse?: boolean;
    } = {},
  ): Promise<Message[]> {
    const { limit = 50, offset = 0, reverse = true } = options;

    const allKeys: string[] = [];
    await this.store.iterate((_, key) => {
      if (key.startsWith(`msg:${conversationId}:`)) {
        allKeys.push(key);
      }
    });

    allKeys.sort();
    if (reverse) allKeys.reverse();

    const paginatedKeys = allKeys.slice(offset, offset + limit);

    const messages = await Promise.all(
      paginatedKeys.map((key) => this.store.getItem<Message>(key)),
    );

    return messages.filter(Boolean) as Message[];
  }

  async getMessageCount(conversationId: string): Promise<number> {
    const meta = await this.metaStore.getItem<ConversationMeta>(
      `meta:${conversationId}`,
    );
    return meta?.count ?? 0;
  }

  async getRecentMessages(
    conversationId: string,
    limit = 50,
  ): Promise<Message[]> {
    return this.getMessages(conversationId, {
      limit,
      offset: 0,
      reverse: true,
    }).then((msgs) => msgs.reverse());
  }
}
