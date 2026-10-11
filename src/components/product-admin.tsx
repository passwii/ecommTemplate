"use client";

import Image from "next/image";
import Link from "next/link";
import { useCallback, useEffect, useMemo, useState, type FormEvent } from "react";
import { collections, products, type Product } from "@/lib/catalog/mock-data";
import { draftToProduct, getStoredProducts, productToDraft, resetStoredProducts, saveStoredProducts, validateProductDraft, type ProductDraft } from "@/lib/catalog/product-admin";

const emptyDraft = (): ProductDraft => ({
  id: `product-${Date.now()}`,
  slug: "",
  name: "",
  category: collections[0].slug,
  shortDescription: "",
  description: "",
  images: [{ url: "", alt: "", sortOrder: 0 }],
  specifications: [{ name: "", value: "" }],
  price: 0,
  currency: "USD",
  status: "draft",
  featured: false,
  sortOrder: 100,
});

export function ProductAdmin() {
  const [items, setItems] = useState<Product[]>(products);
  const [draft, setDraft] = useState<ProductDraft | null>(null);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [notice, setNotice] = useState("");
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    setItems(getStoredProducts());
    setLoaded(true);
  }, []);

  const sortedItems = useMemo(() => [...items].sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name)), [items]);

  const persist = useCallback((next: Product[], message: string) => {
    setItems(next);
    try {
      saveStoredProducts(next);
      setNotice(message);
    } catch {
      setNotice("Could not save in this browser. Check available storage and try again.");
    }
  }, []);

  function startCreate() {
    setErrors({});
    setNotice("");
    setDraft(emptyDraft());
  }

  function startEdit(product: Product) {
    setErrors({});
    setNotice("");
    setDraft(productToDraft(product));
  }

  function updateDraft<K extends keyof ProductDraft>(key: K, value: ProductDraft[K]) {
    setDraft((current) => current ? { ...current, [key]: value } : current);
  }

  function saveDraft(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!draft) return;
    const validation = validateProductDraft(draft);
    const duplicateSlug = items.some((item) => item.slug === draft.slug && item.id !== draft.id);
    if (duplicateSlug) validation.slug = "This product URL is already in use.";
    setErrors(validation);
    if (Object.keys(validation).length) return;
    const existing = items.find((item) => item.id === draft.id);
    const saved = draftToProduct(draft, existing);
    const next = existing ? items.map((item) => item.id === saved.id ? saved : item) : [...items, saved];
    persist(next, `Saved ${saved.name} as ${saved.status}.`);
    setDraft(null);
  }

  function toggleStatus(product: Product) {
    const status: Product["status"] = product.status === "published" ? "draft" : "published";
    if (status === "published") {
      const validation = validateProductDraft(productToDraft(product));
      const firstError = Object.values(validation)[0];
      if (firstError) {
        setNotice(`Cannot publish ${product.name}: ${firstError}`);
        return;
      }
    }
    const next = items.map((item): Product => item.id === product.id ? { ...item, status } : item);
    persist(next, `${product.name} ${status === "published" ? "published" : "unpublished"}.`);
  }

  function resetCatalog() {
    resetStoredProducts();
    setItems(products);
    setDraft(null);
    setNotice("The browser catalog has been reset to the original sample products.");
  }

  function addImage() {
    if (!draft) return;
    updateDraft("images", [...draft.images, { url: "", alt: "", sortOrder: draft.images.length }]);
  }

  function updateImage(index: number, key: "url" | "alt", value: string) {
    if (!draft) return;
    updateDraft("images", draft.images.map((image, imageIndex) => imageIndex === index ? { ...image, [key]: value } : image));
  }

  function removeImage(index: number) {
    if (!draft || draft.images.length <= 1) return;
    updateDraft("images", draft.images.filter((_, imageIndex) => imageIndex !== index).map((image, sortOrder) => ({ ...image, sortOrder })));
  }

  function moveImage(index: number, offset: -1 | 1) {
    if (!draft) return;
    const target = index + offset;
    if (target < 0 || target >= draft.images.length) return;
    const images = [...draft.images];
    [images[index], images[target]] = [images[target], images[index]];
    updateDraft("images", images.map((image, sortOrder) => ({ ...image, sortOrder })));
  }

  function addSpecification() {
    if (!draft) return;
    updateDraft("specifications", [...draft.specifications, { name: "", value: "" }]);
  }

  function updateSpecification(index: number, key: "name" | "value", value: string) {
    if (!draft) return;
    updateDraft("specifications", draft.specifications.map((specification, specIndex) => specIndex === index ? { ...specification, [key]: value } : specification));
  }

  function removeSpecification(index: number) {
    if (!draft) return;
    updateDraft("specifications", draft.specifications.filter((_, specIndex) => specIndex !== index));
  }

  if (!loaded) return <div className="admin-loading" aria-live="polite">Loading local catalog…</div>;

  return (
    <div className="admin-shell">
      <header className="admin-topbar">
        <Link href="/" className="admin-brand">EMBER <span>/</span> FIELD <small>STUDIO</small></Link>
        <nav aria-label="Admin navigation"><Link href="/admin/products" aria-current="page">Products</Link><Link href="/">View storefront ↗</Link></nav>
      </header>
      <main className="admin-main">
        <div className="admin-heading">
          <div><p className="admin-eyebrow">CATALOG MANAGEMENT · GAS-GRILL TEMPLATE</p><h1>Products</h1><p>Curate what appears in your storefront.</p></div>
          <button className="admin-primary" onClick={startCreate}>＋ <span>New product</span></button>
        </div>
        <div className="admin-callout"><span className="admin-callout-mark">i</span><p>This demo saves to this browser only. It has no sign-in or server persistence; changes won’t sync across devices and clearing browser data resets them.</p></div>
        {notice && <p className="admin-notice" role="status">{notice}</p>}
        <section className="admin-table-wrap" aria-label="Product catalog">
          <div className="admin-table-head"><span>PRODUCT <small>{items.length} TOTAL</small></span><span>COLLECTION</span><span>PRICE</span><span>STATUS</span><span>ACTIONS</span></div>
          {sortedItems.map((product) => (
            <article className="admin-row" key={product.id}>
              <div className="admin-product-cell"><div className="admin-thumb"><Image src={product.images[0]?.url || product.image} alt="" width={48} height={48} unoptimized /></div><div><strong>{product.name}</strong><small>/{product.slug}</small></div></div>
              <div className="admin-collection-cell">{product.collection}</div>
              <div className="admin-price-cell">{new Intl.NumberFormat("en-US", { style: "currency", currency: product.currency, maximumFractionDigits: 0 }).format(product.price)}</div>
              <div><span className={`admin-status ${product.status}`}><i />{product.status}</span></div>
              <div className="admin-actions"><button onClick={() => startEdit(product)}>Edit</button><button className="admin-publish" onClick={() => toggleStatus(product)}>{product.status === "published" ? "Unpublish" : "Publish"}</button></div>
            </article>
          ))}
          {!sortedItems.length && <div className="admin-empty"><strong>No products yet.</strong><p>Create a draft to start building your catalog.</p><button className="admin-secondary" onClick={startCreate}>Create first product</button></div>}
        </section>
        <div className="admin-list-footer"><span>Showing {items.length} products · Sorted by display order</span><button onClick={resetCatalog}>Reset sample catalog</button></div>
      </main>

      {draft && <div className="admin-modal-backdrop" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget) setDraft(null); }}>
        <section className="admin-modal" role="dialog" aria-modal="true" aria-labelledby="product-form-title">
          <header className="admin-modal-header"><div><p className="admin-eyebrow">PRODUCT DETAILS</p><h2 id="product-form-title">{items.some((item) => item.id === draft.id) ? "Edit product" : "New product"}</h2></div><button aria-label="Close editor" className="admin-close" onClick={() => setDraft(null)}>×</button></header>
          <form onSubmit={saveDraft} noValidate>
            <div className="admin-form-body">
              <fieldset className="admin-fieldset"><legend>Basics</legend>
                <label className="admin-field">Product name<input autoFocus value={draft.name} onChange={(event) => updateDraft("name", event.target.value)} placeholder="e.g. Ridge 4-Burner" />{errors.name && <small className="admin-error">{errors.name}</small>}</label>
                <label className="admin-field">URL slug<div className="admin-slug-wrap"><span>/products/</span><input value={draft.slug} onChange={(event) => updateDraft("slug", event.target.value.toLowerCase().replace(/\s+/g, "-"))} placeholder="ridge-4-burner" /></div>{errors.slug && <small className="admin-error">{errors.slug}</small>}</label>
                <label className="admin-field">Collection<select value={draft.category} onChange={(event) => updateDraft("category", event.target.value)}>{collections.map((collection) => <option key={collection.slug} value={collection.slug}>{collection.name}</option>)}</select>{errors.category && <small className="admin-error">{errors.category}</small>}</label>
                <label className="admin-field">Short description<input value={draft.shortDescription} onChange={(event) => updateDraft("shortDescription", event.target.value)} placeholder="A concise line for product cards" />{errors.shortDescription && <small className="admin-error">{errors.shortDescription}</small>}</label>
                <label className="admin-field">Description<textarea rows={4} value={draft.description} onChange={(event) => updateDraft("description", event.target.value)} placeholder="Tell customers what makes this grill special." />{errors.description && <small className="admin-error">{errors.description}</small>}</label>
              </fieldset>
              <fieldset className="admin-fieldset"><legend>Images</legend><p className="admin-field-hint">Use a path to an image in the storefront public folder or a public image URL. First image is the cover.</p>
                {draft.images.map((image, index) => <div className="admin-repeater" key={`image-${index}`}><div className="admin-repeater-title"><span>IMAGE {String(index + 1).padStart(2, "0")}</span><div className="admin-repeater-actions"><button type="button" aria-label={`Move image ${index + 1} up`} disabled={index === 0} onClick={() => moveImage(index, -1)}>↑</button><button type="button" aria-label={`Move image ${index + 1} down`} disabled={index === draft.images.length - 1} onClick={() => moveImage(index, 1)}>↓</button>{draft.images.length > 1 && <button type="button" onClick={() => removeImage(index)}>Remove</button>}</div></div><label className="admin-field">Image URL<input value={image.url} onChange={(event) => updateImage(index, "url", event.target.value)} placeholder="/images/grill-example.jpg" /></label><label className="admin-field">Alt text<input value={image.alt} onChange={(event) => updateImage(index, "alt", event.target.value)} placeholder="Describe the image for accessibility" /></label></div>)}
                {errors.images && <small className="admin-error">{errors.images}</small>}<button className="admin-add-row" type="button" onClick={addImage}>＋ Add image</button>
              </fieldset>
              <fieldset className="admin-fieldset"><legend>Specifications</legend><p className="admin-field-hint">Add product facts such as burners, cooking area, fuel and finish.</p>
                {draft.specifications.map((specification, index) => <div className="admin-spec-row" key={`spec-${index}`}><input aria-label={`Specification ${index + 1} name`} value={specification.name} onChange={(event) => updateSpecification(index, "name", event.target.value)} placeholder="Name" /><input aria-label={`Specification ${index + 1} value`} value={specification.value} onChange={(event) => updateSpecification(index, "value", event.target.value)} placeholder="Value" /><button type="button" aria-label={`Remove specification ${index + 1}`} onClick={() => removeSpecification(index)}>×</button></div>)}<button className="admin-add-row" type="button" onClick={addSpecification}>＋ Add specification</button>
              </fieldset>
              <fieldset className="admin-fieldset"><legend>Pricing & publishing</legend>
                <div className="admin-two-col"><label className="admin-field">Price<input type="number" min="0.01" step="0.01" value={draft.price || ""} onChange={(event) => updateDraft("price", Number(event.target.value))} placeholder="2499" />{errors.price && <small className="admin-error">{errors.price}</small>}</label><label className="admin-field">Currency<select value={draft.currency} onChange={(event) => updateDraft("currency", event.target.value)}><option value="USD">USD — US Dollar</option><option value="CAD">CAD — Canadian Dollar</option><option value="EUR">EUR — Euro</option><option value="GBP">GBP — Pound Sterling</option></select></label></div>
                <label className="admin-field">Compare-at price <span className="admin-optional">OPTIONAL</span><input type="number" min="0" step="0.01" value={draft.compareAtPrice ?? ""} onChange={(event) => updateDraft("compareAtPrice", event.target.value ? Number(event.target.value) : undefined)} placeholder="Leave blank to hide" /></label>
                <label className="admin-field">SKU <span className="admin-optional">OPTIONAL</span><input value={draft.sku ?? ""} onChange={(event) => updateDraft("sku", event.target.value)} placeholder="e.g. EF-RDG-04" /></label>
                <div className="admin-two-col"><label className="admin-field">Display order<input type="number" value={draft.sortOrder} onChange={(event) => updateDraft("sortOrder", Number(event.target.value))} /></label><label className="admin-field">Status<select value={draft.status} onChange={(event) => updateDraft("status", event.target.value as ProductDraft["status"])}><option value="draft">Draft</option><option value="published">Published</option></select></label></div>
                <label className="admin-checkbox"><input type="checkbox" checked={draft.featured} onChange={(event) => updateDraft("featured", event.target.checked)} /><span>Feature this product on the storefront</span></label>
              </fieldset>
            </div>
            <footer className="admin-modal-footer"><button className="admin-secondary" type="button" onClick={() => setDraft(null)}>Cancel</button><button className="admin-primary" type="submit">Save {draft.status === "published" ? "product" : "draft"}</button></footer>
          </form>
        </section>
      </div>}
    </div>
  );
}
