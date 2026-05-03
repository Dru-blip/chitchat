import { Alert, AlertDescription } from "@/components/ui/alert";
import { cn } from "@/lib/utils";
import { HugeiconsIcon, IconSvgElement } from "@hugeicons/react";

interface AlertBannerProps {
  message: string | null;
  icon: IconSvgElement;
  variant?: "destructive" | "default";
  className?: string;
}

export function AlertBanner({
  message,
  icon,
  variant = "default",
  className,
}: AlertBannerProps) {
  return (
    <Alert variant={variant} className={cn(className)}>
      <HugeiconsIcon icon={icon} className="size-4" />
      <AlertDescription>{message}</AlertDescription>
    </Alert>
  );
}
