"use client";

import Link from "next/link";

export default function ErrorPage({ error }) {
  return (
    <div style={{ display: "flex", flexDirection: "column", alignItems: "center" , justifyContent: "center", height: "100vh"}}>
      <h1>Error</h1>
      <p>{error.message || "Something went wrong."}</p>
      <Link href="/">
        Go back to home
      </Link>
    </div>
  );
}
