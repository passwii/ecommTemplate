/**
 * Maps storefront API catalog payloads to the view model the gas-grill template
 * renders. Keeping this pure (and tested) means the template never depends on
 * the API transport or the database shape.
 */

/**
 * @param {{
 *   id: string, slug: string, name: string, model?: string, sku?: string,
 *   collection?: { slug?: string, name?: string },
 *   shortDescription?: string, description?: string,
 *   images?: Array<{ url: string, alt: string, sortOrder?: number }>,
 *   specifications?: Array<{ name: string, value: string }>,
 *   applications?: string[], resources?: Array<{ title: string, url: string, fileType?: string, description?: string }>,
 *   priceAmount?: number, priceCurrency?: string, compareAtPrice?: number,
 *   badge?: string, fuel?: string, burners?: number, cookingArea?: string, finish?: string,
 *   minimumOrderQuantity?: string, leadTime?: string,
 *   featured?: boolean, sortOrder?: number
 * }} apiProduct
 */
export function toStorefrontProduct(apiProduct) {
  const images = [...(apiProduct.images ?? [])]
    .filter((image) => image && typeof image.url === "string" && image.url.trim() !== "")
    .sort((a, b) => (a.sortOrder ?? 0) - (b.sortOrder ?? 0));
  const primaryImage = images[0];
  const shortDescription = apiProduct.shortDescription ?? "";

  return {
    id: apiProduct.id,
    slug: apiProduct.slug,
    name: apiProduct.name,
    category: apiProduct.collection?.slug ?? "",
    collection: apiProduct.collection?.name ?? "",
    collectionSlug: apiProduct.collection?.slug ?? "",
    shortDescription,
    description: apiProduct.description ?? "",
    images: images.map((image, index) => ({
      url: image.url,
      alt: image.alt ?? "",
      sortOrder: image.sortOrder ?? index,
    })),
    specifications: (apiProduct.specifications ?? []).map((specification) => ({
      name: specification.name,
      value: specification.value,
    })),
    applications: (apiProduct.applications ?? []).filter((application) => typeof application === "string"),
    resources: (apiProduct.resources ?? []).map((resource) => ({
      title: resource.title,
      url: resource.url,
      fileType: resource.fileType,
      description: resource.description,
    })),
    price: apiProduct.priceAmount ?? 0,
    currency: apiProduct.priceCurrency || "USD",
    compareAtPrice: apiProduct.compareAtPrice,
    status: "published",
    featured: Boolean(apiProduct.featured),
    sortOrder: apiProduct.sortOrder ?? 0,
    sku: apiProduct.sku,
    minimumOrderQuantity: apiProduct.minimumOrderQuantity,
    leadTime: apiProduct.leadTime,
    summary: shortDescription,
    image: primaryImage?.url ?? "",
    imageAlt: primaryImage?.alt ?? "",
    badge: apiProduct.badge,
    fuel: apiProduct.fuel,
    burners: apiProduct.burners,
    cookingArea: apiProduct.cookingArea,
    finish: apiProduct.finish,
  };
}

/**
 * @param {{ slug: string, name: string, description?: string }} apiCollection
 */
export function toStorefrontCollection(apiCollection) {
  return {
    slug: apiCollection.slug,
    name: apiCollection.name,
    description: apiCollection.description ?? "",
  };
}

/**
 * The template shows a price only when the store actually publishes one.
 * @param {{ price?: number, priceCurrency?: string }} product
 */
export function hasPublishedPrice(product) {
  return typeof product.price === "number" && product.price > 0;
}
