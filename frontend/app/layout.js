"use client"
import "./globals.css";
import { usePathname } from "next/navigation";
// import NavigationMenu from "@/components/NavigationMenu";
import Header from "@/components/header";
import { useEffect, useState } from "react";
const metadata = {
  title: "My App",
  description: "A Next.js app",
};

export default function RootLayout({ children }) {
  const [u , setUser] = useState([])
  useEffect(() => {
             fetch('http://localhost:8080/api/userdata', {
                 credentials: "include",
             }).then(data => data.json()).then(d => {
                 setUser(d)
             })
     }, []);
  const pathname = usePathname()
  const hideMenu = pathname === "/login" || pathname === "/register"

  return (
    <html lang="en">
      <body>
        {/* {!hideMenu && <NavigationMenu />} */}
        {!hideMenu &&  <Header userID={u.ID} Avatar={u.Avatar} fname={u.Firstname} lname={u.Lastname} />}
        <main>{children}</main>
        
      </body>
    </html>
  );
}