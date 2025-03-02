"use client"
import {useEffect } from "react";
export default function Sugg() {
useEffect(()=>{
    fetch("http://localhost:8080/api/suggestions", { credentials: "include" }).then(d=>d.json()).then(dd=>{
        console.log(dd);
    })
},[])
    return (
        
        <div>suggestions</div>
    )
}