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
import { cn } from "@/lib/utils";
import { ErrorResponse } from "@/types";
import { useRef, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import * as z from "zod";

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
  const [done, setDone] = useState(false);
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
    const formData = new FormData();
    formData.append("name", values.name);
    formData.append("password", values.password);
    if (values.image) {
      formData.append("image", values.image);
    }

    const res = await fetch(
      process.env.NEXT_PUBLIC_API_URL + "users/onboarding",
      {
        method: "POST",
        credentials: "include",
        body: formData,
      },
    );

    const data = await res.json();

    if (!res.ok) {
      // Map API field errors back into RHF if present
      if (data.details) {
        for (const [field, message] of Object.entries(data.details)) {
          form.setError(field as keyof OnboardingFormValues, {
            message: message as string,
          });
        }
        return;
      }
      throw {
        code: data.code,
        message: data.message,
        details: data.details,
      } satisfies ErrorResponse;
    }

    setDone(true);
  };

  if (done) {
    return (
      <div className={cn("flex flex-col gap-6", className)} {...props}>
        <span>You're all set! Welcome aboard.</span>
      </div>
    );
  }

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
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
