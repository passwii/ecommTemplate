import type { MetadataRoute } from "next";
import { absoluteSiteUrl, siteUrl } from "@/lib/seo/site-url";
import { getPublicStorefrontProducts } from "@/lib/catalog/seo";

export default function robots(): MetadataRoute.Robots {
  const robots: MetadataRoute.Robots = {
    rules: {
      userAgent: "*",
      allow: "/",
      disallow: "/admin/",
    },
  };
  const hasPublishedProducts = getPublicStorefrontProducts().length > 0;
  const sitemapUrl = siteUrl && hasPublishedProducts ? absoluteSiteUrl("/sitemap.xml") : undefined;
  if (sitemapUrl) robots.sitemap = sitemapUrl;
  return robots;
}
