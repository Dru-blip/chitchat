"use client";

import { useRouter } from "next/navigation";
import { HugeiconsIcon } from "@hugeicons/react";
import { ArrowLeft01Icon } from "@hugeicons/core-free-icons";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { cn, getInitials } from "@/lib/utils";
import { Conversation, Participant } from "@/types";

interface ChatHeaderProps {
  conversation: Conversation;
  otherParticipant?: Participant;
  className?: string;
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

export function ChatHeader({
  conversation,
  otherParticipant,
  className,
}: ChatHeaderProps) {
  const router = useRouter();
  const title = getConversationTitle(conversation, otherParticipant);

  return (
    <div
      className={cn(
        "flex items-center gap-3 px-4 py-3 border-b border-border bg-background",
        className,
      )}
    >
      <Button
        variant="ghost"
        size="icon"
        className="md:hidden shrink-0"
        onClick={() => router.push("/chats")}
      >
        <HugeiconsIcon icon={ArrowLeft01Icon} size={20} strokeWidth={1.5} />
      </Button>

      {otherParticipant && (
        <div className="relative shrink-0">
          <Avatar size="lg">
            <AvatarImage
              src={otherParticipant.image}
              alt={otherParticipant.name || otherParticipant.email}
            />
            <AvatarFallback>
              {getInitials(otherParticipant.name)}
            </AvatarFallback>
          </Avatar>
          {conversation.is_online && (
            <span className="absolute bottom-0 right-0 size-3 rounded-full bg-green-500 ring-2 ring-background" />
          )}
        </div>
      )}

      <div className="flex flex-col min-w-0">
        <span className="font-medium text-sm text-foreground truncate">
          {title}
        </span>
        {conversation.is_online && (
          <span className="text-xs text-green-600 dark:text-green-400">
            Online
          </span>
        )}
      </div>
    </div>
  );
}