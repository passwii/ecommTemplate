import Image from "next/image";
import Link from "next/link";
import type { Collection } from "@/lib/catalog/types";

const categoryImages: Record<string, string> = {
  freestanding: "/images/grill-ridge.jpg",
  "built-in": "/images/grill-summit.jpg",
  compact: "/images/grill-pocket.jpg",
};
const categoryNotes: Record<string, string> = {
  freestanding: "A complete cooking station",
  "built-in": "Made to belong outdoors",
  compact: "Full flavor, smaller footprint",
};

export function CategoryLinks({ collections }: { collections: Collection[] }) {
  return (
    <div className="category-grid">
      {collections.map((collection) => (
        <Link className="category-tile" href={`/collections/${collection.slug}`} key={collection.slug}>
          <div className="category-image"><Image src={categoryImages[collection.slug] ?? "/images/grill-hero.jpg"} alt="" fill sizes="(max-width: 700px) 85vw, 30vw" /></div>
          <div className="category-info"><span><strong>{collection.name}</strong><small>{categoryNotes[collection.slug] ?? collection.description}</small></span><span aria-hidden="true">↗</span></div>
        </Link>
      ))}
    </div>
  );
}
