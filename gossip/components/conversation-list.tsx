"use client";

import Link from "next/link";
import { HugeiconsIcon } from "@hugeicons/react";
import { BubbleChatIcon } from "@hugeicons/core-free-icons";
import {
  Avatar,
  AvatarFallback,
  AvatarGroup,
  AvatarGroupCount,
  AvatarImage,
} from "@/components/ui/avatar";
import { useConversations } from "@/context/conversations";
import { useUserContext } from "@/context/user";
import { cn, formatRelativeTime, getInitials } from "@/lib/utils";
import { Conversation, Participant } from "@/types";

interface ConversationListProps {
  activeId?: string;
}

export function ConversationList({ activeId }: ConversationListProps) {
  const { conversations, loading, error } = useConversations();

  if (loading) return <ConversationListSkeleton />;
  if (error) return <ConversationListError message={error} />;
  if (conversations.length === 0) return <ConversationListEmpty />;

  return (
    <div className="flex flex-col overflow-y-auto">
      {conversations.map((conv) => (
        <ConversationItem
          key={conv.id}
          conv={conv}
          isActive={conv.id === activeId}
        />
      ))}
    </div>
  );
}

function ConversationItem({
  conv,
  isActive,
}: {
  conv: Conversation;
  isActive: boolean;
}) {
  const { user } = useUserContext();

  const otherParticipant = conv.participants.find(
    (p) => p.user_id !== user?.id,
  );

  const title = getConversationTitle(conv, otherParticipant);
  const timestamp = formatRelativeTime(
    conv.last_message?.sent_at ?? conv.updated_at,
  );
  const preview = conv.last_message?.text;

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

function ConversationAvatar({
  conv,
  otherParticipant,
}: {
  conv: Conversation;
  otherParticipant?: Participant;
}) {
  if (conv.type === "dm" && otherParticipant) {
    const src = otherParticipant.image;
    const initials = getInitials(otherParticipant.name);

    return (
      <div className="relative shrink-0">
        <Avatar size="lg">
          <AvatarImage src={src} alt={otherParticipant.name || otherParticipant.email} />
          <AvatarFallback>{initials}</AvatarFallback>
        </Avatar>
        {conv.is_online && (
          <span className="absolute bottom-0 right-0 size-3 rounded-full bg-green-500 ring-2 ring-background" />
        )}
      </div>
    );
  }

  const displayParticipants = conv.participants.slice(0, 2);
  const remaining = conv.participants.length - 2;

  return (
    <AvatarGroup>
      {displayParticipants.map((p) => (
        <Avatar key={p.user_id} size="sm">
          <AvatarImage
            src={p.image}
            alt={p.name || p.email}
          />
          <AvatarFallback className="text-[9px]">
            {getInitials(p.name)}
          </AvatarFallback>
        </Avatar>
      ))}
      {remaining > 0 && <AvatarGroupCount>+{remaining}</AvatarGroupCount>}
    </AvatarGroup>
  );
}

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

function ConversationListSkeleton() {
  return (
    <div className="flex flex-col">
      {Array.from({ length: 5 }).map((_, i) => (
        <div key={i} className="flex items-center gap-3 px-4 py-3">
          <div className="size-10 rounded-full bg-muted shrink-0 animate-pulse" />
          <div className="flex-1 min-w-0 space-y-2">
            <div className="h-3.5 w-28 bg-muted rounded animate-pulse" />
            <div className="h-3 w-40 bg-muted rounded animate-pulse" />
          </div>
        </div>
      ))}
    </div>
  );
}

function ConversationListEmpty() {
  return (
    <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
      <HugeiconsIcon
        icon={BubbleChatIcon}
        size={40}
        className="text-muted-foreground/40 mb-3"
        strokeWidth={1.2}
      />
      <p className="text-sm font-medium text-foreground">
        No conversations yet
      </p>
      <p className="text-xs text-muted-foreground mt-1">
        Start a new conversation to get going
      </p>
    </div>
  );
}

function ConversationListError({ message }: { message: string }) {
  return <div className="p-4 text-sm text-destructive">{message}</div>;
}