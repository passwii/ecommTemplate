import Image from "next/image";
import Link from "next/link";
import { getPublishedCollections } from "@/lib/catalog/mock-data";

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

const categories = getPublishedCollections().map((collection) => ({
  ...collection,
  note: categoryNotes[collection.slug] ?? collection.description,
  image: categoryImages[collection.slug] ?? "/images/grill-hero.jpg",
}));

export function CategoryLinks() {
  return (
    <div className="category-grid">
      {categories.map((category) => (
        <Link className="category-tile" href={`/collections/${category.slug}`} key={category.slug}>
          <div className="category-image"><Image src={category.image} alt="" fill sizes="(max-width: 700px) 85vw, 30vw" /></div>
          <div className="category-info"><span><strong>{category.name}</strong><small>{category.note}</small></span><span aria-hidden="true">↗</span></div>
        </Link>
      ))}
    </div>
  );
}
