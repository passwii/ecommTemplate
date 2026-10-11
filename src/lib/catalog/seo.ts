import { getStoredProducts } from "@/lib/catalog/product-admin";
import type { Product } from "@/lib/catalog/mock-data";

function isValidPublishedProduct(product: Product): boolean {
  return product.status === "published" && /^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(product.slug);
}

export function getPublicStorefrontProducts(): Product[] {
  return getStoredProducts().filter(isValidPublishedProduct);
}
