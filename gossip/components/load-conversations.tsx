"use client";

import { useUserContext } from "@/context/user";
import { useWebsocketContext } from "@/context/websockets";
import { conversationStore } from "@/lib/local-stores";
import { useConversationStore } from "@/stores/providers/conversation";
import { EventType } from "@/types";
import { useEffect } from "react";

export function LoadConversations() {
  const { ws } = useWebsocketContext();
  const { user } = useUserContext();
  const setConversations = useConversationStore(
    (store) => store.setConversations,
  );

  async function loadConversations() {
    return await conversationStore.values();
  }

  useEffect(() => {
    loadConversations().then((conversations) => {
      setConversations(conversations);
    });
  }, []);

  useEffect(() => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      loadConversations().then((conversations) => {
        ws.send(
          JSON.stringify({
            event: EventType.QUERY_PRESENCE,
            payload: {
              users: conversations.map((c) => {
                return c.participants[0].user_id !== user?.id
                  ? c.participants[0].user_id
                  : c.participants[1].user_id;
              }),
            },
          }),
        );
      });
    }
  }, [ws]);

  return null;
}
