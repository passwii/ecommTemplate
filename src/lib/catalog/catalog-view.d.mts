import type { Collection, Product } from "./types";

export type ApiProductImage = {
  url: string;
  alt: string;
  sortOrder: number;
};

export type ApiProductSpecification = {
  name: string;
  value: string;
};

export type ApiProductResource = {
  title: string;
  url: string;
  fileType?: string;
  description?: string;
};

export type ApiCollection = {
  id: number;
  slug: string;
  name: string;
  description: string;
  sortOrder: number;
};

export type ApiProduct = {
  id: string;
  slug: string;
  name: string;
  model?: string;
  sku?: string;
  collection: ApiCollection;
  shortDescription: string;
  description: string;
  images: ApiProductImage[];
  specifications: ApiProductSpecification[];
  applications: string[];
  resources: ApiProductResource[];
  priceAmount?: number;
  priceCurrency?: string;
  compareAtPrice?: number;
  badge?: string;
  fuel?: string;
  burners?: number;
  cookingArea?: string;
  finish?: string;
  minimumOrderQuantity?: string;
  leadTime?: string;
  featured: boolean;
  sortOrder: number;
  contentStatus: string;
};

export type ApiCompanyContent = {
  about: { heading: string; body: string; imageUrl: string; imageAlt: string };
  manufacturing: { heading: string; overview: string; qualityHeading: string; qualityDetails: string };
  credentials: { name: string; issuer: string; details: string }[];
  caseStudies: { customer: string; title: string; details: string }[];
  faqs: { question: string; answer: string }[];
  contentStatus: string;
};

export function toStorefrontProduct(apiProduct: ApiProduct): Product;
export function toStorefrontCollection(apiCollection: ApiCollection): Collection;
export function hasPublishedPrice(product: { price?: number }): boolean;
