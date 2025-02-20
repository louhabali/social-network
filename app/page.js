"use client"
import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import styles from './page.module.css';

const HomePage = () => {
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState('');
    const [selectedImage, setSelectedImage] = useState(null);

    useEffect(() => {
        fetch('http://localhost:8080/api/posts')
            .then(response => response.json())
            .then(data => setPosts(data))
            .catch(error => console.error('Error fetching posts:', error));
    }, []);

    const handleImageChange = (event) => {
        const file = event.target.files[0];
        if (file) {
            setSelectedImage(URL.createObjectURL(file));
        }
    };

    const handlePostSubmit = async () => {
        if (!newPost.trim() && !selectedImage) return;

        const formData = new FormData();
        formData.append('content', newPost);
        formData.append('image', selectedImage);

        const response = await fetch('http://localhost:8080/api/posting', {
            method: 'POST',
            body: formData
        });

        const newPostData = await response.json();
        setPosts(prevPosts => [newPostData, ...prevPosts]);
        setNewPost('');
        setSelectedImage(null);

    };


    return (
        <div className={styles.container}>
            <nav className={styles.navbar}>
                <h1>SocialNet</h1>
            </nav>
            <main className={styles.main}>
                <div className={styles.feed}>
                    <div className={styles.createPost}>
                        <Image src="http://localhost:8080/api/avatar/1" alt="User" width={32} height={32} />
                        <input
                            type="text"
                            placeholder="What's on your mind?"
                            value={newPost}
                            onChange={(e) => setNewPost(e.target.value)}
                        />
                        <input
                            type="file"
                            accept="image/*"
                            onChange={handleImageChange}
                        />
                        {selectedImage && <Image src={selectedImage} alt="Selected" width={100} height={100} />}
                        <button onClick={handlePostSubmit}>Post</button>
                    </div>
                    {posts.map((post) => (
                        <div key={post.id} className={styles.post}>
                            <div className={styles.postHeader}>
                                <Image src={post.user.avatar} alt={post.user.name} width={32} height={32} />
                                <div>
                                    <strong>{post.user.name}</strong>
                                    <p>{post.user.username} Â· {post.timestamp}</p>
                                </div>
                            </div>
                            <p>{post.content}</p>
                            {post.image && <Image src={post.image} alt="Post Image" width={300} height={300} />}
                            <div className={styles.postActions}>
                                <button>Like ({post.likes})</button>
                                <button>Dislike ({post.Dislike})</button>
                                <button>Comment ({post.comments})</button>
                            </div>
                        </div>
                    ))}
                </div>
            </main>
        </div>
    );
};

export default HomePage;