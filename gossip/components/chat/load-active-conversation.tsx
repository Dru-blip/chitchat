"use client";

import { useUserContext } from "@/context/user";
import { messageStore } from "@/lib/local-stores";
import { useActiveConversationStore } from "@/stores/providers/active-conversation";
import { useConversationStore } from "@/stores/providers/conversation";
import { useParams } from "next/navigation";
import { useEffect } from "react";

export function LoadActiveConversation() {
  const params = useParams();
  const { user } = useUserContext();
  const setConversation = useActiveConversationStore(
    (state) => state.setConversation,
  );
  const setMessages = useActiveConversationStore((state) => state.setMessages);
  const setOtherParticipant = useActiveConversationStore(
    (state) => state.setOtherParticipant,
  );
  const conversations = useConversationStore((state) => state.conversations);

  useEffect(() => {
    const conversation = conversations.find(
      (c) => c.id === params.conversationId,
    );

    if (conversation) {
      setConversation(conversation);
      setOtherParticipant(
        conversation.participants.find((p) => p.user_id !== user?.id)!,
      );
      messageStore.getRecentMessages(conversation.id).then((msgs) => {
        setMessages(msgs);
      });

      messageStore.fetchLatestMessages(conversation.id).then((latest) => {
        setMessages(latest);
      });
    }
  }, [params.conversationId, conversations]);

  return null;
}
