import { NextRequest, NextResponse } from "next/server";

const publicRoutes = ["/login", "/verify-link", "/"];

export default async function proxy(req: NextRequest) {
  const sessionCookie = req.cookies.get("chisession");

  const path = req.nextUrl.pathname;
  if (!sessionCookie && !publicRoutes.includes(path)) {
    return NextResponse.redirect(new URL("/login", req.nextUrl));
  }

  if (sessionCookie && publicRoutes.includes(path)) {
    return NextResponse.redirect(new URL("/chat", req.nextUrl));
  }

  return NextResponse.next();
}
