import { collections, products, type Product, type ProductImage, type ProductResource, type ProductSpecification } from "@/lib/catalog/mock-data";

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
  applications: string[];
  minimumOrderQuantity?: string;
  leadTime?: string;
  resources: ProductResource[];
};

export type ProductValidationErrors = Partial<Record<"slug" | "name" | "category" | "shortDescription" | "description" | "price" | "images" | "applications" | "minimumOrderQuantity" | "leadTime" | "resources", string>>;

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
    applications: [...(product.applications ?? [])],
    minimumOrderQuantity: product.minimumOrderQuantity,
    leadTime: product.leadTime,
    resources: (product.resources ?? []).map((resource) => ({ ...resource })),
  };
}

export function isSafeProductResourceUrl(value: string): boolean {
  const url = value.trim();
  if (!url) return false;
  if (url.startsWith("/")) return !url.startsWith("//") && !url.includes("\\\\") && !/[\u0000-\u001f]/.test(url);
  try {
    const parsed = new URL(url);
    return (parsed.protocol === "https:" || parsed.protocol === "http:") && Boolean(parsed.hostname);
  } catch {
    return false;
  }
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
  if (draft.applications.some((application) => !application.trim())) errors.applications = "Complete or remove each application entry.";
  if (draft.minimumOrderQuantity !== undefined && draft.minimumOrderQuantity.trim() === "") errors.minimumOrderQuantity = "Complete this field or clear it to leave it optional.";
  if (draft.leadTime !== undefined && draft.leadTime.trim() === "") errors.leadTime = "Complete this field or clear it to leave it optional.";
  if (draft.resources.some((resource) => !resource.title.trim() || !resource.url.trim())) errors.resources = "Each resource needs a title and URL, or remove the incomplete row.";
  else if (draft.resources.some((resource) => !isSafeProductResourceUrl(resource.url))) errors.resources = "Use an http(s) URL or a path within this storefront.";
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
    applications: draft.applications.map((application) => application.trim()).filter(Boolean),
    minimumOrderQuantity: draft.minimumOrderQuantity?.trim() || undefined,
    leadTime: draft.leadTime?.trim() || undefined,
    resources: draft.resources.filter((resource) => resource.title.trim() && resource.url.trim()).map((resource) => ({
      title: resource.title.trim(),
      url: resource.url.trim(),
      fileType: resource.fileType?.trim() || undefined,
      description: resource.description?.trim() || undefined,
    })),
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
