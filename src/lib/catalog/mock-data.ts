export type ProductImage = {
  url: string;
  alt: string;
  sortOrder: number;
};

export type ProductSpecification = {
  name: string;
  value: string;
};

export type ProductResource = {
  title: string;
  url: string;
  fileType?: string;
  description?: string;
};

export type Product = {
  id: string;
  slug: string;
  name: string;
  category: string;
  collection: string;
  collectionSlug: string;
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
  /** Optional trade details; older catalog records may not include them. */
  applications?: string[];
  minimumOrderQuantity?: string;
  leadTime?: string;
  resources?: ProductResource[];
  /** Storefront compatibility fields: aliases for the primary image and short description. */
  summary: string;
  image: string;
  imageAlt: string;
  badge?: string;
  fuel?: string;
  burners?: number;
  cookingArea?: string;
  finish?: string;
};

export type Collection = {
  slug: string;
  name: string;
  description: string;
};

export const collections: Collection[] = [
  { slug: "built-in", name: "Built-in grills", description: "A considered centerpiece for an outdoor kitchen that feels like part of home." },
  { slug: "freestanding", name: "Freestanding grills", description: "A complete cooking station, ready to make the patio your favorite room." },
  { slug: "compact", name: "Compact grills", description: "Full-flavor cooking, thoughtfully sized for smaller outdoor spaces." },
];

export const products: Product[] = [
  {
    id: "grill-01", slug: "ridge-4-burner", name: "Ridge 4-Burner", category: "freestanding", collection: "Freestanding grills", collectionSlug: "freestanding",
    price: 2499, currency: "USD", shortDescription: "Room to cook for the whole table.", summary: "Room to cook for the whole table.", description: "A capable centerpiece for generous weekends outside. Four independently controlled burners give you the range to sear, roast and keep a side warm, all at once.",
    images: [{ url: "/images/grill-ridge.jpg", alt: "Patio lounge seating around an outdoor fire pit beneath a covered, tree-lined terrace", sortOrder: 0 }], image: "/images/grill-ridge.jpg", imageAlt: "Patio lounge seating around an outdoor fire pit beneath a covered, tree-lined terrace",
    specifications: [{ name: "Burners", value: "4" }, { name: "Cooking area", value: "748 sq in" }, { name: "Fuel", value: "Natural gas or propane" }, { name: "Finish", value: "Brushed stainless steel" }],
    badge: "Bestseller", fuel: "Natural gas or propane", burners: 4, cookingArea: "748 sq in", finish: "Brushed stainless steel", status: "published", featured: true, sortOrder: 10,
  },
  {
    id: "grill-02", slug: "hearth-3-burner", name: "Hearth 3-Burner", category: "freestanding", collection: "Freestanding grills", collectionSlug: "freestanding",
    price: 1899, currency: "USD", shortDescription: "Everyday cooking, made effortless.", summary: "Everyday cooking, made effortless.", description: "The grill you reach for on a Tuesday and gather around on Saturday. Three responsive burners make reliable heat feel simple.",
    images: [{ url: "/images/grill-hearth.jpg", alt: "Warm wood-paneled cabin kitchen with a dining area beside large windows", sortOrder: 0 }], image: "/images/grill-hearth.jpg", imageAlt: "Warm wood-paneled cabin kitchen with a dining area beside large windows",
    specifications: [{ name: "Burners", value: "3" }, { name: "Cooking area", value: "560 sq in" }, { name: "Fuel", value: "Propane" }, { name: "Finish", value: "Graphite and stainless steel" }],
    fuel: "Propane", burners: 3, cookingArea: "560 sq in", finish: "Graphite and stainless steel", status: "published", featured: false, sortOrder: 20,
  },
  {
    id: "grill-03", slug: "summit-built-in", name: "Summit Built-In", category: "built-in", collection: "Built-in grills", collectionSlug: "built-in",
    price: 3299, currency: "USD", shortDescription: "Made to belong in your outdoor kitchen.", summary: "Made to belong in your outdoor kitchen.", description: "Clean lines and steady, even heat bring a thoughtful finish to the outdoor kitchen. Designed to sit flush with the space around it.",
    images: [{ url: "/images/grill-summit.jpg", alt: "Two people seated outdoors beside a fire pit beneath string lights", sortOrder: 0 }], image: "/images/grill-summit.jpg", imageAlt: "Two people seated outdoors beside a fire pit beneath string lights",
    specifications: [{ name: "Burners", value: "5" }, { name: "Cooking area", value: "860 sq in" }, { name: "Fuel", value: "Natural gas or propane" }, { name: "Finish", value: "Brushed stainless steel" }],
    badge: "New", fuel: "Natural gas or propane", burners: 5, cookingArea: "860 sq in", finish: "Brushed stainless steel", status: "published", featured: true, sortOrder: 30,
  },
  {
    id: "grill-04", slug: "pocket-2-burner", name: "Pocket 2-Burner", category: "compact", collection: "Compact grills", collectionSlug: "compact",
    price: 1099, currency: "USD", shortDescription: "A smaller footprint. No smaller appetite.", summary: "A smaller footprint. No smaller appetite.", description: "Thoughtfully sized for balconies, courtyards and first homes. Two burners give you direct and indirect heat without taking over the space.",
    images: [{ url: "/images/grill-pocket.jpg", alt: "Person cooking at an outdoor grill in a leafy garden", sortOrder: 0 }], image: "/images/grill-pocket.jpg", imageAlt: "Person cooking at an outdoor grill in a leafy garden",
    specifications: [{ name: "Burners", value: "2" }, { name: "Cooking area", value: "360 sq in" }, { name: "Fuel", value: "Propane" }, { name: "Finish", value: "Matte graphite" }],
    fuel: "Propane", burners: 2, cookingArea: "360 sq in", finish: "Matte graphite", status: "published", featured: false, sortOrder: 40,
  },
];

export const formatPrice = (price: number, currency = "USD") => new Intl.NumberFormat("en-US", { style: "currency", currency, maximumFractionDigits: 0 }).format(price);
export const getProduct = (slug: string) => products.find((product) => product.slug === slug);
export const getCollection = (slug: string) => collections.find((collection) => collection.slug === slug);
export const getPublishedProducts = () => products.filter((product) => product.status === "published");
export const getPublishedCollections = () => {
  const publishedCollectionSlugs = new Set(getPublishedProducts().map((product) => product.collectionSlug));
  return collections.filter((collection) => publishedCollectionSlugs.has(collection.slug));
};
