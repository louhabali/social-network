"use client";
import React, { useState, useEffect, useContext } from "react";
import Image from "next/image";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileArrowUp, faShare } from "@fortawesome/free-solid-svg-icons";
import { WebSocketContext } from "@/context/WebSocketContext";
import "./postform.css";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

const PostForm = ({ userInfo, setPosts, setOffset, groupId }) => {
    const [newPost, setNewPost] = useState("");
    const [selectedImage, setSelectedImage] = useState(null);
    const [sendImage, setSendImage] = useState(null);
    const [showprivate,setShowprivate]= useState(false)
    const [ops, setOp] = useState("public")
    const [privateUsers, setPrivateUsers] = useState([])
    const [selectedOptions, setSelectedOptions] = useState([]);

    const { showNotification } = useContext(WebSocketContext);

    // select 
    const handleToggle = (option) => {
        setSelectedOptions((prev) =>
            prev.includes(option) ? prev.filter((o) => o !== option) : [...prev, option]
        );
    };
    useEffect(() => {
        if (ops === "private") {
            setShowprivate(true)

            fetch(`${apiUrl}/api/getfollowedusers`, {
                credentials: "include",
            })
            .then(res => {
                if (!res.ok) {
                    showNotification("Error fetching followed users", "error");
                    return null;
                }
                return res.json();
            })
            .then(data => {
                setPrivateUsers(data);
            })
            .catch(() => {
                showNotification("Network error while fetching followed users", "error");
            });

        } else {
           setShowprivate(false)
        }
    }, [ops]);



    //
    const handleImageChange = (event) => {
        const file = event.target.files[0];
        if (file) {
            if (file.size > 4*1024 * 1024) {
                showNotification("Image size must be under 4MB", "error");
                return;
              }
              const allowedTypes = ['image/jpeg', 'image/png', 'image/gif'];
              if (!allowedTypes.includes(file.type)) {
                showNotification("Only JPG, PNG, or GIF images are allowed", "error");
                return;
              }
            setSelectedImage(URL.createObjectURL(file));
            setSendImage(file);
        }
    };

    const triggerFileInput = () => {
        document.getElementById("fileInput").click();
    };

    const handlePostSubmit = async () => {
        if (!newPost.trim() && !selectedImage) {
            showNotification("Post cannot be empty!", "error");
            return;
        }
        if (newPost.length > 200) {
            showNotification("Post content must be less than 200 characters.", "error");
            return;
        }
        if (ops === "private" && selectedOptions.length === 0) {
            showNotification("Please select at least one user for a private post.", "error");
            return;
        }
        if ((ops === "almost_private" || ops === "public") && selectedOptions.length > 0) {
            showNotification("You can't select users when the post is almost private or public.", "error");
            return;
        }
        if (ops !== "public" && ops !== "private" && ops !== "almost_private") {
            showNotification("Invalid privacy option selected.", "error");
            return;
        }
        
        const formData = new FormData();
        formData.append("content", newPost);
        formData.append("groupId", groupId);
        formData.append("privacy", ops)
        if (ops == "private") {
            selectedOptions.forEach(userId => {
                formData.append("selectedUsers", userId);
            });
        }

        if (sendImage) formData.append("image", sendImage);

        try {
            const response = await fetch(`${apiUrl}/api/posting`, {
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
            setOp("public")
            setPrivateUsers([])
            setSelectedOptions([])
            showNotification("Post created successfully!", "info");

        } catch (error) {
            showNotification("Error creating post. Try again!", "error");
        }
    };

    return (
        <div className="post-form-container">
            <div className="create-post">
                <div className="user-avatar">
                    {userInfo && userInfo.Avatar ? (
                        <img src={`${apiUrl}/api/images/${userInfo.Avatar}`} className="userAvatar" alt="User Avatar" width={40} height={40} />
                    ) : (
                        <img src="/image/no-avatar.png" alt="No Avatar" width={40} height={40} className="userAvatar" />
                    )}
                </div>

                <div className="post-input-container">
                    <div className="input-with-post-btn">
                        <input
                            type="text"
                            className="post-input"
                            placeholder="What's on your mind?"
                            value={newPost}
                            maxLength={200}
                            onChange={(e) => setNewPost(e.target.value)}
                        />
                        <div className="post-icon-btn" onClick={handlePostSubmit}>
                            <FontAwesomeIcon icon={faShare} />
                        </div>
                    </div>

                    <div className="post-actions">
                        <div className="post-tools">
                            <div className="upload-button" onClick={triggerFileInput}>
                                <FontAwesomeIcon icon={faFileArrowUp} />
                                <span>Photo</span>
                            </div>

                            <div className="privacy-selector">
                                <select
                                    id="slct"
                                    className="privacy-dropdown"
                                    value={ops}
                                    onChange={(o) => setOp(o.target.value)}
                                >
                                    <option value="public">Public</option>
                                    <option value="private">Only for</option>
                                    <option value="almost_private">Only followers</option>
                                </select>
                            </div>
                        </div>

                        <input
                            id="fileInput"
                            type="file"
                            accept="image/*"
                            onChange={handleImageChange}
                            style={{ display: "none" }}
                        />
                    </div>
                </div>
            </div>

            {privateUsers && showprivate && (
                <div id='popup' className="private-users-popup">
                    {privateUsers.map((p) => (
                        <label key={p.id} className="user-checkbox">
                            <input
                                type="checkbox"
                                checked={selectedOptions.includes(p.id)}
                                onChange={() => handleToggle(p.id)}
                            />
                            <span>{`${p.fname} ${p.lname}`}</span>
                        </label>
                    ))}
                </div>
            )}

            {selectedImage && (
                <div className="selected-image-container">
                    <img
                        src={selectedImage}
                        alt="Selected"
                        width={400}
                        height={200}
                    />
                </div>
            )}
        </div>);
};

export default PostForm;
