import { Message } from "@/types";

const fallbackTexts = [
  "Hey! How are you?",
  "I'm doing great, thanks for asking!",
  "Did you see the latest update?",
  "Yes, it looks really promising!",
  "We should discuss it further",
  "Absolutely, let's chat later today",
  "Sounds good, I'll be around",
  "Perfect, talk soon!",
];

function generateFallbackMessages(conversationId: string): Message[] {
  const now = new Date();
  return Array.from({ length: 8 }, (_, i) => ({
    id: `${conversationId}-msg-${i}`,
    conversation_id: conversationId,
    sender_id: i % 2 === 0 ? "self" : "other",
    text: fallbackTexts[i],
    sent_at: new Date(now.getTime() - (8 - i) * 3600000).toISOString(),
  }));
}

const mockMessagesByConversation: Record<string, Message[]> = {};

export function getMockMessages(conversationId: string): Message[] {
  return mockMessagesByConversation[conversationId] ?? generateFallbackMessages(conversationId);
}