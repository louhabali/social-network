"use client"
import {useEffect, useState } from "react";
import Link from "next/link";
import Image from "next/image";
export default function Sugg() {
    const [suggs , setSuggs ]= useState([])
    useEffect(()=>{
    fetch("http://localhost:8080/api/suggestions", { credentials: "include" }).then(d=>d.json()).then(dd=>{
        setSuggs(dd)
    })
},[])
console.log(suggs);

    return (
        <div>
            <div style={{marginTop:"100px"}}>suggestions</div>
            <ul style={{listStyle:"none"}}>
            {suggs.map((v) => (
                 <li style={{ listStyle:"none",marginBottom:"10px"}} key={v.id}>
                 <Link style={{ textDecoration:"none",color:"white",display:"flex",alignItems:"center" ,gap:"10px" }} href={`/profile/${v.id}`}>
                  {v.Avatar ? (
                     <Image src={`data:image/png;base64,${v.Avatar}`} alt="user" style={{ borderRadius:"50%" }} width={32} height={32} />

                     ) : (
                     <Image src="/image/no-avatar.png" alt="" width={32} height={32} />
                     )} {v.fname} {v.lname}</Link></li> 
            ))}
        </ul>
        </div>

    )
}