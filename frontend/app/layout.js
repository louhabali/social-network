"use client"
import "./globals.css";
import { usePathname } from "next/navigation";
import { WebSocketProvider } from "@/components/websocket";
import Header from "@/components/header";
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
        {!hideMenu ? <WebSocketProvider>
           <Header/>
          <main>{children}</main>
        </WebSocketProvider>:
        <main>{children}</main>}
      </body>
    </html>
  );
}