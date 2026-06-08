"use client";

import { useActiveConversationStore } from "@/stores/providers/active-conversation";
import { ConversationNotFound } from "../conversation-not-found";
import { ChatHeader } from "./chat-header";
import { MessageInput } from "./message-input";
import { MessageList } from "./message-list";

export function ChatView() {
  const conversation = useActiveConversationStore(
    (state) => state.conversation,
  );

  if (!conversation) return <ConversationNotFound />;

  return (
    <div className="flex flex-col h-full pb-16 md:pb-0 animate-slide-in-right md:animate-none">
      <ChatHeader />
      <MessageList />
      <MessageInput />
    </div>
  );
}
