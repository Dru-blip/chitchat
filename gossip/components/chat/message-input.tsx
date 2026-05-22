"use client";

import { useState } from "react";
import { HugeiconsIcon } from "@hugeicons/react";
import { SentIcon } from "@hugeicons/core-free-icons";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface MessageInputProps {
  conversationId: string;
  onSend: (text: string) => void;
  disabled?: boolean;
  className?: string;
}

export function MessageInput({
  conversationId,
  onSend,
  disabled = false,
  className,
}: MessageInputProps) {
  const [text, setText] = useState("");

  const handleSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    const trimmed = text.trim();
    if (!trimmed) return;
    onSend(trimmed);
    setText("");
  };

  return (
    <form
      onSubmit={handleSubmit}
      className={cn(
        "flex items-center gap-2 px-4 py-3 border-t border-border bg-background",
        className,
      )}
    >
      <input
        type="text"
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="Type a message..."
        disabled={disabled}
        className={cn(
          "flex-1 h-10 px-4 rounded-3xl border border-transparent bg-input/50",
          "text-sm outline-none placeholder:text-muted-foreground",
          "focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/30",
          "disabled:opacity-50 disabled:pointer-events-none",
        )}
      />
      <Button type="submit" size="icon" disabled={!text.trim() || disabled}>
        <HugeiconsIcon icon={SentIcon} size={18} strokeWidth={1.5} />
      </Button>
    </form>
  );
}
