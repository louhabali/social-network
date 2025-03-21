"use client";
import React, { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import styles from "./page.module.css";
import Post from "@/components/post";
import PostForm from "@/components/postForm";
import { websoketProvider } from "./oldwebSocket";

const HomePage = () => {
    
    const [posts, setPosts] = useState([]);
    const [userInfo, setUserInfo] = useState(null);
    const [offset, setOffset] = useState(0);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const observer = useRef(null);
    const router = useRouter();
    useEffect(() => {
        websoketProvider.connect().
            then(res => {
                console.log(res);
                window.ws = res
            }).catch(error => {
                console.error("WebSocket error:", error);
            });
    }, [])

    // posts
    const fetchPosts = async () => {
        if (loading || !hasMore) return;

        setLoading(true);
        try {
            const response = await fetch(`http://localhost:8080/api/posts?offset=${offset}`, {
                credentials: "include",
            });

            if (!response.ok) {
                setHasMore(false);
                throw Error("gg")
                return;
            }

            const data = await response.json();
            if (data === null) {
                setHasMore(false);
                return
            }
            if (data.length < 4) {
                setHasMore(false);
            };
            setPosts((prevPosts) => [...prevPosts, ...data]);
            setOffset((prevOffset) => prevOffset + data?.length);
        } catch (error) {
            
            router.push("/login");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetch("http://localhost:8080/api/userdata", { credentials: "include" })
            .then((response) => response.json())
            .then((data) => setUserInfo(data))
            .catch((error) => console.error("Error fetching user data:", error));
    }, []);


    useEffect(() => {
        if (loading || !hasMore) return;

        observer.current = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting) {
                    fetchPosts();
                }
            },
            { threshold: 1.0 }
        );

        if (observer.current && document.getElementById("load-more")) {
            observer.current.observe(document.getElementById("load-more"));
        }

        return () => observer.current?.disconnect();
    }, [fetchPosts, loading, hasMore]);

    
    
    console.log(posts)
    return (
        <div className={styles.container}>
            {/* Remove the Header component from here */}

            <main className={styles.main}>
                <div className={styles.feed}>
                    <PostForm userInfo={userInfo} setPosts={setPosts} setOffset={setOffset} />
                </div>

                {posts.map((post) => (
                    <Post key={post.Pid} post={post} />
                ))}

                {hasMore && <div id="load-more" style={{ height: "10px" }} />}
                {loading && <p className={styles.loading}>Loading...</p>}
            </main>
        </div>
    );
};

export default HomePage;