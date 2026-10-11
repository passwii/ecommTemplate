const configuredSiteUrl = process.env.NEXT_PUBLIC_SITE_URL?.trim();

function parseSiteUrl(value: string | undefined): URL | undefined {
  if (!value) return undefined;

  try {
    const url = new URL(value);
    if (url.protocol !== "http:" && url.protocol !== "https:") return undefined;
    return new URL(url.origin);
  } catch {
    return undefined;
  }
}

/** The public storefront origin. Leave unset until the store's real domain is confirmed. */
export const siteUrl = parseSiteUrl(configuredSiteUrl);

export function absoluteSiteUrl(pathname: string): string | undefined {
  return siteUrl ? new URL(pathname, siteUrl).toString() : undefined;
}
