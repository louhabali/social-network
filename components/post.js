import React from 'react';
import Image from 'next/image';
import styles from '../app/page.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHeart, faComment } from '@fortawesome/free-solid-svg-icons';
import Link from 'next/link';

export default function Post({ post, handleLike }) {
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
                <FontAwesomeIcon
                    style={{ cursor: "pointer", color: post.UserInteraction === 1 ? 'red' : 'white' }}
                    onClick={handleLike}
                    icon={faHeart}

                />{post.Likes}

                <Link href={"/post/" + post.Pid} style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '6px',
                    color: '#B0B3B8',
                    textDecoration: 'none',
                    // padding: '8px 12px',
                    borderRadius: '4px',
                    transition: 'background-color 0.2s ease'
                }}>
                    <FontAwesomeIcon
                        icon={faComment}
                        style={{
                            fontSize: '16px',
                            color: '#B0B3B8'
                        }}
                    />
                    <span style={{ fontSize: '14px', fontWeight: '500' }}>
                        {post.NbComment}
                    </span>
                </Link>            </div>
        </section>
    );
};
