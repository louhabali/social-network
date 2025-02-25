"use client"
import React, { useState, useEffect } from 'react';
import Image from 'next/image';
import styles from './page.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFileArrowUp, faShare,faHeart ,faComment} from '@fortawesome/free-solid-svg-icons';
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
    
    useEffect(() => {
        fetch('http://localhost:8080/api/userdata', {
            credentials: "include",
        })
            .then(response => response.json())
            .then(data => setUserInfo(data))
            .catch(error => console.error('Error fetching posts:', error));
    }, []);
    
    const handleLike = async (e) => {
        
        const pd = e.target.closest('section').id
        const likeIcon = e.target
        
        const postIndex = posts.findIndex(post => post.Pid.toString() === pd)
        const uPosts = [...posts]

        let lk = uPosts[postIndex].UserInteraction === 1 ? 0 : 1
    
       
    
        const response = await fetch(`http://localhost:8080/api/like`, {
            method: 'POST',
            body: JSON.stringify({ pd, lk }),
            credentials: "include",
        });

        if (response.ok) {
        
            if (uPosts[postIndex].UserInteraction === 1) {
                likeIcon.style.color = 'white'
                uPosts[postIndex].Likes -= 1
                uPosts[postIndex].UserInteraction = 0
            } else{
                likeIcon.style.color = 'red'
                uPosts[postIndex].Likes += 1
                uPosts[postIndex].UserInteraction = 1
            }
            setPosts(uPosts)
        }
    }
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
                        
                        <div className={styles.btnsub} onClick={handlePostSubmit}><FontAwesomeIcon icon={faShare} /></div>
                      
                    </div>
                    {selectedImage && <Image src={selectedImage}style={{borderRadius:"5px"}} alt="Selected" width={400} height={200} />}
                </div>
                {typeof posts!=="string" && posts && posts.map((post) => (
                    <section id={post.Pid} key={post.Pid} className={styles.post}>
                        <div className={styles.postHeader}>
                            {post.Avatar ? (
                                <Image style={{borderRadius:"50%"}} src={`data:image/png;base64,${post.Avatar}`} alt="" width={32} height={32} />
                            ) : (
                                <Image style={{borderRadius:"50%"}} src="/image/no-avatar.png" alt="" width={32} height={32} />

                            )}
                            <div>
                                <strong>{post.Username}</strong>
                            </div>
                            <div><span>{post.CreatedAt}</span></div>
                        </div>
                        <div>{post.Content}</div>
                        <p>{post.content}</p>
                        {post.Image && <Image src={"http://localhost:8080/api/images/"+post.Image} alt="" width={400} height={200} />}
                        <div className={styles.postActions}>
                            <FontAwesomeIcon 
                                style={{ cursor: "pointer", color: post.UserInteraction === 1 ? 'red' : 'white' }} 
                                onClick={handleLike} 
                                icon={faHeart} 
                                
                            />{post.Likes}
                            <FontAwesomeIcon icon={faComment} />
                            {post.NbComment}
                        </div>
                    </section>
                ))}
            </main>
        </div>
    );
};

export default HomePage;