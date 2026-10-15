"use client";

import { useEffect } from "react";
import { usePathname } from "next/navigation";

/**
 * Marks the document as JS-capable during HTML parse, before first paint, so
 * reveal targets only start hidden when scripting can actually reveal them.
 * Without JS the content renders normally instead of staying invisible.
 */
const ARM_SCRIPT = `document.documentElement.classList.add("js")`;

/**
 * Reveals blocks marked with `data-reveal` once, as they enter the viewport.
 * transform/opacity only, IntersectionObserver only, no scroll listeners.
 *
 * A MutationObserver picks up blocks that appear after the first paint (the
 * company page renders its sections once browser-local content is read), so a
 * late block never gets stranded in its hidden state.
 *
 * Under `prefers-reduced-motion: reduce` the stylesheet resolves every reveal
 * to its end state, so no observer is started at all.
 */
export function Reveal() {
  const pathname = usePathname();

  useEffect(() => {
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
    if (!("IntersectionObserver" in window)) {
      document.documentElement.classList.remove("js");
      return;
    }

    const pending = new Set<HTMLElement>();

    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (!entry.isIntersecting) continue;
          entry.target.classList.add("is-visible");
          observer.unobserve(entry.target);
          pending.delete(entry.target as HTMLElement);
        }
      },
      { rootMargin: "0px 0px -12% 0px", threshold: 0.08 },
    );

    const register = (node: Node) => {
      if (!(node instanceof HTMLElement)) return;
      const targets = node.matches("[data-reveal]") ? [node, ...Array.from(node.querySelectorAll<HTMLElement>("[data-reveal]"))] : Array.from(node.querySelectorAll<HTMLElement>("[data-reveal]"));
      for (const target of targets) {
        if (target.classList.contains("is-visible") || pending.has(target)) continue;
        pending.add(target);
        observer.observe(target);
      }
    };

    register(document.body);

    if (pending.size === 0) {
      document.documentElement.classList.remove("js");
    }

    const mutations = new MutationObserver((records) => {
      for (const record of records) {
        for (const node of record.addedNodes) register(node);
      }
    });
    mutations.observe(document.body, { childList: true, subtree: true });

    return () => {
      observer.disconnect();
      mutations.disconnect();
    };
  }, [pathname]);

  return <script dangerouslySetInnerHTML={{ __html: ARM_SCRIPT }} />;
}
