"use client";

import { type ReactNode, createContext, use, useState } from "react";
import { useStore } from "zustand";

import {
  type ConversationStore,
  createConversationStore,
} from "@/stores/conversations";

export type ConversationStoreApi = ReturnType<typeof createConversationStore>;

export const ConversationStoreContext = createContext<
  ConversationStoreApi | undefined
>(undefined);

export interface ConversationStoreProviderProps {
  children: ReactNode;
}

export const ConversationStoreProvider = ({
  children,
}: ConversationStoreProviderProps) => {
  const [store] = useState(() => createConversationStore());
  return (
    <ConversationStoreContext value={store}>
      {children}
    </ConversationStoreContext>
  );
};

export const useConversationStore = <T,>(
  selector: (store: ConversationStore) => T,
): T => {
  const conversationStoreContext = use(ConversationStoreContext);
  if (!conversationStoreContext) {
    throw new Error(
      `useConversationStore must be used within ConversationStoreProvider`,
    );
  }
  return useStore(conversationStoreContext, selector);
};
