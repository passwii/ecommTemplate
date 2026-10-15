import Image from "next/image";
import Link from "next/link";
import type { Collection } from "@/lib/catalog/types";

const categoryImages: Record<string, { src: string; alt: string }> = {
  freestanding: { src: "/images/grill-ridge.jpg", alt: "Patio lounge seating around an outdoor fire pit beneath a covered terrace" },
  "built-in": { src: "/images/grill-summit.jpg", alt: "Two people seated outdoors beside a fire pit beneath string lights" },
  compact: { src: "/images/grill-pocket.jpg", alt: "Person cooking at an outdoor grill in a leafy garden" },
};
const categoryNotes: Record<string, string> = {
  freestanding: "A complete cooking station",
  "built-in": "Made to belong outdoors",
  compact: "Full flavor, smaller footprint",
};

const FALLBACK_IMAGE = { src: "/images/grill-hero.jpg", alt: "Outdoor patio lounge with a fire pit beneath string lights" };

export function CategoryLinks({ collections }: { collections: Collection[] }) {
  return (
    <div className="category-grid">
      {collections.map((collection) => {
        const image = categoryImages[collection.slug] ?? FALLBACK_IMAGE;
        return (
          <Link className="category-tile" href={`/collections/${collection.slug}`} key={collection.slug}>
            <div className="category-image"><Image src={image.src} alt={image.alt} fill sizes="(max-width: 767px) 90vw, 30vw" /></div>
            <div className="category-info"><span><strong>{collection.name}</strong><small>{categoryNotes[collection.slug] ?? collection.description}</small></span><span aria-hidden="true">↗</span></div>
          </Link>
        );
      })}
    </div>
  );
}
