"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { DEFAULT_TEMPLATE } from "@/lib/templates/registry";

const NAV_LINKS = [
  { href: "/collections/freestanding", label: "Freestanding" },
  { href: "/collections/built-in", label: "Built-in" },
  { href: "/collections/compact", label: "Compact" },
  { href: "/company", label: "Company" },
];

export function SiteHeader() {
  const template = DEFAULT_TEMPLATE;
  const pathname = usePathname();
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    setMenuOpen(false);
  }, [pathname]);

  useEffect(() => {
    if (!menuOpen) return;
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setMenuOpen(false);
    };
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    window.addEventListener("keydown", onKeyDown);
    return () => {
      document.body.style.overflow = previousOverflow;
      window.removeEventListener("keydown", onKeyDown);
    };
  }, [menuOpen]);

  return (
    <header className="site-header">
      <Link className="wordmark" href="/" aria-label="Ember Field home">EMBER <span>/</span> FIELD</Link>
      <nav aria-label="Main navigation">
        {NAV_LINKS.map((link) => <Link key={link.href} href={link.href}>{link.label}</Link>)}
      </nav>
      <div className="header-actions">
        <Link className="header-cta" href="/inquiries">Get in touch</Link>
        <button
          className="menu-toggle"
          type="button"
          aria-expanded={menuOpen}
          aria-controls="site-menu"
          onClick={() => setMenuOpen((open) => !open)}
        >
          <span className="menu-toggle-bars" aria-hidden="true" />
          <span className="sr-only">{menuOpen ? "Close menu" : "Open menu"}</span>
        </button>
      </div>
      <div className="site-menu" id="site-menu" data-open={menuOpen} aria-hidden={!menuOpen}>
        <button className="site-menu-close" type="button" onClick={() => setMenuOpen(false)} aria-label="Close menu">×</button>
        {NAV_LINKS.map((link) => <Link key={link.href} href={link.href}>{link.label}</Link>)}
      </div>
      <span className="sr-only">Active storefront template: {template}</span>
    </header>
  );
}
