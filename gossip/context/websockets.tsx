"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";

interface WebsocketContextType {
  ws: WebSocket | null;
  connected: boolean;
}

const EventType = {
  CONNECTED: 0,
  DISCONNECTED: 1,
  PING: 2,
  PONG: 3,
  MESSAGE: 4,
};

const websocketContext = createContext<WebsocketContextType | undefined>(
  undefined,
);

export function WebsocketProvider({ children }: { children: React.ReactNode }) {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const pingIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

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
      const msg = JSON.parse(evt.data);
      switch (msg.event) {
        case EventType.PONG: {
          console.log("server: connected");
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
