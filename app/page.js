"use client"
import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import styles from './page.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFileArrowUp, faShare } from '@fortawesome/free-solid-svg-icons';
import { useRouter } from "next/navigation"

const HomePage = () => {
    const [posts, setPosts] = useState([]);
    const [newPost, setNewPost] = useState('');
    const [selectedImage, setSelectedImage] = useState(null);
    const [sendImage, setSendImage] = useState(null);
    const [userInfo, setUserInfo] = useState(null);


  const router = useRouter()

    useEffect(() => {
        fetch('http://localhost:8080/api/posts', {
            credentials: "include",
        })
            .then(response => {
                if (!response.ok) {
                    router.push("/login")  
                }    
                return response.json()
            })
            .then(data => setPosts(data))
            .catch(error => console.error('Error fetching posts:', error));
    }, []);
    console.log("psts: ",posts);
    useEffect(() => {
        fetch('http://localhost:8080/api/userdata', {
            credentials: "include",
        })
            .then(response => response.json())
            .then(data => setUserInfo(data))
            .catch(error => console.error('Error fetching posts:', error));
    }, []);


    const handleImageChange = (event) => {
        const file = event.target.files[0];        
        if (file) {
            setSelectedImage(URL.createObjectURL(file));
            setSendImage(file);
        }
    };
    const triggerFileInput = () => {
        document.getElementById('fileInput').click()
    }
    const handlePostSubmit = async () => {
        if (!newPost.trim() && !selectedImage) return;

        const formData = new FormData();
        formData.append('content', newPost);
        formData.append('image', sendImage);
        
        const response = await fetch('http://localhost:8080/api/posting', {
            method: 'POST',
            body: formData,
            credentials: "include",

        });
        
        const newPostData = await response.json();
        
        setPosts(prevPosts => prevPosts ? [newPostData, ...prevPosts]:[newPostData]);
        setNewPost('');
        setSelectedImage(null);
        setSendImage(null);
        document.getElementById('fileInput').value = ""

    };

  
    

    return (
        <div className={styles.container}>
            <nav className={styles.navbar}>
            </nav>
            <main className={styles.main}>
                <div className={styles.feed}>
                    <div className={styles.createPost}>

                        {userInfo && userInfo.Avatar ? (
                            <Image src={`data:image/png;base64,${userInfo.Avatar}`} alt="" width={32} height={32} />
                        ) : (
                            <Image src="/image/no-avatar.png" alt="" width={32} height={32} />
                        )}

                        <input
                            type="text"
                            placeholder="What's on your mind?"
                            value={newPost}
                            onChange={(e) => setNewPost(e.target.value)}
                        />
                        <div style={{ cursor: 'pointer' }} onClick={triggerFileInput}><FontAwesomeIcon icon={faFileArrowUp} /></div>
                        <input
                            id="fileInput"
                            className={styles.btnsub}
                            type="file"
                            accept="image/*"
                            onChange={handleImageChange}
                            style={{ display: 'none' }}
                        />
                        {selectedImage && <Image src={selectedImage} alt="Selected" width={100} height={100} />}
                        <div className={styles.btnsub} onClick={handlePostSubmit}><FontAwesomeIcon icon={faShare} /></div>
                    </div>
                </div>
                {posts && posts.map((post) => (
                    <div key={post.Pid} className={styles.post}>
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
                        <p>{post.content}</p>
                        {post.Image && <Image src={"http://localhost:8080/api/images/"+post.Image} alt="" width={300} height={300} />}
                        <div className={styles.postActions}>
                            <button>Like {post.Likes}</button>
                            <button>Dislike {post.Dislikes}</button>
                            <button>Comment {post.NbComment}</button>
                        </div>
                    </div>
                ))}
            </main>
        </div>
    );
};

export default HomePage;