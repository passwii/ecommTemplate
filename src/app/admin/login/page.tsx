import type { Metadata } from "next";
import Link from "next/link";
import { AdminLoginForm } from "@/components/admin/login-form";

export const metadata: Metadata = {
  title: "Admin sign in",
  robots: { index: false, follow: false },
};

export default function AdminLoginPage() {
  return (
    <div className="admin-page page-wrap">
      <nav className="breadcrumbs" aria-label="Breadcrumb"><Link href="/">Home</Link><span>/</span><span>Admin</span></nav>
      <header className="admin-heading">
        <p className="eyebrow">STORE ADMIN</p>
        <h1>Sign in</h1>
        <p>Manage storefront inquiries. This area is for store operators only.</p>
      </header>
      <div className="admin-login-card">
        <AdminLoginForm />
      </div>
    </div>
  );
}
