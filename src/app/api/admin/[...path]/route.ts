import { NextResponse } from "next/server";
import { catalogApiBaseUrl } from "@/lib/catalog/api";

/**
 * Same-origin proxy for the admin API.
 *
 * The browser never sees the Go API address; the session cookie stays on this
 * origin and is forwarded per request. Server-side authentication in the Go API
 * remains the only authorization boundary.
 */
const ALLOWED_METHODS = new Set(["GET", "POST", "PATCH", "DELETE"]);
const MAX_BODY_BYTES = 64 * 1024;

function errorResponse(status: number, code: string, message: string) {
  return NextResponse.json({ error: { code, message } }, { status });
}

async function forward(request: Request, context: { params: Promise<{ path?: string[] }> }) {
  if (!ALLOWED_METHODS.has(request.method)) {
    return errorResponse(405, "method_not_allowed", "That method is not supported.");
  }

  const { path } = await context.params;
  const segments = path ?? [];
  // Only the admin API surface may be proxied, and never with path traversal.
  if (segments.length === 0 || segments.some((segment) => segment === "" || segment === "." || segment === ".." || segment.includes("/") || segment.includes("\\"))) {
    return errorResponse(404, "not_found", "The requested resource was not found.");
  }

  const base = catalogApiBaseUrl();
  if (!base) {
    return errorResponse(503, "service_unavailable", "The admin API is not available right now.");
  }

  const headers: Record<string, string> = { accept: "application/json" };
  const cookie = request.headers.get("cookie");
  if (cookie) headers.cookie = cookie;
  const forwardedFor = request.headers.get("x-forwarded-for");
  if (forwardedFor) headers["x-forwarded-for"] = forwardedFor;

  let body: string | undefined;
  if (request.method !== "GET" && request.method !== "DELETE") {
    body = await request.text();
    if (new TextEncoder().encode(body).length > MAX_BODY_BYTES) {
      return errorResponse(413, "validation_error", "The request is too large.");
    }
    headers["content-type"] = request.headers.get("content-type") ?? "application/json";
  }

  let upstream: Response;
  try {
    upstream = await fetch(new URL(`/api/v1/admin/${segments.join("/")}`, base), {
      method: request.method,
      headers,
      body,
      cache: "no-store",
      redirect: "manual",
    });
  } catch {
    return errorResponse(503, "service_unavailable", "The admin API is not available right now.");
  }

  const text = await upstream.text();
  const responseHeaders = new Headers();
  responseHeaders.set("content-type", upstream.headers.get("content-type") ?? "application/json");
  // Forward session cookies set or cleared by the API.
  for (const setCookie of upstream.headers.getSetCookie()) {
    responseHeaders.append("set-cookie", setCookie);
  }
  return new NextResponse(text, { status: upstream.status, headers: responseHeaders });
}

export async function GET(request: Request, context: { params: Promise<{ path?: string[] }> }) {
  return forward(request, context);
}

export async function POST(request: Request, context: { params: Promise<{ path?: string[] }> }) {
  return forward(request, context);
}

export async function PATCH(request: Request, context: { params: Promise<{ path?: string[] }> }) {
  return forward(request, context);
}

export async function DELETE(request: Request, context: { params: Promise<{ path?: string[] }> }) {
  return forward(request, context);
}
