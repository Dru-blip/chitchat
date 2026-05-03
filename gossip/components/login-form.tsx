"use client";

import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { getSerializedIdentityKeys } from "@/lib/local-stores";
import { apiFetch, cn } from "@/lib/utils";
import { Loading01FreeIcons } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import { differenceInSeconds, intervalToDuration } from "date-fns";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { ErrorBanner } from "./error-banner";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [serverError, setServerError] = useState("");
  const [secondsLeft, setSecondsLeft] = useState(0);
  const [cooldownUntil, setCooldownUntil] = useState<Date | null>(null);

  useEffect(() => {
    if (!cooldownUntil) return;

    const counter = setInterval(() => {
      const seconds = differenceInSeconds(cooldownUntil, new Date());
      setSecondsLeft(seconds);

      if (seconds <= 0) {
        clearInterval(counter);
        setCooldownUntil(null);
        return;
      }
    }, 1000);

    return () => clearInterval(counter);
  }, [cooldownUntil]);

  const formatSeconds = (seconds: number) => {
    const duration = intervalToDuration({
      start: 0,
      end: seconds * 1000,
    });
    return `${duration.minutes ?? 0}:${duration.seconds ?? 0}`;
  };

  const sendMagicLink = async (e: React.SubmitEvent<HTMLFormElement>) => {
    setServerError("");
    e.preventDefault();
    setLoading(true);

    const keys = await getSerializedIdentityKeys();
    const { data, error } = await apiFetch<
      { message: string; retryAfter: string; email: string },
      { message: string } & { retryAfter: string }
    >("auth/send-magic-link", {
      method: "POST",
      body: JSON.stringify({ email, pubkey: keys.pubKey }),
    });

    if (error) {
      if (error.retryAfter) setCooldownUntil(new Date(error.retryAfter));
      setServerError(error.message ?? "Failed to send magic link");
      setLoading(false);
      return;
    }

    setCooldownUntil(new Date(data.retryAfter));

    toast.success("Magic link sent to your email.");
    setLoading(false);
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      {serverError && <ErrorBanner message={serverError} />}
      <form onSubmit={sendMagicLink}>
        <FieldGroup>
          <Field>
            <FieldLabel htmlFor="email">Email</FieldLabel>
            <Input
              id="email"
              type="email"
              onChange={(e) => setEmail(e.target.value)}
              placeholder="m@example.com"
              required
            />
          </Field>
          <Field>
            <Button type="submit" disabled={loading || secondsLeft > 0}>
              {loading ? (
                <>
                  <HugeiconsIcon icon={Loading01FreeIcons} /> Sending
                </>
              ) : secondsLeft > 0 ? (
                `Resend in ${formatSeconds(secondsLeft)}`
              ) : (
                "Login"
              )}
            </Button>
          </Field>
        </FieldGroup>
      </form>
    </div>
  );
}
