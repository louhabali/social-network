import { NextResponse } from "next/server";

export async function middleware(request) {
  const token = request.cookies.get("Token");

  const { pathname } = request.nextUrl;

  if (!token && pathname !== "/login" && pathname !== "/register") {
    return NextResponse.redirect(new URL("/login", request.url));
  }

  if (token) {
    const apiResponse = await fetch("http://localhost:8080/api/validate-token", {
      method: "POST",
      body: JSON.stringify({ token : token.value }),
    });
    
    if (!apiResponse.ok) {
      const response = NextResponse.redirect(new URL("/login", request.url));
      response.cookies.delete("Token");
      return response;
    }

    if (pathname === "/login" || pathname === "/register") {
      return NextResponse.redirect(new URL("/", request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};
