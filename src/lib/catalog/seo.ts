import { getPublishedProducts, type Product } from "@/lib/catalog/mock-data";

function hasPublicSlug(product: Product): boolean {
  return /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(product.slug);
}

/**
 * Products in the build-time catalog are the public source of truth for routes and SEO.
 * Browser-local admin edits are previews only and must never affect indexability.
 */
export function getIndexableProducts(): Product[] {
  return getPublishedProducts().filter(hasPublicSlug);
}
