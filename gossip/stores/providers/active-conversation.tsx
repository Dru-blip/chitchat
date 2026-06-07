"use client";

import {
  type ReactNode,
  createContext,
  useState,
  useContext,
  use,
} from "react";
import { useStore } from "zustand";

import {
  type ActiveConversationStore,
  createActiveConversationStore,
} from "@/stores/active-conversation";

export type ActiveConversationStoreApi = ReturnType<
  typeof createActiveConversationStore
>;

export const ActiveConversationStoreContext = createContext<
  ActiveConversationStoreApi | undefined
>(undefined);

export interface ActiveConversationStoreProviderProps {
  children: ReactNode;
}

export const ActiveConversationStoreProvider = ({
  children,
}: ActiveConversationStoreProviderProps) => {
  const [store] = useState(() => createActiveConversationStore());
  return (
    <ActiveConversationStoreContext value={store}>
      {children}
    </ActiveConversationStoreContext>
  );
};

export const useActiveConversationStore = <T,>(
  selector: (store: ActiveConversationStore) => T,
): T => {
  const conversationStoreContext = use(ActiveConversationStoreContext);
  if (!conversationStoreContext) {
    throw new Error(
      `useActiveConversationStore must be used within ActiveConversationStoreProvider`,
    );
  }
  return useStore(conversationStoreContext, selector);
};
