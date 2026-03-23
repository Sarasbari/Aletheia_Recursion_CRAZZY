"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import WalletButton from "./WalletButton";

export default function Navbar() {
  const [scrolled, setScrolled] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 50);
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  // Handle hash scrolling after navigation
  useEffect(() => {
    if (pathname === "/" && window.location.hash) {
      const hash = window.location.hash;
      setTimeout(() => {
        const el = document.querySelector(hash);
        if (el) {
          window.scrollTo({
            top: el.getBoundingClientRect().top + window.pageYOffset - 75,
            behavior: "smooth",
          });
        }
      }, 300);
    }
  }, [pathname]);

  const scrollToHash = (hash) => {
    const el = document.querySelector(hash);
    if (el) {
      window.scrollTo({
        top: el.getBoundingClientRect().top + window.pageYOffset - 75,
        behavior: "smooth",
      });
    }
  };

  const handleAnchorClick = (e, hash) => {
    e.preventDefault();
    setMobileOpen(false);

    if (pathname !== "/") {
      router.push("/" + hash);
    } else {
      scrollToHash(hash);
    }
  };

  return (
    <nav className={`navbar${scrolled ? " scrolled" : ""}`} id="navbar">
      <div className="container">
        <Link href="/" className="brand">
          ALETHEIA<span>.</span>
        </Link>
        <ul className={`nav-links${mobileOpen ? " open" : ""}`}>
          <li>
            <a href="/#how-it-works" onClick={(e) => handleAnchorClick(e, "#how-it-works")}>
              process
            </a>
          </li>
          <li>
            <Link href="/verify" onClick={() => setMobileOpen(false)}>
              verify
            </Link>
          </li>
          <li>
            <a href="/#pricing" onClick={(e) => handleAnchorClick(e, "#pricing")}>
              cost
            </a>
          </li>
          <li>
            <WalletButton />
          </li>
          <li>
            <Link href="/anchor" className="nav-cta" onClick={() => setMobileOpen(false)}>
              Launch App
            </Link>
          </li>
        </ul>
        <button
          className="mobile-toggle"
          aria-label="Toggle navigation"
          onClick={() => setMobileOpen(!mobileOpen)}
        >
          {mobileOpen ? "✕" : "☰"}
        </button>
      </div>
    </nav>
  );
}
