import { ChatsShell } from "@/components/chat/chat-shell";
import { LoadConversations } from "@/components/load-conversations";
import { UploadPrekeys } from "@/components/upload-prekeys";

export default function ChatsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <LoadConversations />
      <UploadPrekeys />
      <ChatsShell>{children}</ChatsShell>
    </>
  );
}
