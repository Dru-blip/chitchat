"use client";

import { useConversationStore } from "@/stores/providers/conversation";
import { ConversationListItem } from "./conversation-list-item";
import { ConversationListEmpty } from "./conversation-list-empty";

interface ConversationListProps {
  activeId?: string;
}

export function ConversationList({ activeId }: ConversationListProps) {
  const conversations = useConversationStore((state) => state.conversations);

  if (conversations.length === 0) return <ConversationListEmpty />;

  return (
    <div className="flex flex-col overflow-y-auto">
      {conversations.map((conv) => (
        <ConversationListItem
          key={conv.id}
          conv={conv}
          isActive={conv.id === activeId}
        />
      ))}
    </div>
  );
}
