import type { Metadata } from "next";
import { absoluteSiteUrl, siteUrl } from "@/lib/seo/site-url";

type PageMetadataOptions = {
  title: string;
  description: string;
  pathname: string;
  image: { url: string; alt: string };
};

export function createPageMetadata({ title, description, pathname, image }: PageMetadataOptions): Metadata {
  const canonicalUrl = absoluteSiteUrl(pathname);
  const imageUrl = absoluteSiteUrl(image.url);
  const metadata: Metadata = {
    title,
    description,
    openGraph: {
      type: "website",
      siteName: "EMBER / FIELD",
      title,
      description,
    },
    twitter: {
      card: imageUrl ? "summary_large_image" : "summary",
      title,
      description,
    },
  };

  if (siteUrl) metadata.metadataBase = siteUrl;
  if (canonicalUrl) {
    metadata.alternates = { canonical: canonicalUrl };
    if (metadata.openGraph && !Array.isArray(metadata.openGraph)) metadata.openGraph.url = canonicalUrl;
  }
  if (imageUrl) {
    if (metadata.openGraph && !Array.isArray(metadata.openGraph)) {
      metadata.openGraph.images = [{ url: imageUrl, alt: image.alt }];
    }
    if (metadata.twitter && !Array.isArray(metadata.twitter)) {
      metadata.twitter.images = [{ url: imageUrl, alt: image.alt }];
    }
  }

  return metadata;
}
