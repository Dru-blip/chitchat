"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn, formatRelativeTime, getInitials } from "@/lib/utils";
import { Message, Participant } from "@/types";

interface MessageBubbleProps {
  message: Message;
  isOwn: boolean;
  showAvatar: boolean;
  otherParticipant?: Participant;
  className?: string;
}

export function MessageBubble({
  message,
  isOwn,
  showAvatar,
  otherParticipant,
  className,
}: MessageBubbleProps) {
  const time = formatRelativeTime(message.sent_at);

  return (
    <div
      className={cn(
        "flex gap-2 max-w-[75%]",
        isOwn ? "ml-auto flex-row-reverse" : "mr-auto",
        className,
      )}
    >
      {!isOwn && showAvatar && otherParticipant ? (
        <div className="shrink-0 self-end">
          <Avatar size="sm">
            <AvatarImage
              src={otherParticipant.image}
              alt={otherParticipant.name}
            />
            <AvatarFallback className="text-[9px]">
              {getInitials(otherParticipant.name)}
            </AvatarFallback>
          </Avatar>
        </div>
      ) : (
        !isOwn && <div className="shrink-0 w-6" />
      )}

      <div
        className={cn(
          "px-3 py-2 rounded-3xl text-sm",
          isOwn
            ? "bg-primary text-primary-foreground"
            : "bg-card text-card-foreground border border-border",
        )}
      >
        <p className="whitespace-pre-wrap wrap-break-word">{message.text}</p>
        <span
          className={cn(
            "block text-[10px] mt-0.5 tabular-nums",
            isOwn
              ? "text-primary-foreground/60 text-right"
              : "text-muted-foreground",
          )}
        >
          {time}
        </span>
      </div>
    </div>
  );
}
