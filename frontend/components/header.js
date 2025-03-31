"use client";
import React from 'react';
import Link from 'next/link';
import styles from './Header.module.css';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from "react";
import Image from 'next/image';


const Header = () => {

    const router = useRouter();
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [u, setUser] = useState([])

    useEffect(() => {
        fetch('http://localhost:8080/api/userdata', {
            credentials: "include",
        }).then(data => data.json()).then(d => {
            setUser(d)
        })
    }, []);

    const handleLogout = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/logout', {
                method: 'POST',
                credentials: 'include',
            });

            if (response.ok) {
                router.push('/login');
            }
        } catch (error) {
            console.error('Logout error:', error);
        }
    };

    const [searchData, setSearchData] = useState("");
    const [searchResults, setSearchResults] = useState([]);

    function handleSearch(e) {
        setSearchData(e.target.value);
    }
    function handleKeyDown(e) {
        if (e.key === "Enter" && searchData.trim()) {
            setSearchData("");
            setSearchResults([]);
            router.push(`/search?q=${encodeURIComponent(searchData)}`);
        }
    }
    useEffect(() => {
        if (!searchData.trim()) {
            setSearchResults([]);
            return;
        }
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
                if (data) {
                    if (data) {
                        setSearchResults(Object.values(data));
                    }
                    console.log(data);
                }


            }
            )
            .catch(error => console.error('error', error));
    }, [searchData])
    return (
        <>
            {/* Desktop Header */}
            <header className={styles.header}>
                <div className={styles.headerLeft}>
                    <Link href="/" className={styles.logo}>
                        Connectify
                    </Link>
                    <div className={styles.searchContainer} >
                        <input
                            type="text"
                            placeholder="Search Connectify..."
                            className={styles.searchInput}
                            value={searchData}
                            onChange={handleSearch}
                            onKeyDown={handleKeyDown}
                        />
                        {searchResults.length > 0 && (
                            <div style={{
                                marginTop: "45px", top: "0",
                                width: "206px",
                                background: "rgb(0, 0, 0)",
                                paddingLeft: "5px",
                                position: "absolute",
                                color: "white"
                            }}>
                                {searchResults.map((user, index) => (
                                    <p key={index}>
                                        <Link style={{ textDecoration:"none",color:"white",display:"flex",alignItems:"center" ,gap:"10px" }} href={`/profile/${user.Uid}`}>
                                        {user.Avatar ? (
                                        <Image src={`http://localhost:8080/api/images/${user.Avatar}`} alt="user" style={{ borderRadius: "50%" }} width={32} height={32} />

                                    ) : (
                                        <Image src="/image/no-avatar.png" alt="" width={32} height={32} />
                                    )} {user.Fname} {user.Lname}</Link></p>
                                ))}
                            </div>
                        )}
                        <div id="s" style={{ marginTop: "45px", top: "0", width: "206px", background: "rgb(0, 0, 0)", paddingLeft: "5px", position: "absolute", color: "white" }}></div>

                        <button className={styles.searchButton}>
                            {/* <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.searchIcon}>
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                            </svg> */}
                        </button>
                    </div>
                </div>

                {/* Desktop Navigation */}
                <div className={styles.headerCenter}>
                    <nav className={styles.mainNav}>
                        <Link href="/" className={styles.navItem}>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.navIcon}>
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                            </svg>
                        </Link>
                        <Link href="/groups" className={styles.navItem}>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.navIcon}>
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                            </svg>
                        </Link>
                        <Link href="/chat" className={styles.navItem}>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.navIcon}>
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                            </svg>
                        </Link>
                        <Link href="/notifications" className={styles.navItem}>
                            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.navIcon}>
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                            </svg>
                        </Link>
                    </nav>
                </div>

                <div className={styles.headerRight}>
                    <div className={styles.userMenu} onClick={() => setIsDropdownOpen(!isDropdownOpen)}>
                        {u.Avatar ? (
                            <img
                                src={`http://localhost:8080/api/images/${u.Avatar}`} alt="user"
                                className={styles.userAvatar}
                                width={20} height={20}
                            />
                        ) : (
                            <Image src="/image/no-avatar.png" alt="" className={styles.userAvatar} width={20} height={20} />

                        )}<span>{u.Firstname} {u.Lastname}</span>

                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            className={`${styles.dropdownIcon} ${isDropdownOpen ? styles.dropdownIconOpen : ''}`}
                        >
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>

                        {isDropdownOpen && (
                            <div className={styles.dropdown}>
                                <Link href={`/profile/${u.ID}`} className={styles.dropdownItem}>
                                    Profile
                                </Link>
                                {/* <Link href="/settings" className={styles.dropdownItem}>
                                    Settings
                                </Link> */}
                                <button onClick={handleLogout} className={styles.dropdownItem}>
                                    Logout
                                </button>
                            </div>
                        )}
                    </div>
                </div>
            </header>

            {/* Mobile Bottom Navigation */}
            <nav className={styles.mobileNav}>
                <div className={styles.mobileNavInner}>
                    <Link href="/" className={styles.mobileNavItem}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.mobileNavIcon}>
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                        </svg>
                        <span className={styles.mobileNavLabel}>Home</span>
                    </Link>

                    <Link href="/groups" className={styles.mobileNavItem}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.mobileNavIcon}>
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                        </svg>
                        <span className={styles.mobileNavLabel}>Groups</span>
                    </Link>

                    <Link href="/messages" className={styles.mobileNavItem}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.mobileNavIcon}>
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                        </svg>
                        <span className={styles.mobileNavLabel}>Chat</span>
                    </Link>

                    <Link href="/notifications" className={styles.mobileNavItem}>
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" className={styles.mobileNavIcon}>
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                        </svg>
                        <span className={styles.mobileNavLabel}>Notifications</span>
                    </Link>

                    <Link href="/profile" className={styles.mobileNavItem}>
                        <div className={styles.userMenuIcon}>
                            {u.Avatar ? (
                                <img
                                    src={`http://localhost:8080/api/images/${u.Avatar}`}
                                    alt={u.Firstname}
                                    className={styles.mobileUserAvatar}
                                    width={20} height={20}
                                />
                            ) : (
                                <Image src="/image/no-avatar.png" alt="" className={styles.mobileUserAvatar} width={20} height={20} />

                            )}
                        </div>
                        <span className={styles.mobileNavLabel}>Profile</span>
                    </Link>
                </div>
            </nav>
        </>
    );
};


export default Header;