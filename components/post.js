import React from 'react';
import Image from 'next/image';
import styles from '../app/page.module.css';

export default function Post ({ post }) {
    return (
        <div className={styles.post}>
            <div className={styles.postHeader}>
                {post.Avatar ? (
                    <Image src={`data:image/png;base64,${post.Avatar}`} alt="" width={32} height={32} />
                ) : (
                    <Image src="/image/no-avatar.png" alt="" width={32} height={32} />
                )}
                <div>
                    <strong>{post.Username}</strong>
                    <span> {post.CreatedAt}</span>
                </div>
            </div>
            <div>{post.Content}</div>
            {post.Image && <Image src={"http://localhost:8080/api/images/"+post.Image} alt="Post Image" width={300} height={300} />}
            <div className={styles.postActions}>
                <button>Like {post.Likes}</button>
                <button>Dislike {post.Dislikes}</button>
                <button>Comment {post.NbComment}</button>
            </div>
        </div>
    );
};

