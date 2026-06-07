import { messageStore } from "@/lib/local-stores";
import { Conversation, Message, Participant } from "@/types";
import { createStore } from "zustand";

interface ActiveConversationState {
  conversation: Conversation | null;
  messages: Message[];
  otherParticipant: Participant | null;
}

interface ActiveConversationActions {
  setConversation: (c: Conversation) => void;
  setMessages: (messages: Message[]) => void;
  setOtherParticipant: (p: Participant) => void;
  addMessage: (conversationId: string, message: Message) => void;
}

export const defaultInitState: ActiveConversationState = {
  conversation: null,
  messages: [],
  otherParticipant: null,
};

export type ActiveConversationStore = ActiveConversationState &
  ActiveConversationActions;

export const createActiveConversationStore = (
  initState: ActiveConversationState = defaultInitState,
) => {
  return createStore<ActiveConversationStore>()((set) => ({
    ...initState,
    setConversation: (c: Conversation) => set(() => ({ conversation: c })),
    setOtherParticipant: (p: Participant) =>
      set(() => ({ otherParticipant: p })),
    setMessages: (messages: Message[]) => set(() => ({ messages })),
    addMessage: (conversationId: string, message: Message) => {
      messageStore.appendMessage(conversationId, message);
      set((state) => ({ messages: [...state.messages, message] }));
    },
  }));
};
