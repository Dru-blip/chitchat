"use client";

import { Button } from "@/components/ui/button";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { getSerializedIdentityKeys } from "@/lib/local-stores";
import { cn } from "@/lib/utils";
import { ErrorResponse } from "@/types";
import { useState } from "react";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);

  const sendMagicLink = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const keys = await getSerializedIdentityKeys();
    const res = await fetch(
      process.env.NEXT_PUBLIC_API_URL + "auth/send-magic-link",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ email, pubkey: keys.pubKey }),
      },
    );

    const data = await res.json();

    if (!res.ok) {
      //TODO: error handling
      throw {
        code: data.code,
        message: data.message,
        details: data.details,
      } satisfies ErrorResponse;
    }
    setSent(true);
    setLoading(false);
  };
  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      {sent && <span>Magic link sent to your email.</span>}
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
            <Button type="submit" disabled={loading}>
              Login
            </Button>
          </Field>
        </FieldGroup>
      </form>
    </div>
  );
}
