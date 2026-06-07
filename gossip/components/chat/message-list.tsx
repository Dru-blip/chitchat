"use client";

import { useUserContext } from "@/context/user";
import { cn } from "@/lib/utils";
import { useActiveConversationStore } from "@/stores/providers/active-conversation";
import { Participant } from "@/types";
import { useEffect, useRef } from "react";
import { MessageBubble } from "./message-bubble";

interface MessageListProps {
  className?: string;
}

export function MessageList({ className }: MessageListProps) {
  const bottomRef = useRef<HTMLDivElement>(null);
  const { user } = useUserContext();

  const otherParticipant = useActiveConversationStore(
    (state) => state.otherParticipant,
  ) as Participant;

  const messages = useActiveConversationStore((state) => state.messages);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages.length]);

  if (messages.length === 0) {
    return (
      <div className={cn("flex-1 flex items-center justify-center", className)}>
        <p className="text-sm text-muted-foreground">
          No messages yet. Say hello!
        </p>
      </div>
    );
  }

  return (
    <div
      className={cn(
        "flex-1 overflow-y-auto px-4 py-3 flex flex-col gap-1",
        className,
      )}
    >
      {messages.map((message, index) => {
        const isOwn = message.sender_id === user?.id;
        const nextMessage = messages[index + 1];
        const showAvatar =
          !nextMessage || nextMessage.sender_id !== message.sender_id;

        return (
          <MessageBubble
            key={message.id}
            message={message}
            isOwn={isOwn}
            showAvatar={showAvatar}
            otherParticipant={otherParticipant}
          />
        );
      })}
      <div ref={bottomRef} />
    </div>
  );
}
