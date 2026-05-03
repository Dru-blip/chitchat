"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { apiFetch, cn } from "@/lib/utils";
import { useRef, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import * as z from "zod";
import { getSerializedIdentityKeys } from "@/lib/local-stores";
import { useRouter } from "next/navigation";
import { ErrorBanner } from "./error-banner";
import { toast } from "sonner";

const onboardingSchema = z.object({
  name: z
    .string()
    .min(2, "Name must be at least 2 characters.")
    .max(64, "Name must be at most 64 characters."),
  password: z
    .string()
    .min(8, "Password must be at least 8 characters.")
    .max(128, "Password must be at most 128 characters."),
  image: z
    .instanceof(File)
    .refine((f) => f.size <= 5 * 1024 * 1024, "Image must be under 5MB.")
    .refine(
      (f) => ["image/jpeg", "image/png", "image/webp"].includes(f.type),
      "Only JPEG, PNG, or WebP images are allowed.",
    )
    .optional(),
});

type OnboardingFormValues = z.infer<typeof onboardingSchema>;

export function UserOnboardingForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const router = useRouter();
  const [serverError, setServerError] = useState("");
  const fileInputRef = useRef<HTMLInputElement>(null);

  const form = useForm<OnboardingFormValues>({
    resolver: zodResolver(onboardingSchema),
    defaultValues: {
      name: "",
      password: "",
      image: undefined,
    },
  });

  const onSubmit = async (values: OnboardingFormValues) => {
    setServerError("");

    const idKeys = await getSerializedIdentityKeys();
    const { error } = await apiFetch<
      void,
      { message: string; details?: Record<keyof OnboardingFormValues, string> }
    >("users/onboard", {
      method: "PATCH",
      body: JSON.stringify({ ...values, pubkey: idKeys.pubKey }),
    });

    if (error) {
      setServerError(error.message);

      if (error.details) {
        for (const [field, message] of Object.entries(error.details)) {
          form.setError(field as keyof OnboardingFormValues, { message });
        }
      }
      return;
    }

    toast.success("Successfully onboarded");
    router.push("/chat");
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <ErrorBanner message={serverError} />
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <FieldGroup>
          {/* Name */}
          <Controller
            name="name"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="name">Name</FieldLabel>
                <Input
                  {...field}
                  id="name"
                  type="text"
                  placeholder="Jane Doe"
                  autoComplete="name"
                  aria-invalid={fieldState.invalid}
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />

          {/* Password */}
          <Controller
            name="password"
            control={form.control}
            render={({ field, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="password">Password</FieldLabel>
                <Input
                  {...field}
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  autoComplete="new-password"
                  aria-invalid={fieldState.invalid}
                />
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />

          {/* Profile Image (optional) */}
          <Controller
            name="image"
            control={form.control}
            render={({ field: { onChange, value, ...rest }, fieldState }) => (
              <Field data-invalid={fieldState.invalid}>
                <FieldLabel htmlFor="image">
                  Profile Image{" "}
                  <span className="text-muted-foreground font-normal">
                    (optional)
                  </span>
                </FieldLabel>
                <Input
                  {...rest}
                  ref={fileInputRef}
                  id="image"
                  type="file"
                  accept="image/jpeg,image/png,image/webp"
                  aria-invalid={fieldState.invalid}
                  onChange={(e) => onChange(e.target.files?.[0])}
                />
                <FieldDescription>
                  JPEG, PNG, or WebP — max 5 MB.
                </FieldDescription>
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            )}
          />

          <Field>
            <Button type="submit" disabled={form.formState.isSubmitting}>
              {form.formState.isSubmitting ? "Saving…" : "Get started"}
            </Button>
          </Field>
        </FieldGroup>
      </form>
    </div>
  );
}
