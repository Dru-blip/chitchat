"use client";

import { useActiveConversationStore } from "@/stores/providers/active-conversation";
import { useConversationStore } from "@/stores/providers/conversation";
import { EventType } from "@/types";
import { createContext, useContext, useEffect, useRef, useState } from "react";

interface WebsocketContextType {
  ws: WebSocket | null;
  connected: boolean;
}

const websocketContext = createContext<WebsocketContextType | undefined>(
  undefined,
);

export function WebsocketProvider({ children }: { children: React.ReactNode }) {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const pingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const addConversation = useConversationStore(
    (state) => state.addConversation,
  );
  const setPresence = useConversationStore((state) => state.setPresence);
  const addMessage = useActiveConversationStore((state) => state.addMessage);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:5050/ws");

    socket.onopen = () => {
      setConnected(true);
      setWs(socket);
      pingIntervalRef.current = setInterval(() => {
        if (socket.readyState === WebSocket.OPEN) {
          socket.send(JSON.stringify({ event: EventType.PING }));
        }
      }, 3000);
    };

    socket.onclose = () => {
      setConnected(false);
      setWs(null);
    };

    socket.onerror = (err) => {
      console.error("websocket error:", err);
    };

    socket.onmessage = (evt) => {
      const data = JSON.parse(evt.data);
      const payload = data.payload;
      switch (data.event) {
        case EventType.PONG: {
          console.log("ws: connected");
          break;
        }
        case EventType.NEW_CONVERSATION: {
          addConversation(payload);
          break;
        }
        case EventType.MESSAGE: {
          addMessage(payload.conversation_id, payload);
          break;
        }
        case EventType.PRESENCE_RESPONSE: {
          setPresence(payload);
          break;
        }
      }
    };

    return () => {
      socket.close();
      if (pingIntervalRef.current) {
        clearInterval(pingIntervalRef.current);
        pingIntervalRef.current = null;
      }
    };
  }, []);

  return (
    <websocketContext.Provider value={{ ws, connected }}>
      {children}
    </websocketContext.Provider>
  );
}

export function useWebsocketContext() {
  const context = useContext(websocketContext);
  if (context === undefined) {
    throw new Error(
      "useWebsocketContext must be used within a WebsocketProvider",
    );
  }
  return context;
}
