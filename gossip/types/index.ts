export interface ErrorResponse {
  code: string;
  message: string;
  details: Record<string, unknown>;
}

export interface User {
  id: string;
  email: string;
  name: string | null;
  image: string | null;
  createdAt: string | null;
  onboarding: boolean;
}

export interface Participant {
  user_id: string;
  email: string;
  name?: string;
  image?: string;
}

export interface Conversation {
  id: string;
  type: string;
  name?: string;
  initiator_id: string;
  created_at: string;
  updated_at: string;
  participants: Participant[];

  last_message?: { text: string; sender_id: string; sent_at: string };
  unread_count?: number;
  is_online?: boolean;
  is_pinned?: boolean;
}

export interface Message {
  id: string;
  conversation_id: string;
  sender_id: string;
  text: string;
  sent_at: string;
}

export interface WebsocketEvent {
  event: 100;
  payload: Record<string, unknown>;
}

export interface KeyBundle {
  deviceId: string;
  clientId: number;
  signedPreKey: string;
  signature: string;
  prekeyId: number;
  prekey: string;
}

export interface ConversationMeta {
  conversationId: string;
  count: number;
  lastTimestamp: string;
  lastMessageId: string;
  lastFetchedAt: string;
}

export const EventType = {
  CONNECTED: 0,
  DISCONNECTED: 1,
  PING: 2,
  PONG: 3,
  NEW_CONVERSATION: 4,
  MESSAGE: 5,
  ERROR: 6,
  QUERY_PRESENCE: 7,
  PRESENCE_RESPONSE: 8,
};
