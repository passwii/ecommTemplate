import { test } from "node:test";
import assert from "node:assert/strict";
import { hasPublishedPrice, toStorefrontCollection, toStorefrontProduct } from "./catalog-view.mjs";

const apiProduct = {
  id: "grill-01",
  slug: "ridge-4-burner",
  name: "Ridge 4-Burner",
  collection: { id: 24, slug: "freestanding", name: "Freestanding grills", description: "Desc", sortOrder: 20 },
  shortDescription: "Room to cook for the whole table.",
  description: "Long description.",
  images: [
    { url: "/images/secondary.jpg", alt: "Secondary", sortOrder: 2 },
    { url: "/images/primary.jpg", alt: "Primary", sortOrder: 1 },
  ],
  specifications: [{ name: "Burners", value: "4" }],
  applications: ["Patios"],
  resources: [{ title: "Spec sheet", url: "/files/spec.pdf" }],
  priceAmount: 2499,
  priceCurrency: "USD",
  featured: true,
  sortOrder: 10,
  contentStatus: "sample",
};

test("toStorefrontProduct maps API payload to the template view model", () => {
  const product = toStorefrontProduct(apiProduct);

  assert.equal(product.collectionSlug, "freestanding");
  assert.equal(product.collection, "Freestanding grills");
  assert.equal(product.summary, "Room to cook for the whole table.");
  assert.equal(product.image, "/images/primary.jpg");
  assert.equal(product.imageAlt, "Primary");
  assert.equal(product.price, 2499);
  assert.equal(product.currency, "USD");
  assert.equal(product.status, "published");
  assert.equal(product.featured, true);
});

test("toStorefrontProduct sorts images by sortOrder", () => {
  const product = toStorefrontProduct(apiProduct);
  assert.deepEqual(
    product.images.map((image) => image.url),
    ["/images/primary.jpg", "/images/secondary.jpg"],
  );
});

test("toStorefrontProduct tolerates missing optional fields", () => {
  const product = toStorefrontProduct({ id: "x", slug: "x", name: "X" });

  assert.deepEqual(product.images, []);
  assert.equal(product.image, "");
  assert.equal(product.price, 0);
  assert.equal(product.currency, "USD");
  assert.equal(product.collectionSlug, "");
  assert.deepEqual(product.applications, []);
});

test("toStorefrontProduct drops blank image urls", () => {
  const product = toStorefrontProduct({
    id: "x",
    slug: "x",
    name: "X",
    images: [
      { url: "   ", alt: "blank", sortOrder: 1 },
      { url: "/images/real.jpg", alt: "real", sortOrder: 2 },
    ],
  });

  assert.deepEqual(
    product.images.map((image) => image.url),
    ["/images/real.jpg"],
  );
  assert.equal(product.image, "/images/real.jpg");
});

test("toStorefrontCollection maps collection payload", () => {
  const collection = toStorefrontCollection({ slug: "built-in", name: "Built-in grills" });
  assert.deepEqual(collection, { slug: "built-in", name: "Built-in grills", description: "" });
});

test("hasPublishedPrice only reports a positive price", () => {
  assert.equal(hasPublishedPrice({ price: 2499 }), true);
  assert.equal(hasPublishedPrice({ price: 0 }), false);
  assert.equal(hasPublishedPrice({}), false);
});
