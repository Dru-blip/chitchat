"use client";

import { conversationStore } from "@/lib/local-stores";
import { useConversationStore } from "@/stores/providers/conversation";
import { useEffect } from "react";

export function LoadConversations() {
  const setConversations = useConversationStore(
    (store) => store.setConversations,
  );

  async function loadConversations() {
    return await conversationStore.values();
  }

  useEffect(() => {
    loadConversations().then(setConversations);
  }, []);

  return null;
}
