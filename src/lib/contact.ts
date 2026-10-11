export const CONTACT_PATH = "/contact";

export function getInquiryPath(productSlug?: string) {
  if (!productSlug) return CONTACT_PATH;
  return `${CONTACT_PATH}?product=${encodeURIComponent(productSlug)}`;
}
