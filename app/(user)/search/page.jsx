"use client";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

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

    return (
        <div style={{ margin: "80px" }}>
            <h1>results for "{query}"</h1>
            {results.length > 0 ? (
                <ul>
                    {results.map((user, index) => (
                        <li key={index}>{user.Fname} {user.Lname}</li>
                    ))}
                </ul>
            ) : (
                <p>No users found.</p>
            )}
        </div>
    );
}
