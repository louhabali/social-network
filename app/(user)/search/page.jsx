"use client";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";

export default function SearchResults() {
    const searchParams = useSearchParams();
    const query = searchParams.get("q");
    const [results, setResults] = useState([]);

    useEffect(() => {
        if (!query) return;
        console.log("query", query);
        const formDataToSend = new FormData();
        formDataToSend.append("searchValue", query);
        fetch(`http://localhost:8080/api/searchData`, {
            method: "POST",
            credentials: "include",
            body: formDataToSend,
            
        })
            .then((res) => res.json())
            .then((data) =>{
                if (data){
                setResults(Object.values(data))
                }else{
                    setResults([])
                }
            } )
            .catch((err) => console.error("Error:", err));
    }, [query]);
    console.log("results", results);
    return (
        <div style={{ margin: "80px" }}>
            <h1>results for "{query}"</h1>
            {results.length > 0 ? (
                <ul>
                    {results.map((user, index) => (
                        <li style={{ listStyle:"none",marginBottom:"10px"}} key={index}>
                            <Link style={{ textDecoration:"none",color:"white",display:"flex",alignItems:"center" ,gap:"10px" }} href={`/user/${user.Uid}`}>
                             {user.Avatar ? (
                                <Image src={`data:image/png;base64,${user.Avatar}`} alt="user" style={{ borderRadius:"50%" }} width={32} height={32} />
                               
                                ) : (
                                <Image src="/image/no-avatar.png" alt="" width={32} height={32} />
                                )} {user.Fname} {user.Lname}</Link></li> 
                            
                    ))}
                </ul>
            ) : (
                <p>No users found.</p>
            )}
        </div>
    );
}
