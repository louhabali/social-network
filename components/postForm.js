"use client";
import React, { useState,useEffect } from "react";
import Image from "next/image";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileArrowUp, faShare } from "@fortawesome/free-solid-svg-icons";
import styles from "../app/page.module.css";

const PostForm = ({ userInfo, setPosts,setOffset,groupId }) => {
    const [newPost, setNewPost] = useState("");
    const [selectedImage, setSelectedImage] = useState(null);
    const [sendImage, setSendImage] = useState(null);

    const [ops, setOp] = useState("public")
    const [privateUsers,setPrivateUsers] = useState([])
    const [selectedOptions, setSelectedOptions] = useState([]);
 // select 
    const handleToggle = (option) => {
        setSelectedOptions((prev) =>
        prev.includes(option) ? prev.filter((o) => o !== option) : [...prev, option]
         );
    };
    useEffect(() => {
        if (ops === "private") {
            if (document.getElementById("popup")) document.getElementById("popup").style.display="flex"
          
            fetch('http://localhost:8080/api/getfollowedusers', {
                credentials: "include",
            }).then(data => data.json()).then(d => {
                
                setPrivateUsers(d)
              
            })
           
            
            
            
        }else {
            if (document.getElementById("popup")) document.getElementById("popup").style.display="none"
        }
    }, [ops]);
    
    
    
//
    const handleImageChange = (event) => {
        const file = event.target.files[0];
        if (file) {
            setSelectedImage(URL.createObjectURL(file));
            setSendImage(file);
        }
    };

    const triggerFileInput = () => {
        document.getElementById("fileInput").click();
    };

    const handlePostSubmit = async () => {
        if (!newPost.trim() && !selectedImage) return;

        const formData = new FormData();
        formData.append("content", newPost);
        formData.append("groupId", groupId);
        formData.append("privacy",ops)
        if (ops=="private"){
            formData.append("selectedUsers",selectedOptions)
        }
        
        if (sendImage) formData.append("image", sendImage);

        try {
            const response = await fetch("http://localhost:8080/api/posting", {
                method: "POST",
                body: formData,
                credentials: "include",
            });

            if (!response.ok) throw new Error("Failed to post");

            const newPostData = await response.json();
            setPosts((prevPosts) => (prevPosts ? [newPostData, ...prevPosts] : [newPostData]));
            setOffset((prevOffset) => prevOffset + 1);

            setNewPost("");
            setSelectedImage(null);
            setSendImage(null);
            document.getElementById("fileInput").value = "";
            //document.getElementById("popup").style.display="none"
            

        } catch (error) {
            console.error("Error posting:", error);
        }
    };

    return (
        <div className={styles.createPost}>
            {userInfo && userInfo.Avatar ? (
                <Image src={`data:image/png;base64,${userInfo.Avatar}`} alt="User Avatar" width={32} height={32} />
            ) : (
                <Image src="/image/no-avatar.png" alt="No Avatar" width={32} height={32} />
            )}

            <input type="text" placeholder="What's on your mind?" value={newPost} onChange={(e) => setNewPost(e.target.value)} />

            <div style={{ cursor: "pointer" }} onClick={triggerFileInput}>
                <FontAwesomeIcon icon={faFileArrowUp} />
            </div>

            <input id="fileInput" type="file" accept="image/*" onChange={handleImageChange} style={{ display: "none" }} />
            <select id="slct" onChange={(o)=> setOp(o.target.value)}>
                    <option value="public">public</option>
                    <option value="private">only for</option>
                    <option value="almost_private">only followers</option>
            </select>
          {privateUsers && <div id='popup'  style={{
                  
                    height: "45px",
                    display:"flex",
                    flexDirection:"column",
                    overflowY: "auto",
                    border: "1px solid #ccc", 
                    padding: "5px",
                    backgroundColor: "black",
                    borderRadius: "8px", 
                    boxShadow: "0 2px 10px rgba(154, 148, 148, 0.77)",
  }} >
                    {privateUsers && privateUsers.map((p) => (
                        <label key={p.id} className="flex items-center space-x-2">
                    <input
                    type="checkbox"
                    checked={selectedOptions.includes(p.id)}
                    onChange={() => handleToggle(p.id)}
                    />
                    <span>{`${p.fname} ${p.lname}`}</span>
                        </label>
                    ))}
                </div>}
            <div className={styles.btnsub} onClick={handlePostSubmit}>
                <FontAwesomeIcon icon={faShare} />
            </div>

            {selectedImage && <Image src={selectedImage} style={{ borderRadius: "5px" }} alt="Selected" width={400} height={200} />}
        </div>
    );
};

export default PostForm;
