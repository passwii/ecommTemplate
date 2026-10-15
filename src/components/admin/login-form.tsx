"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

export function AdminLoginForm() {
  const router = useRouter();
  const [status, setStatus] = useState<"idle" | "submitting" | "error">("idle");
  const [message, setMessage] = useState("");

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);

    setStatus("submitting");
    setMessage("");

    let response: Response;
    try {
      response = await fetch("/api/admin/sessions", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          email: String(data.get("email") ?? ""),
          password: String(data.get("password") ?? ""),
        }),
      });
    } catch {
      setStatus("error");
      setMessage("We could not reach the store just now. Please try again.");
      return;
    }

    if (response.ok) {
      form.reset();
      router.replace("/admin/inquiries");
      router.refresh();
      return;
    }

    setStatus("error");
    if (response.status === 401) {
      setMessage("The email or password is incorrect.");
    } else if (response.status === 503) {
      setMessage("Sign-in is not available right now.");
    } else {
      setMessage("Sign-in failed. Please try again.");
    }
  }

  return (
    <form className="admin-login-form" onSubmit={handleSubmit}>
      {status === "error" && <p className="admin-error" role="alert">{message}</p>}

      <div className="admin-field">
        <label htmlFor="admin-email">Email</label>
        <input id="admin-email" name="email" type="email" required autoComplete="username" />
      </div>

      <div className="admin-field">
        <label htmlFor="admin-password">Password</label>
        <input id="admin-password" name="password" type="password" required autoComplete="current-password" />
      </div>

      <button className="button button-dark admin-submit" type="submit" disabled={status === "submitting"}>
        {status === "submitting" ? "Signing in…" : "Sign in"}
      </button>
    </form>
  );
}
