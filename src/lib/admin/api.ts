import { cookies } from "next/headers";
import { catalogApiBaseUrl } from "@/lib/catalog/api";

/**
 * Server-side admin API client. Server components read the session cookie from
 * the incoming request and forward it, so admin pages never expose the API
 * address or the session token to client JavaScript.
 */
export type AdminResult<T> =
  | { ok: true; data: T; meta?: Record<string, unknown> }
  | { ok: false; reason: "unauthorized" | "not_found" | "unavailable" | "error" };

export type AdminInquiry = {
  id: string;
  status: string;
  contact: { name: string; email: string; company: string; phone: string };
  destinationCountry: string;
  product: { id?: string; slug: string; model: string; name: string };
  quantity: number | null;
  customizationRequest: string;
  sourcePath: string;
  privacyNoticeVersion: string;
  privacyConsentedAt: string;
  owner: string;
  lastContactedAt: string | null;
  nextAction: string;
  nextActionDueAt: string | null;
  lossReason: string;
  createdAt: string;
  updatedAt: string;
  activities?: AdminActivity[];
};

export type AdminActivity = {
  id: number;
  type: string;
  body: string;
  actor: string;
  createdAt: string;
};

export const INQUIRY_STATUSES = ["new", "in_progress", "quoted", "won", "lost", "on_hold"] as const;

async function adminRequest<T>(path: string): Promise<AdminResult<T>> {
  const base = catalogApiBaseUrl();
  if (!base) return { ok: false, reason: "unavailable" };

  const cookieStore = await cookies();
  const cookieHeader = cookieStore
    .getAll()
    .map((cookie) => `${cookie.name}=${cookie.value}`)
    .join("; ");

  let response: Response;
  try {
    response = await fetch(new URL(path, base), {
      headers: { accept: "application/json", ...(cookieHeader ? { cookie: cookieHeader } : {}) },
      cache: "no-store",
    });
  } catch {
    return { ok: false, reason: "unavailable" };
  }

  if (response.status === 401) return { ok: false, reason: "unauthorized" };
  if (response.status === 404) return { ok: false, reason: "not_found" };
  if (!response.ok) return { ok: false, reason: "error" };

  try {
    const payload = (await response.json()) as { data?: T; meta?: Record<string, unknown> };
    if (payload?.data === undefined) return { ok: false, reason: "error" };
    return { ok: true, data: payload.data, meta: payload.meta };
  } catch {
    return { ok: false, reason: "error" };
  }
}

export async function fetchAdminInquiries(options: { status?: string; q?: string; limit?: number; offset?: number } = {}) {
  const search = new URLSearchParams();
  if (options.status) search.set("status", options.status);
  if (options.q) search.set("q", options.q);
  search.set("limit", String(options.limit ?? 50));
  search.set("offset", String(options.offset ?? 0));
  return adminRequest<AdminInquiry[]>(`/api/v1/admin/inquiries?${search.toString()}`);
}

export async function fetchAdminInquiry(id: string) {
  return adminRequest<AdminInquiry>(`/api/v1/admin/inquiries/${encodeURIComponent(id)}`);
}

export async function fetchAdminSession() {
  return adminRequest<{ email: string }>("/api/v1/admin/session");
}

export function formatInquiryTimestamp(value: string | null): string {
  if (!value) return "—";
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return "—";
  return parsed.toLocaleString("en-GB", { dateStyle: "medium", timeStyle: "short", timeZone: "UTC" });
}

export function statusLabel(status: string): string {
  switch (status) {
    case "new": return "New";
    case "in_progress": return "In progress";
    case "quoted": return "Quoted";
    case "won": return "Won";
    case "lost": return "Lost";
    case "on_hold": return "On hold";
    default: return status;
  }
}
