"use client";

import { useRef, useEffect } from "react";
import { Message, Participant } from "@/types";
import { MessageBubble } from "./message-bubble";
import { cn } from "@/lib/utils";

interface MessageListProps {
  messages: Message[];
  currentUserId: string;
  otherParticipant?: Participant;
  className?: string;
}

export function MessageList({
  messages,
  currentUserId,
  otherParticipant,
  className,
}: MessageListProps) {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages.length]);

  if (messages.length === 0) {
    return (
      <div
        className={cn(
          "flex-1 flex items-center justify-center",
          className,
        )}
      >
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
        const isOwn = message.sender_id === currentUserId;
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