"use client";

import { useSearchParams } from "next/navigation";
import { useEffect, useState, Suspense } from "react";
import Image from "next/image";
import Link from "next/link";

const SearchResults = () => {
    const params = useSearchParams();
    const query = params.get("q");
    const [results, setResults] = useState([]);

    useEffect(() => {
        console.log(query);

        if (!query) return;

        const formDataToSend = new FormData();
        formDataToSend.append("searchValue", query);

        fetch(`http://localhost:8080/api/searchData`, {
            method: "POST",
            credentials: "include",
            body: formDataToSend,
        })
            .then((res) => res.json())
            .then((data) => {
                if (data) {
                    setResults(Object.values(data));
                } else {
                    setResults([]);
                }
            })
            .catch((err) => console.error("Error:", err));
    }, [query]);

    console.log("results", results);

    return (
        <div style={{ margin: "80px" }}>
            <h1>Results for "{query}"</h1>
            {results.length > 0 ? (
                <ul>
                    {results.map((user, index) => (
                        <li style={{ listStyle: "none", marginBottom: "10px" }} key={index}>
                            <Link
                                style={{
                                    textDecoration: "none",
                                    color: "white",
                                    display: "flex",
                                    alignItems: "center",
                                    gap: "10px",
                                }}
                                href={`/profile/${user.Uid}`}
                            >
                                {user.Avatar ? (
                                    <img
                                        src={`http://localhost:8080/api/images/${user.Avatar}`}
                                        alt="user"
                                        style={{ borderRadius: "50%" }}
                                        width={32}
                                        height={32}
                                    />
                                ) : (
                                    <img src="/image/no-avatar.png" alt="" width={32} height={32} />
                                )}
                                {user.Fname} {user.Lname}
                            </Link>
                        </li>
                    ))}
                </ul>
            ) : (
                <p>No users found.</p>
            )}
        </div>
    );
};

// Wrap the component inside Suspense to handle CSR properly
const SearchResultsWithSuspense = () => {
    return (
        <Suspense fallback={<div>Loading...</div>}>
            <SearchResults />
        </Suspense>
    );
};

export default SearchResultsWithSuspense;
