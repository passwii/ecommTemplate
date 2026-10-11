export type Product = {
  id: string;
  slug: string;
  name: string;
  collection: string;
  collectionSlug: string;
  price: number;
  summary: string;
  description: string;
  image: string;
  imageAlt: string;
  badge?: string;
  fuel: string;
  burners: number;
  cookingArea: string;
  finish: string;
};

export type Collection = {
  slug: string;
  name: string;
  description: string;
};

export const collections: Collection[] = [
  { slug: "built-in", name: "Built-in grills", description: "A considered centerpiece for an outdoor kitchen that feels like part of home." },
  { slug: "freestanding", name: "Freestanding grills", description: "A complete cooking station, ready to make the patio your favorite room." },
  { slug: "compact", name: "Compact grills", description: "Full-flavor cooking, thoughtfully sized for smaller outdoor spaces." },
];

export const products: Product[] = [
  {
    id: "grill-01", slug: "ridge-4-burner", name: "Ridge 4-Burner", collection: "Freestanding grills", collectionSlug: "freestanding",
    price: 2499, summary: "Room to cook for the whole table.", description: "A capable centerpiece for generous weekends outside. Four independently controlled burners give you the range to sear, roast and keep a side warm, all at once.",
    image: "/images/grill-ridge.jpg", imageAlt: "Stainless steel gas grill staged outdoors against a wood deck", badge: "Bestseller", fuel: "Natural gas or propane", burners: 4, cookingArea: "748 sq in", finish: "Brushed stainless steel",
  },
  {
    id: "grill-02", slug: "hearth-3-burner", name: "Hearth 3-Burner", collection: "Freestanding grills", collectionSlug: "freestanding",
    price: 1899, summary: "Everyday cooking, made effortless.", description: "The grill you reach for on a Tuesday and gather around on Saturday. Three responsive burners make reliable heat feel simple.",
    image: "/images/grill-hearth.jpg", imageAlt: "Outdoor gas grill on a timber patio", fuel: "Propane", burners: 3, cookingArea: "560 sq in", finish: "Graphite and stainless steel",
  },
  {
    id: "grill-03", slug: "summit-built-in", name: "Summit Built-In", collection: "Built-in grills", collectionSlug: "built-in",
    price: 3299, summary: "Made to belong in your outdoor kitchen.", description: "Clean lines and steady, even heat bring a thoughtful finish to the outdoor kitchen. Designed to sit flush with the space around it.",
    image: "/images/grill-summit.jpg", imageAlt: "Stainless steel built-in grill in an outdoor setting", badge: "New", fuel: "Natural gas or propane", burners: 5, cookingArea: "860 sq in", finish: "Brushed stainless steel",
  },
  {
    id: "grill-04", slug: "pocket-2-burner", name: "Pocket 2-Burner", collection: "Compact grills", collectionSlug: "compact",
    price: 1099, summary: "A smaller footprint. No smaller appetite.", description: "Thoughtfully sized for balconies, courtyards and first homes. Two burners give you direct and indirect heat without taking over the space.",
    image: "/images/grill-pocket.jpg", imageAlt: "Compact grill for cooking outdoors", fuel: "Propane", burners: 2, cookingArea: "360 sq in", finish: "Matte graphite",
  },
];

export const formatPrice = (price: number) => new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 0 }).format(price);
export const getProduct = (slug: string) => products.find((product) => product.slug === slug);
export const getCollection = (slug: string) => collections.find((collection) => collection.slug === slug);
