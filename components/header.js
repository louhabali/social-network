"use client";
import Link from "next/link";
import '../app/styles/header.css';
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";

export default function Header() {
    const [searchData, setSearchData] = useState('');
    const router = useRouter();
    if (document.getElementById("s")){
        document.getElementById("s").innerHTML=""
    }
    function handleSearch(e) {
        setSearchData(e.target.value);
    }

    function handleKeyDown(e) {
        if (e.key === "Enter") {
            // reset search input
             setSearchData("");
             document.getElementById("s").innerHTML=""
            if (searchData.trim()) {
                router.push(`/search?q=${searchData}`);
            }
        }
    }
    useEffect(() => {
            
            const formDataToSend = new FormData();
            formDataToSend.append("searchValue", searchData);
            fetch('http://localhost:8080/api/searchData', {
                method: 'POST',
                credentials: "include",
                body: formDataToSend,
            })
                .then(response => {
                   return response.json();
                    
                })
                .then(data => {
                    if (data){
                        Object.entries(data).forEach(x=>{
                                if (searchData){
                                    document.getElementById("s").innerHTML+= `
                                    <p>${x[1].Fname+" "+x[1].Lname}</p>
                                    ` 
                                }else{
                                    document.getElementById("s").innerHTML=""
                                }
                            
                          
                        })
                    }
               
                     //console.log(users);
                     
                }
                )
                .catch(error => console.error('error', error));
    }, [searchData])

    return (
        <header className="header">
            <input
                type="text"
                placeholder="search for users..."
                value={searchData}
                onChange={handleSearch}
                onKeyDown={handleKeyDown} 
            />
            <div className="profile">
                <Link href="/api/profile">profile</Link>
            </div>
            <div id="s" style={{marginTop:"45px",opacity:"50%",top:"0", background:"rgb(36, 35, 35)", position:"absolute" ,color:"white"}}></div>
        </header>
    );
}
