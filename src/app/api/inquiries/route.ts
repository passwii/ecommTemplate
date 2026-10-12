import { NextResponse } from "next/server";
import { catalogApiBaseUrl } from "@/lib/catalog/api";

/**
 * Same-origin proxy for inquiry submissions.
 *
 * The browser posts to this route so the Go API stays a server-side detail:
 * no CORS setup, no public API URL in client code, and the visitor's address is
 * forwarded for the API's rate limiter when a trusted proxy set it.
 */
const MAX_BODY_BYTES = 64 * 1024;

function errorResponse(status: number, code: string, message: string) {
  return NextResponse.json({ error: { code, message } }, { status });
}

export async function POST(request: Request) {
  const base = catalogApiBaseUrl();
  if (!base) {
    return errorResponse(503, "service_unavailable", "Inquiry submissions are not available right now.");
  }

  const contentType = request.headers.get("content-type") ?? "";
  if (!contentType.toLowerCase().startsWith("application/json")) {
    return errorResponse(415, "unsupported_media_type", "Send the request as JSON.");
  }

  const body = await request.text();
  if (new TextEncoder().encode(body).length > MAX_BODY_BYTES) {
    return errorResponse(413, "validation_error", "The submission is too large.");
  }

  const headers: Record<string, string> = {
    "content-type": "application/json",
    accept: "application/json",
  };
  // Only forwarded when the hosting platform already provided it; the Go API
  // ignores this header unless TRUST_PROXY_HEADER is enabled.
  const forwardedFor = request.headers.get("x-forwarded-for");
  if (forwardedFor) headers["x-forwarded-for"] = forwardedFor;

  try {
    const upstream = await fetch(new URL("/api/v1/inquiries", base), {
      method: "POST",
      headers,
      body,
      cache: "no-store",
    });
    const text = await upstream.text();
    return new NextResponse(text, {
      status: upstream.status,
      headers: { "content-type": "application/json" },
    });
  } catch {
    return errorResponse(503, "service_unavailable", "Inquiry submissions are not available right now.");
  }
}
