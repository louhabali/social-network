"use client";
import React, { useState, useEffect } from "react";
import Image from "next/image";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileArrowUp, faShare } from "@fortawesome/free-solid-svg-icons";
import styles from "../app/page.module.css";
import "./postform.css";

const PostForm = ({ userInfo, setPosts, setOffset, groupId }) => {
    const [newPost, setNewPost] = useState("");
    const [selectedImage, setSelectedImage] = useState(null);
    const [sendImage, setSendImage] = useState(null);

    const [ops, setOp] = useState("public")
    const [privateUsers, setPrivateUsers] = useState([])
    const [selectedOptions, setSelectedOptions] = useState([]);
    // select 
    const handleToggle = (option) => {
        setSelectedOptions((prev) =>
            prev.includes(option) ? prev.filter((o) => o !== option) : [...prev, option]
        );
    };
    useEffect(() => {
        if (ops === "private") {
            if (document.getElementById("popup")) document.getElementById("popup").style.display = "flex"

            fetch('http://localhost:8080/api/getfollowedusers', {
                credentials: "include",
            }).then(data => data.json()).then(d => {
<<<<<<< HEAD

                setPrivateUsers(d)

=======
                setPrivateUsers(d)
>>>>>>> alouhab
            })




        } else {
            if (document.getElementById("popup")) document.getElementById("popup").style.display = "none"
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
        formData.append("privacy", ops)
        if (ops == "private") {
            selectedOptions.forEach(userId => {
                formData.append("selectedUsers", userId);
            });
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
        <div className="post-form-container">
            <div className="create-post">
                <div className="user-avatar">
                    {userInfo && userInfo.Avatar ? (
                        <Image src={`data:image/png;base64,${userInfo.Avatar}`} alt="User Avatar" width={40} height={40} />
                    ) : (
                        <Image src="/image/no-avatar.png" alt="No Avatar" width={40} height={40} />
                    )}
                </div>

                <div className="post-input-container">
                    <div className="input-with-post-btn">
                        <input
                            type="text"
                            className="post-input"
                            placeholder="What's on your mind?"
                            value={newPost}
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

            {privateUsers && (
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
                    <Image
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
