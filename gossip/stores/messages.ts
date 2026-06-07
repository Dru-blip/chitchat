import { messageStore } from "@/lib/local-stores";
import { Message } from "@/types";
import { createStore } from "zustand";

interface MessageState {
  messages: Map<string, Message[]>;
}

interface MessageActions {
  addMessage: (conversationId: string, m: Message) => void;
  getMessages: (conversationId: string) => Message[];
  setMessages: (conversationId: string, messages: Message[]) => void;
}

export const defaultInitState: MessageState = {
  messages: new Map(),
};

export type MessageStore = MessageState & MessageActions;

export const createMessageStore = (
  initState: MessageState = defaultInitState,
) => {
  return createStore<MessageStore>()((set, get) => ({
    ...initState,
    addMessage: (conversationId, m) => {
      messageStore.appendMessage(conversationId, m);
      set((state) => ({
        messages: new Map(state.messages).set(conversationId, [
          ...(state.messages.get(conversationId) ?? []),
          m,
        ]),
      }));
    },
    setMessages: (conversationId: string, messages: Message[]) => {
      set((state) => ({
        messages: new Map(state.messages).set(conversationId, messages),
      }));
    },
    getMessages: (conversationId: string) => {
      return get().messages.get(conversationId) ?? [];
    },
  }));
};
