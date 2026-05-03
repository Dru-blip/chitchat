import { cn } from "@/lib/utils";
import { AlertCircleIcon } from "@hugeicons/core-free-icons";
import { AlertBanner } from "./alert-banner";

interface ErrorBannerProps {
  message: string | null;
  className?: string;
}

export function ErrorBanner({ message, className }: ErrorBannerProps) {
  if (!message) return null;

  return (
    <AlertBanner
      message={message}
      icon={AlertCircleIcon}
      variant="destructive"
      className={cn(className)}
    />
  );
}
