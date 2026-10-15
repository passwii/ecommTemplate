import Link from "next/link";

export function SiteFooter() {
  return (
    <footer className="site-footer">
      <Link className="wordmark footer-mark" href="/">EMBER <span>/</span> FIELD</Link>
      <p>Good food takes time. Make room for it.</p>
      <div className="footer-links">
        <Link href="/collections/freestanding">Explore grills</Link>
        <Link href="/company">Company</Link>
        <a href="mailto:hello@emberandfield.example">Get in touch</a>
        <Link href="/admin/company">Company content</Link>
        <Link href="/admin/products">Product admin</Link>
      </div>
      <small>© 2026 Ember / Field. Sample storefront content. Product data is illustrative, nothing here is connected to checkout or payment.</small>
    </footer>
  );
}
