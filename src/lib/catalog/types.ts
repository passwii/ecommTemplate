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
