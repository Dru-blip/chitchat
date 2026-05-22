"use client";

import { useUserContext } from "@/context/user";
import { Conversation, Message, Participant } from "@/types";
import { ChatHeader } from "./chat-header";
import { MessageList } from "./message-list";
import { MessageInput } from "./message-input";

interface ChatViewProps {
  conversation: Conversation;
  messages: Message[];
}

export function ChatView({ conversation, messages }: ChatViewProps) {
  const { user } = useUserContext();

  const otherParticipant = conversation.participants.find(
    (p) => p.user_id !== user?.id,
  );

  const handleSend = (text: string) => {
    // Stub — will be replaced with real API call
    console.log("Send message:", text, "to conversation:", conversation.id);
  };

  return (
    <div className="flex flex-col h-full pb-16 md:pb-0 animate-slide-in-right md:animate-none">
      <ChatHeader
        conversation={conversation}
        otherParticipant={otherParticipant}
      />
      <MessageList
        messages={messages}
        currentUserId={user?.id ?? ""}
        otherParticipant={otherParticipant}
      />
      <MessageInput
        conversationId={conversation.id}
        onSend={handleSend}
      />
    </div>
  );
}