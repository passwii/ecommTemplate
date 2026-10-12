import type { MetadataRoute } from "next";
import { absoluteSiteUrl, siteUrl } from "@/lib/seo/site-url";
import { getIndexableProducts } from "@/lib/catalog/seo";

export const revalidate = 3600;

export default async function robots(): Promise<MetadataRoute.Robots> {
  const robots: MetadataRoute.Robots = {
    rules: {
      userAgent: "*",
      allow: "/",
      disallow: "/admin/",
    },
  };
  const hasPublishedProducts = (await getIndexableProducts()).length > 0;
  const sitemapUrl = siteUrl && hasPublishedProducts ? absoluteSiteUrl("/sitemap.xml") : undefined;
  if (sitemapUrl) robots.sitemap = sitemapUrl;
  return robots;
}
