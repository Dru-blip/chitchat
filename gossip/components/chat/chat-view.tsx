"use client";

import { useActiveConversationStore } from "@/stores/providers/active-conversation";
import { ConversationNotFound } from "../conversation-not-found";
import { ChatHeader } from "./chat-header";
import { MessageInput } from "./message-input";
import { MessageList } from "./message-list";

export function ChatView() {
  const conversation = useActiveConversationStore(
    (state) => state.conversation,
  );

  if (!conversation) return <ConversationNotFound />;

  return (
    <div className="flex flex-col h-full pb-16 md:pb-0 animate-slide-in-right md:animate-none">
      <ChatHeader />
      <MessageList />
      <MessageInput />
    </div>
  );
}

// const setMessages = useMessageStore((store) => store.setMessages);
// const addMessage = useMessageStore((store) => store.addMessage);
// const messages = useMessageStore(
//   (store) => store.getMessages(conversation.id) ?? [],
// );

// const otherParticipant = conversation.participants.find(
//   (p) => p.user_id !== user?.id,
// );

// useEffect(() => {
//   messageStore.getMessages(conversation.id).then((msgs) => {
//     setMessages(conversation.id, msgs);
//   });
// }, [conversation.id, setMessages]);
