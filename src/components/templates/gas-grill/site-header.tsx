"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { CONTACT_PATH } from "@/lib/contact";
import { DEFAULT_TEMPLATE } from "@/lib/templates/registry";

export function SiteHeader() {
  const template = DEFAULT_TEMPLATE;
  const [menuOpen, setMenuOpen] = useState(false);
  const menuButtonRef = useRef<HTMLButtonElement>(null);
  const firstMenuLinkRef = useRef<HTMLAnchorElement>(null);

  useEffect(() => {
    if (!menuOpen) return;

    firstMenuLinkRef.current?.focus();
    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setMenuOpen(false);
        menuButtonRef.current?.focus();
      }
    }

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [menuOpen]);

  function closeMenu() {
    setMenuOpen(false);
  }

  return (
    <header className={`site-header${menuOpen ? " menu-open" : ""}`}>
      <Link className="wordmark" href="/" aria-label="Ember Field home">EMBER <span>/</span> FIELD</Link>
      <nav className="desktop-navigation" aria-label="Main navigation">
        <Link href="/collections/freestanding">Freestanding</Link>
        <Link href="/collections/built-in">Built-in</Link>
        <Link href="/collections/compact">Compact</Link>
        <Link href="/#our-approach">Our approach</Link>
      </nav>
      <Link className="header-cta" href={CONTACT_PATH}>Contact <span aria-hidden="true">↗</span></Link>
      <button
        ref={menuButtonRef}
        className="mobile-menu-toggle"
        type="button"
        aria-expanded={menuOpen}
        aria-controls="mobile-main-navigation"
        aria-label={menuOpen ? "Close navigation menu" : "Open navigation menu"}
        onClick={() => setMenuOpen((open) => !open)}
      >
        {menuOpen ? "Close" : "Menu"}
      </button>
      <nav id="mobile-main-navigation" className="mobile-nav" aria-label="Mobile navigation" hidden={!menuOpen}>
        <Link ref={firstMenuLinkRef} href="/collections/freestanding" onClick={closeMenu}>Freestanding</Link>
        <Link href="/collections/built-in" onClick={closeMenu}>Built-in</Link>
        <Link href="/collections/compact" onClick={closeMenu}>Compact</Link>
        <Link href="/#our-approach" onClick={closeMenu}>Our approach</Link>
      </nav>
      <span className="sr-only">Active storefront template: {template}</span>
    </header>
  );
}
