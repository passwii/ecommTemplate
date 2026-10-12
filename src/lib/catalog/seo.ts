import { fetchCollections, fetchProducts } from "@/lib/catalog/api";
import type { Product } from "@/lib/catalog/types";

function hasPublicSlug(product: Product): boolean {
  return /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(product.slug);
}

/** Products published by this store's catalog API are the public source of truth. */
export function isIndexableProduct(product: Product): boolean {
  return hasPublicSlug(product);
}

/**
 * Published, indexable products from the catalog API. Browser-local admin
 * previews are not part of this source and never affect indexability.
 */
export async function getIndexableProducts(): Promise<Product[]> {
  const result = await fetchProducts({ limit: 100 });
  if (!result.ok) return [];
  return result.data.filter(isIndexableProduct);
}

/** Published collection slugs, used for sitemap and robots decisions. */
export async function getPublishedCollectionSlugs(): Promise<string[]> {
  const result = await fetchCollections();
  if (!result.ok) return [];
  return result.data.map((collection) => collection.slug);
}
