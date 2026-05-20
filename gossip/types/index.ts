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
}
