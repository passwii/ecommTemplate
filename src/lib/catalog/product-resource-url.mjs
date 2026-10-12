/**
 * Accept external HTTP(S) URLs or absolute paths that resolve within this storefront.
 * @param {string} value
 * @returns {boolean}
 */
export function isSafeProductResourceUrl(value) {
  const url = value.trim();
  if (!url) return false;
  if (url.startsWith("/")) {
    if (url.startsWith("//") || url.includes("\\") || /[\u0000-\u001f]/.test(url)) return false;
    try {
      return new URL(url, "https://storefront.invalid").origin === "https://storefront.invalid";
    } catch {
      return false;
    }
  }
  try {
    const parsed = new URL(url);
    return (parsed.protocol === "https:" || parsed.protocol === "http:") && Boolean(parsed.hostname);
  } catch {
    return false;
  }
}
