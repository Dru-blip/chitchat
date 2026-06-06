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
