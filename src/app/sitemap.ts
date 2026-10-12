import type { MetadataRoute } from "next";
import { getIndexableProducts, getPublishedCollectionSlugs } from "@/lib/catalog/seo";
import { absoluteSiteUrl, siteUrl } from "@/lib/seo/site-url";

export const revalidate = 3600;

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  // Canonical sitemap entries require the confirmed public store origin.
  if (!siteUrl) return [];

  const homeUrl = absoluteSiteUrl("/");
  if (!homeUrl) return [];

  const entries: MetadataRoute.Sitemap = [
    { url: homeUrl, changeFrequency: "weekly", priority: 1 },
  ];

  const [publishedProducts, collectionSlugs] = await Promise.all([
    getIndexableProducts(),
    getPublishedCollectionSlugs(),
  ]);

  for (const slug of collectionSlugs) {
    const url = absoluteSiteUrl(`/collections/${slug}`);
    if (url) entries.push({ url, changeFrequency: "weekly", priority: 0.8 });
  }

  for (const { slug } of publishedProducts) {
    const url = absoluteSiteUrl(`/products/${slug}`);
    if (url) entries.push({ url, changeFrequency: "weekly", priority: 0.7 });
  }

  return entries;
}
