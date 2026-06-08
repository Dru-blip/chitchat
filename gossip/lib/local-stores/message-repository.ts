import { ConversationMeta, Message } from "@/types";
import localforage from "localforage";
import { apiFetch } from "../utils";

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
      lastTimestamp: "0",
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

  async getConversationMeta(
    conversationId: string,
  ): Promise<ConversationMeta | null> {
    return await this.metaStore.getItem<ConversationMeta>(
      `meta:${conversationId}`,
    );
  }

  async getLastMessageById(
    conversationId: string,
    timestamp: string,
    messageId: string,
  ): Promise<Message | null> {
    return await this.store.getItem<Message>(
      `msg:${conversationId}:${timestamp}:${messageId}`,
    );
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

  async fetchLatestMessages(
    conversationId: string,
    limit = 50,
  ): Promise<Message[]> {
    const convMeta = await this.getConversationMeta(conversationId);
    const lastMessageId = convMeta?.lastMessageId;
    if (!lastMessageId) {
      return this.getRecentMessages(conversationId, limit);
    }

    const lastMessage = await this.getLastMessageById(
      conversationId,
      convMeta.lastTimestamp,
      lastMessageId,
    );
    if (!lastMessage) {
      return this.getRecentMessages(conversationId, limit);
    }

    const { data, error } = await apiFetch<Message[]>(
      `conversations/${conversationId}/messages?timestamp=${new Date(lastMessage.sent_at).toISOString()}`,
      {},
    );

    if (error) {
      return [];
    }

    if (data.length === 0) {
      return this.getRecentMessages(conversationId, limit);
    }

    for (const msg of data) {
      await this.appendMessage(conversationId, msg);
    }

    return this.getRecentMessages(conversationId, limit);
  }
}
