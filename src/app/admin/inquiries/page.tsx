import type { Metadata } from "next";
import Link from "next/link";
import { redirect } from "next/navigation";
import { fetchAdminInquiries, formatInquiryTimestamp, INQUIRY_STATUSES, statusLabel } from "@/lib/admin/api";

export const metadata: Metadata = {
  title: "Inquiries",
  robots: { index: false, follow: false },
};

type Props = { searchParams: Promise<{ status?: string; q?: string; offset?: string }> };

export default async function AdminInquiriesPage({ searchParams }: Props) {
  const params = await searchParams;
  const status = params.status && INQUIRY_STATUSES.includes(params.status as (typeof INQUIRY_STATUSES)[number]) ? params.status : undefined;
  const query = params.q?.trim() ?? "";
  const offset = Number.parseInt(params.offset ?? "0", 10) || 0;

  const result = await fetchAdminInquiries({ status, q: query || undefined, limit: 50, offset });

  if (!result.ok && result.reason === "unauthorized") redirect("/admin/login");

  const total = result.ok && typeof result.meta?.total === "number" ? result.meta.total : 0;
  const limit = result.ok && typeof result.meta?.limit === "number" ? result.meta.limit : 50;

  return (
    <div className="admin-page page-wrap">
      <nav className="breadcrumbs" aria-label="Breadcrumb"><Link href="/">Home</Link><span>/</span><span>Admin</span></nav>
      <header className="admin-heading">
        <p className="eyebrow">STORE ADMIN</p>
        <h1>Inquiries</h1>
        <p>Visitor enquiries stored in this store&apos;s database.</p>
      </header>

      <form className="admin-filters" method="get">
        <div className="admin-field">
          <label htmlFor="filter-status">Status</label>
          <select id="filter-status" name="status" defaultValue={status ?? ""}>
            <option value="">All statuses</option>
            {INQUIRY_STATUSES.map((value) => (
              <option key={value} value={value}>{statusLabel(value)}</option>
            ))}
          </select>
        </div>
        <div className="admin-field">
          <label htmlFor="filter-q">Search</label>
          <input id="filter-q" name="q" type="search" defaultValue={query} placeholder="Name, email, company or product" />
        </div>
        <button className="button button-dark" type="submit">Apply</button>
      </form>

      {!result.ok ? (
        <p className="admin-empty" role="alert">
          {result.reason === "unavailable"
            ? "The admin API is not available right now."
            : "Inquiries could not be loaded. Please try again."}
        </p>
      ) : result.data.length === 0 ? (
        <p className="admin-empty">No inquiries match this view yet.</p>
      ) : (
        <div className="admin-table-wrap">
          <table className="admin-table">
            <caption className="sr-only">Storefront inquiries</caption>
            <thead>
              <tr>
                <th scope="col">Received</th>
                <th scope="col">Contact</th>
                <th scope="col">Product</th>
                <th scope="col">Country</th>
                <th scope="col">Status</th>
                <th scope="col"><span className="sr-only">Open</span></th>
              </tr>
            </thead>
            <tbody>
              {result.data.map((inquiry) => (
                <tr key={inquiry.id}>
                  <td>{formatInquiryTimestamp(inquiry.createdAt)}</td>
                  <td>
                    <strong>{inquiry.contact.name}</strong>
                    <span className="admin-sub">{inquiry.contact.company || inquiry.contact.email}</span>
                  </td>
                  <td>{inquiry.product.name || "—"}</td>
                  <td>{inquiry.destinationCountry || "—"}</td>
                  <td><span className={`admin-status admin-status-${inquiry.status}`}>{statusLabel(inquiry.status)}</span></td>
                  <td><Link className="text-link" href={`/admin/inquiries/${inquiry.id}`}>Open</Link></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {result.ok && total > limit && (
        <nav className="admin-pagination" aria-label="Pagination">
          {offset > 0 && (
            <Link className="text-link" href={`/admin/inquiries?offset=${Math.max(0, offset - limit)}${status ? `&status=${status}` : ""}${query ? `&q=${encodeURIComponent(query)}` : ""}`}>Newer</Link>
          )}
          <span>{offset + 1}–{Math.min(offset + limit, total)} of {total}</span>
          {offset + limit < total && (
            <Link className="text-link" href={`/admin/inquiries?offset=${offset + limit}${status ? `&status=${status}` : ""}${query ? `&q=${encodeURIComponent(query)}` : ""}`}>Older</Link>
          )}
        </nav>
      )}

      <p className="admin-note">Local MVP administration. Sign-in and inquiry data are served by this store&apos;s own API and database.</p>
    </div>
  );
}
