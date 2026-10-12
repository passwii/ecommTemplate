import { toStorefrontCollection, toStorefrontProduct } from "./catalog-view.mjs";
import type { ApiCollection, ApiCompanyContent, ApiProduct } from "./catalog-view.mjs";
import type { Collection, Product } from "./types";

/**
 * Storefront data comes from this store's own Go API. The URL is server-side
 * configuration only: it is deliberately not a NEXT_PUBLIC_* value.
 */
const REVALIDATE_SECONDS = 60;

export type CatalogResult<T> =
  | { ok: true; data: T }
  | { ok: false; reason: "not_found" | "unavailable" };

export type CompanyContent = ApiCompanyContent;

function apiBaseUrl(): string | undefined {
  const value = process.env.CATALOG_API_URL?.trim();
  if (!value) return undefined;
  try {
    const url = new URL(value);
    if (url.protocol !== "http:" && url.protocol !== "https:") return undefined;
    return url.origin;
  } catch {
    return undefined;
  }
}

/** True when the storefront is configured to read from a catalog API. */
export function isCatalogApiConfigured(): boolean {
  return apiBaseUrl() !== undefined;
}

async function fetchCatalog<T>(path: string): Promise<CatalogResult<T>> {
  const base = apiBaseUrl();
  if (!base) return { ok: false, reason: "unavailable" };

  let response: Response;
  try {
    response = await fetch(new URL(path, base), {
      headers: { accept: "application/json" },
      next: { revalidate: REVALIDATE_SECONDS },
    });
  } catch {
    // Never silently fall back to stale sample data: report the API as unavailable.
    return { ok: false, reason: "unavailable" };
  }

  if (response.status === 404) return { ok: false, reason: "not_found" };
  if (!response.ok) return { ok: false, reason: "unavailable" };

  try {
    const payload = (await response.json()) as { data?: T };
    if (payload?.data === undefined) return { ok: false, reason: "unavailable" };
    return { ok: true, data: payload.data };
  } catch {
    return { ok: false, reason: "unavailable" };
  }
}

export async function fetchProducts(options: { collection?: string; featured?: boolean; limit?: number } = {}): Promise<CatalogResult<Product[]>> {
  const search = new URLSearchParams();
  if (options.collection) search.set("collection", options.collection);
  if (options.featured) search.set("featured", "true");
  if (options.limit) search.set("limit", String(options.limit));
  const query = search.toString();

  const result = await fetchCatalog<ApiProduct[]>(`/api/v1/products${query ? `?${query}` : ""}`);
  if (!result.ok) return result;
  return { ok: true, data: result.data.map(toStorefrontProduct) };
}

export async function fetchProduct(slug: string): Promise<CatalogResult<Product>> {
  const result = await fetchCatalog<ApiProduct>(`/api/v1/products/${encodeURIComponent(slug)}`);
  if (!result.ok) return result;
  return { ok: true, data: toStorefrontProduct(result.data) };
}

export async function fetchCollections(): Promise<CatalogResult<Collection[]>> {
  const result = await fetchCatalog<ApiCollection[]>("/api/v1/collections");
  if (!result.ok) return result;
  return { ok: true, data: result.data.map(toStorefrontCollection) };
}

export async function fetchCollection(slug: string): Promise<CatalogResult<Collection>> {
  const result = await fetchCatalog<ApiCollection>(`/api/v1/collections/${encodeURIComponent(slug)}`);
  if (!result.ok) return result;
  return { ok: true, data: toStorefrontCollection(result.data) };
}

export async function fetchCompanyContent(): Promise<CatalogResult<CompanyContent>> {
  return fetchCatalog<CompanyContent>("/api/v1/company-content");
}
