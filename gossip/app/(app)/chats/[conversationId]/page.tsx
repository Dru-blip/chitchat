"use client";

import { useParams } from "next/navigation";
import { getMockMessages } from "@/lib/mock-messages";
import { ChatView } from "@/components/chat/chat-view";
import { HugeiconsIcon } from "@hugeicons/react";
import { BubbleChatIcon } from "@hugeicons/core-free-icons";
import { useConversationStore } from "@/stores/providers/conversation";

export default function ConversationPage() {
  const params = useParams<{ conversationId: string }>();
  const conversations = useConversationStore((state) => state.conversations);

  const conversation = conversations.find(
    (c) => c.id === params.conversationId,
  );

  if (!conversation) {
    return <ConversationNotFound />;
  }

  const messages = getMockMessages(params.conversationId);

  return <ChatView conversation={conversation} messages={messages} />;
}

function ConversationNotFound() {
  return (
    <div className="flex-1 flex flex-col items-center justify-center gap-3 p-8">
      <div className="size-16 rounded-4xl bg-muted flex items-center justify-center">
        <HugeiconsIcon
          icon={BubbleChatIcon}
          size={28}
          className="text-muted-foreground"
          strokeWidth={1.2}
        />
      </div>
      <div className="text-center">
        <p className="text-sm font-medium text-foreground">
          Conversation not found
        </p>
        <p className="text-xs text-muted-foreground mt-1">
          This conversation may have been deleted
        </p>
      </div>
    </div>
  );
}
