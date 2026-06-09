"use client";

import {
  Avatar,
  AvatarFallback,
  AvatarGroup,
  AvatarGroupCount,
  AvatarImage,
  AvatarBadge,
} from "@/components/ui/avatar";
import { getInitials } from "@/lib/utils";
import { Conversation, Participant } from "@/types";

export function ConversationAvatar({
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
          <AvatarImage
            src={src}
            alt={otherParticipant.name || otherParticipant.email}
          />
          <AvatarFallback>{initials}</AvatarFallback>
          {conv.is_online && (
            <AvatarBadge className="bg-green-600 dark:bg-green-800" />
          )}
        </Avatar>
      </div>
    );
  }

  const displayParticipants = conv.participants.slice(0, 2);
  const remaining = conv.participants.length - 2;

  return (
    <AvatarGroup>
      {displayParticipants.map((p) => (
        <Avatar key={p.user_id} size="sm">
          <AvatarImage src={p.image} alt={p.name || p.email} />
          <AvatarFallback className="text-[9px]">
            {getInitials(p.name)}
          </AvatarFallback>
        </Avatar>
      ))}
      {remaining > 0 && <AvatarGroupCount>+{remaining}</AvatarGroupCount>}
    </AvatarGroup>
  );
}
