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
