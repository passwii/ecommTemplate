import { collections, products, type Product, type ProductImage, type ProductSpecification } from "@/lib/catalog/mock-data";

export type ProductDraft = {
  id: string;
  slug: string;
  name: string;
  category: string;
  shortDescription: string;
  description: string;
  images: ProductImage[];
  specifications: ProductSpecification[];
  price: number;
  currency: string;
  compareAtPrice?: number;
  status: "draft" | "published";
  featured: boolean;
  sortOrder: number;
  sku?: string;
};

export type ProductValidationErrors = Partial<Record<"slug" | "name" | "category" | "shortDescription" | "description" | "price" | "images", string>>;

export function productToDraft(product: Product): ProductDraft {
  return {
    id: product.id,
    slug: product.slug,
    name: product.name,
    category: product.category,
    shortDescription: product.shortDescription,
    description: product.description,
    images: product.images.map((image) => ({ ...image })),
    specifications: product.specifications.map((specification) => ({ ...specification })),
    price: product.price,
    currency: product.currency,
    compareAtPrice: product.compareAtPrice,
    status: product.status,
    featured: product.featured,
    sortOrder: product.sortOrder,
    sku: product.sku,
  };
}

export function validateProductDraft(draft: ProductDraft): ProductValidationErrors {
  const errors: ProductValidationErrors = {};
  if (!draft.name.trim()) errors.name = "Enter a product name.";
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(draft.slug)) errors.slug = "Use lowercase letters, numbers and hyphens.";
  if (!collections.some((collection) => collection.slug === draft.category)) errors.category = "Choose a valid collection.";
  if (!draft.shortDescription.trim()) errors.shortDescription = "Enter a short description.";
  if (!draft.description.trim()) errors.description = "Enter a product description.";
  if (!Number.isFinite(draft.price) || draft.price <= 0) errors.price = "Enter a price greater than zero.";
  if (!draft.images.length || draft.images.some((image) => !image.url.trim() || !image.alt.trim())) errors.images = "Add at least one image with a URL and alt text.";
  return errors;
}

export function draftToProduct(draft: ProductDraft, existing?: Product): Product {
  const category = collections.find((item) => item.slug === draft.category)!;
  const orderedImages = draft.images.map((image, index) => ({ ...image, sortOrder: index }));
  const primaryImage = orderedImages[0];
  const specifications = draft.specifications.filter((item) => item.name.trim() && item.value.trim()).map((item) => ({ name: item.name.trim(), value: item.value.trim() }));
  const compatibilitySpecification = (name: string) => specifications.find((item) => item.name.toLowerCase() === name.toLowerCase())?.value;

  return {
    id: draft.id,
    slug: draft.slug.trim(),
    name: draft.name.trim(),
    category: draft.category,
    collection: category.name,
    collectionSlug: category.slug,
    shortDescription: draft.shortDescription.trim(),
    summary: draft.shortDescription.trim(),
    description: draft.description.trim(),
    images: orderedImages,
    image: primaryImage.url.trim(),
    imageAlt: primaryImage.alt.trim(),
    specifications,
    price: draft.price,
    currency: draft.currency,
    compareAtPrice: draft.compareAtPrice,
    status: draft.status,
    featured: draft.featured,
    sortOrder: draft.sortOrder,
    sku: draft.sku?.trim() || undefined,
    badge: existing?.badge,
    fuel: compatibilitySpecification("Fuel") ?? existing?.fuel,
    burners: Number(compatibilitySpecification("Burners")) || existing?.burners,
    cookingArea: compatibilitySpecification("Cooking area") ?? existing?.cookingArea,
    finish: compatibilitySpecification("Finish") ?? existing?.finish,
  };
}

export function getStoredProducts(): Product[] {
  if (typeof window === "undefined") return products;
  try {
    const stored = window.localStorage.getItem("ember-field.catalog.products.v1");
    return stored ? (JSON.parse(stored) as Product[]) : products;
  } catch {
    return products;
  }
}

export function saveStoredProducts(items: Product[]) {
  window.localStorage.setItem("ember-field.catalog.products.v1", JSON.stringify(items));
}

export function resetStoredProducts() {
  window.localStorage.removeItem("ember-field.catalog.products.v1");
}
