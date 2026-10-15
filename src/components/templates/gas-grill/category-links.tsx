import Image from "next/image";
import Link from "next/link";
import { getPublishedCollections } from "@/lib/catalog/mock-data";

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

const categories = getPublishedCollections().map((collection) => ({
  ...collection,
  note: categoryNotes[collection.slug] ?? collection.description,
  image: categoryImages[collection.slug] ?? { src: "/images/grill-hero.jpg", alt: "Outdoor patio lounge with a fire pit beneath string lights" },
}));

export function CategoryLinks() {
  return (
    <div className="category-grid">
      {categories.map((category) => (
        <Link className="category-tile" href={`/collections/${category.slug}`} key={category.slug}>
          <div className="category-image"><Image src={category.image.src} alt={category.image.alt} fill sizes="(max-width: 767px) 90vw, 30vw" /></div>
          <div className="category-info"><span><strong>{category.name}</strong><small>{category.note}</small></span><span aria-hidden="true">↗</span></div>
        </Link>
      ))}
    </div>
  );
}
