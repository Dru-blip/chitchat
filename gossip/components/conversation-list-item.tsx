"use client";

import { useUserContext } from "@/context/user";
import { cn, formatRelativeTime } from "@/lib/utils";
import { Conversation, ConversationMeta, Message, Participant } from "@/types";
import Link from "next/link";
import { ConversationAvatar } from "./conversation-avatar";
import { useEffect, useState } from "react";
import { messageStore } from "@/lib/local-stores";

function UnreadBadge({ count }: { count: number }) {
  return (
    <span className="inline-flex items-center justify-center min-w-5 h-5 px-1.5 rounded-full bg-primary text-primary-foreground text-[11px] font-semibold leading-none shrink-0">
      {count > 99 ? "99+" : count}
    </span>
  );
}

function getConversationTitle(
  conv: Conversation,
  otherParticipant?: Participant,
): string {
  if (conv.name) return conv.name;
  if (conv.type === "dm" && otherParticipant) {
    return otherParticipant.name || otherParticipant.email || "Direct Message";
  }
  return "Group Chat";
}

export function ConversationListItem({
  conv,
  isActive,
}: {
  conv: Conversation;
  isActive: boolean;
}) {
  const { user } = useUserContext();
  const [convMeta, setConvMeta] = useState<ConversationMeta | null>(null);
  const [lastMessage, setLastMessage] = useState<Message | null>(null);

  const otherParticipant = conv.participants.find(
    (p) => p.user_id !== user?.id,
  );

  useEffect(() => {
    messageStore.getConversationMeta(conv.id).then((meta) => {
      if (meta) {
        setConvMeta(meta);
        messageStore
          .getLastMessageById(conv.id, meta.lastTimestamp, meta.lastMessageId)
          .then((lastMessage) => {
            if (lastMessage) {
              setLastMessage(lastMessage);
            }
          });
      }
    });
  }, [conv.id]);

  const title = getConversationTitle(conv, otherParticipant);
  const timestamp = formatRelativeTime(convMeta?.lastTimestamp);
  const preview = lastMessage?.text;

  return (
    <Link
      href={`/chats/${conv.id}`}
      className={cn(
        "flex items-center gap-3 px-4 py-3 rounded-xl transition-colors",
        "hover:bg-accent/50",
        isActive && "bg-accent",
      )}
    >
      <ConversationAvatar conv={conv} otherParticipant={otherParticipant} />

      <div className="flex-1 min-w-0">
        <div className="flex items-baseline justify-between gap-2">
          <span className="font-medium text-sm text-foreground truncate">
            {title}
          </span>
          {timestamp && (
            <span className="text-xs text-muted-foreground shrink-0 tabular-nums">
              {timestamp}
            </span>
          )}
        </div>

        <div className="flex items-center justify-between gap-2 mt-0.5">
          {preview ? (
            <span className="text-sm text-muted-foreground truncate">
              {preview}
            </span>
          ) : (
            <span className="text-sm text-muted-foreground/50 truncate">
              No messages yet
            </span>
          )}
          {conv.unread_count != null && conv.unread_count > 0 && (
            <UnreadBadge count={conv.unread_count} />
          )}
        </div>
      </div>
    </Link>
  );
}
