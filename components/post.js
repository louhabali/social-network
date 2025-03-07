'use client';
import React from 'react';
import Image from 'next/image';
import styles from '../app/page.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHeart, faComment } from '@fortawesome/free-solid-svg-icons';
import Link from 'next/link';

export default function Post({ post, handleLike }) {
    const handleLikeClick = () => {
        handleLike(Number(post.Pid));
    };
    console.log(post.NbComment);
    
    return (
        <section id={post.Pid} key={post.Pid} className={styles.post}>
            <div className={styles.postHeader}>
                {post.Avatar ? (
                    <Image style={{ borderRadius: "50%" }} src={`data:image/png;base64,${post.Avatar}`} alt="" width={32} height={32} />
                ) : (
                    <Image style={{ borderRadius: "50%" }} src="/image/no-avatar.png" alt="" width={32} height={32} />
                )}
                <div>
                    <strong>{post.Username}</strong>
                </div>
                <div><span>{post.CreatedAt}</span></div>
            </div>
            <div>{post.Content}</div>
            <Link href={"/post/" + post.Pid}>
                {post.Image && <Image src={"http://localhost:8080/api/images/" + post.Image} alt="" width={400} height={200} />}
            </Link>
            <div className={styles.postActions}>
                <div style={{ flexDirection: "column", textAlign: "center" }}>
                    <FontAwesomeIcon
                        style={{ cursor: "pointer", color: post.UserInteraction === 1 ? 'red' : 'white' }}
                        onClick={handleLikeClick}
                        icon={faHeart} />
                    <div>
                        {post.Likes} {post.Likes !== 1 ? "likes" : "like"}
                    </div>
                </div>
                <div style={{ flexDirection: "column", textAlign: "center" }}>
                    <Link style={{ textDecoration: "none", color: "white" }} href={"/post/" + post.Pid}>
                        <FontAwesomeIcon icon={faComment} />
                    </Link>
                    <div>
                        {post.NbComment === 0 ? "no" : post.NbComment}{post.NbComment != 1?" comments":" comment"}
                    </div>
                </div>
            </div>
        </section>
    );
}

