"use client"
import "./globals.css";
import { usePathname } from "next/navigation";
import NavigationMenu from "@/components/NavigationMenu"; 

const metadata = {
  title: "My App",
  description: "A Next.js app",
};

export default function RootLayout({ children }) {
  const pathname = usePathname()
  const hideMenu = pathname === "/login" || pathname === "/register"

  return (
    <html lang="en">
      <body>
        {!hideMenu && <NavigationMenu />} 
        <main>{children}</main> {/* Render page content */}
      </body>
    </html>
  );
}