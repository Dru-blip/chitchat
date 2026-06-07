"use client";

import {
  type ReactNode,
  createContext,
  useState,
  useContext,
  use,
} from "react";
import { useStore } from "zustand";

import { type MessageStore, createMessageStore } from "@/stores/messages";

export type MessageStoreApi = ReturnType<typeof createMessageStore>;

export const MessageStoreContext = createContext<MessageStoreApi | undefined>(
  undefined,
);

export interface MessageStoreProviderProps {
  children: ReactNode;
}

export const MessageStoreProvider = ({
  children,
}: MessageStoreProviderProps) => {
  const [store] = useState(() => createMessageStore());
  return <MessageStoreContext value={store}>{children}</MessageStoreContext>;
};

export const useMessageStore = <T,>(
  selector: (store: MessageStore) => T,
): T => {
  const messageStoreContext = use(MessageStoreContext);
  if (!messageStoreContext) {
    throw new Error(`useMessageStore must be used within MessageStoreProvider`);
  }
  return useStore(messageStoreContext, selector);
};
