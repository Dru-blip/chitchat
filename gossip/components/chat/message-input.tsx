"use client";

import { useState } from "react";
import { HugeiconsIcon } from "@hugeicons/react";
import { SentIcon } from "@hugeicons/core-free-icons";
import { Button } from "@/components/ui/button";
import { apiFetch, cn } from "@/lib/utils";
import { useActiveConversationStore } from "@/stores/providers/active-conversation";
import { KeyBundle, Message } from "@/types";

interface MessageInputProps {
  disabled?: boolean;
  className?: string;
}

export function MessageInput({
  disabled = false,
  className,
}: MessageInputProps) {
  const [text, setText] = useState("");

  const conversation = useActiveConversationStore(
    (state) => state.conversation,
  );
  const otherParticipant = useActiveConversationStore(
    (state) => state.otherParticipant,
  );
  const addMessage = useActiveConversationStore((state) => state.addMessage);

  const sendMessage = async (text: string) => {
    if (!otherParticipant?.user_id) return;

    const { data, error } = await apiFetch<{
      bundle: KeyBundle[];
    }>(`keys/${otherParticipant.user_id}`, {
      method: "POST",
    });

    if (error) {
      console.error("Failed to fetch key bundle:", error);
      return;
    }

    const messageEnvelopes = [];
    for (const bundle of data.bundle) {
      messageEnvelopes.push({
        is_incoming: false,
        recipient_user_id: otherParticipant.user_id,
        recipient_device_id: bundle.deviceId,
        content: text,
      });
    }

    const messagePayload = {
      content_type: "text",
      envelopes: messageEnvelopes,
    };

    const { data: sentMessage, error: sendError } = await apiFetch<Message>(
      `conversations/${conversation?.id}/messages`,
      {
        method: "POST",
        body: JSON.stringify(messagePayload),
      },
    );
    if (sendError) {
      //TODO: toast
      console.error("Failed to send message:", sendError);
    }

    if (sentMessage) {
      addMessage(conversation?.id as string, { ...sentMessage, text: text });
    }
  };

  const handleSubmit = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    const trimmed = text.trim();
    if (!trimmed) return;
    sendMessage(trimmed);
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
