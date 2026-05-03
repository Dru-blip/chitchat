import { NextRequest, NextResponse } from "next/server";

const publicRoutes = ["/login", "/verify-link", "/"];

export default async function proxy(req: NextRequest) {
  const sessionCookie = req.cookies.get("chisession")?.value;
  const onboardingCookie = req.cookies.get("onboarding")?.value;

  const path = req.nextUrl.pathname;
  if (!sessionCookie && !publicRoutes.includes(path)) {
    return NextResponse.redirect(new URL("/login", req.nextUrl));
  }

  if (onboardingCookie && !path.startsWith("/onboarding")) {
    return NextResponse.redirect(new URL("/onboarding", req.nextUrl));
  }

  if (sessionCookie && publicRoutes.includes(path)) {
    return NextResponse.redirect(new URL("/chat", req.nextUrl));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};
