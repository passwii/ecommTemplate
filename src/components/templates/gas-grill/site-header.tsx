import Link from "next/link";
import { DEFAULT_TEMPLATE } from "@/lib/templates/registry";

export function SiteHeader() {
  const template = DEFAULT_TEMPLATE;
  return (
    <header className="site-header">
      <Link className="wordmark" href="/" aria-label="Ember Field home">EMBER <span>/</span> FIELD</Link>
      <nav aria-label="Main navigation">
        <Link href="/collections/freestanding">Freestanding</Link>
        <Link href="/collections/built-in">Built-in</Link>
        <Link href="/collections/compact">Compact</Link>
        <Link href="/company">Company</Link>
      </nav>
      <Link className="header-cta" href="/inquiries">Request a quote <span aria-hidden="true">↗</span></Link>
      <span className="sr-only">Active storefront template: {template}</span>
    </header>
  );
}
