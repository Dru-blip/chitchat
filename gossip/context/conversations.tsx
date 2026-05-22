"use client";

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
} from "react";
import { apiFetch } from "@/lib/utils";
import { Conversation } from "@/types";

interface ConversationsContextType {
  conversations: Conversation[];
  loading: boolean;
  error: string | null;
  addConversation: (conversation: Conversation) => void;
}

const ConversationsContext = createContext<
  ConversationsContextType | undefined
>(undefined);

export function ConversationsProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchConversations = async () => {
      const { data, error: fetchError } = await apiFetch<
        Conversation[],
        { message: string }
      >("conversations");

      if (fetchError) {
        setError(fetchError.message);
      } else if (data) {
        setConversations(data);
      }
      setLoading(false);
    };
    fetchConversations();
  }, []);

  const addConversation = useCallback((conversation: Conversation) => {
    setConversations((prev) => [conversation, ...prev]);
  }, []);

  return (
    <ConversationsContext.Provider
      value={{ conversations, loading, error, addConversation }}
    >
      {children}
    </ConversationsContext.Provider>
  );
}

export function useConversations() {
  const context = useContext(ConversationsContext);
  if (!context) {
    throw new Error(
      "useConversations must be used within a ConversationsProvider",
    );
  }
  return context;
}
