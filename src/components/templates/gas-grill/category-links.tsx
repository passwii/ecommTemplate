import Image from "next/image";
import Link from "next/link";

const categories = [
  { slug: "freestanding", name: "Freestanding grills", note: "A complete cooking station", image: "/images/grill-ridge.jpg" },
  { slug: "built-in", name: "Built-in grills", note: "Made to belong outdoors", image: "/images/grill-summit.jpg" },
  { slug: "compact", name: "Compact grills", note: "Full flavor, smaller footprint", image: "/images/grill-pocket.jpg" },
];

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
