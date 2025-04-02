'use client';
import React, { useState, useContext } from 'react';
import Image from 'next/image';
import styles from '../app/page.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHeart, faComment } from '@fortawesome/free-solid-svg-icons';
import { WebSocketContext } from "@/context/WebSocketContext";
import Link from 'next/link';

export default function Post({ post }) {
    const [likes, setLikes] = useState(post.Likes);
    const [userInteraction, setUserInteraction] = useState(post.UserInteraction);
    const { showNotification } = useContext(WebSocketContext);

    const handleLike = async (postId) => {
        let likeStatus = userInteraction === 1 ? 0 : 1;

        const response = await fetch("http://localhost:8080/api/like", {
            method: "POST",
            body: JSON.stringify({ postId, likeStatus }),
            credentials: "include",
        });

        if (response.ok) {
            setUserInteraction(likeStatus);
            setLikes((prevLikes) => prevLikes + (likeStatus ? 1 : -1));
        }else {
            showNotification("Error liking the post", "error");
        }
    };

    return (
        <section id={post.Pid} className={styles.post}>
            <div className={styles.postHeader}>
                <Link href={`/profile/${post.Uid}`}>
                    {post.Avatar ? (
                        <Image style={{ cursor: "pointer", borderRadius: "50%" }} src={`http://localhost:8080/api/images/${post.Avatar}`} alt="" width={32} height={32} />
                    ) : (
                        <Image style={{ cursor: "pointer", borderRadius: "50%" }} src="/image/no-avatar.png" alt="" width={32} height={32} />
                    )}
                </Link>
                <div>
                    <strong>{post.Username}</strong>
                </div>
                <div><span>{post.CreatedAt}</span></div>
            </div>
            <div>{post.Content}</div>
            <Link href={`/post/${post.Pid}`}>
                {post.Image && <Image src={`http://localhost:8080/api/images/${post.Image}`} alt="" width={400} height={200} />}
            </Link>
            <div className={styles.postActions}>
                <div style={{ flexDirection: "column", textAlign: "center" }}>
                    <FontAwesomeIcon
                        style={{ cursor: "pointer", color: userInteraction === 1 ? 'red' : 'white' }}
                        onClick={() => handleLike(post.Pid)}
                        icon={faHeart}
                    />
                    <div>
                        {likes} {likes !== 1 ? "likes" : "like"}
                    </div>
                </div>
                <div style={{ flexDirection: "column", textAlign: "center" }}>
                    <Link style={{ textDecoration: "none", color: "white" }} href={`/post/${post.Pid}`}>
                        <FontAwesomeIcon icon={faComment} />
                    </Link>
                    <div>
                        {post.NbComment === 0 ? "no" : post.NbComment} {post.NbComment !== 1 ? "comments" : "comment"}
                    </div>
                </div>
            </div>
        </section>
    );
}
