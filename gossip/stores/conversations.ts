import { conversationStore } from "@/lib/local-stores";
import { Conversation } from "@/types";
import { createStore } from "zustand";

interface ConversationState {
  conversations: Conversation[];
}

interface ConversationActions {
  addConversation: (c: Conversation) => void;
  setConversations: (c: Conversation[]) => void;
}

export const defaultInitState: ConversationState = {
  conversations: [],
};

export type ConversationStore = ConversationState & ConversationActions;

export const createConversationStore = (
  initState: ConversationState = defaultInitState,
) => {
  return createStore<ConversationStore>()((set) => ({
    ...initState,
    addConversation: (c) => {
      conversationStore.set(c.id, c);
      set((state) => ({
        conversations: [...state.conversations, c],
      }));
    },
    setConversations: (c) =>
      set(() => ({
        conversations: c,
      })),
  }));
};
