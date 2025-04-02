import React, { useState, useContext } from "react";
import Image from "next/image";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFileArrowUp, faShare } from "@fortawesome/free-solid-svg-icons";
import { WebSocketContext } from "@/context/WebSocketContext";
import "./postform.css";

const GroupsPostForm = ({ userInfo, setPosts, setOffset, groupId }) => {
    const [newPost, setNewPost] = useState("");
    const [selectedImage, setSelectedImage] = useState(null);
    const [sendImage, setSendImage] = useState(null);
    const { showNotification } = useContext(WebSocketContext);
    

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
        const formData = new FormData();
        formData.append("content", newPost);
        formData.append("groupId", groupId);
        if (newPost.length > 200 || newPost.length < 3) {
            showNotification("Post must be between 3 and 200 characters.", "error");
            return;
        }

        if (sendImage) formData.append("image", sendImage);
        formData.append("privacy", "public");
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
            showNotification("Post created successfully!", "info");
        } catch (error) {
            showNotification("Failed to post", "error");
        }
    };

    return (
        <div className="post-form-container">
            <div className="create-post">
                <div className="user-avatar">
                    {userInfo && userInfo.Avatar ? (
                        <Image src={`http://localhost:8080/api/images/${userInfo.Avatar}`} alt="User Avatar" width={40} height={40} />
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
                        </div>

                        <input
                            id="fileInput"
                            type="file"
                            accept="image/*"
                            onChange={handleImageChange}
                            style={{ display: "none" }}
                        />
                        {selectedImage && (
                            <button type="button" onClick={() => { setSelectedImage(null);setSendImage(null);document.getElementById("fileInput").value = ""; }}>
                                Cancel
                            </button>
                        )}
                    </div>
                </div>
            </div>

            {selectedImage && (
                <div className="selected-image-container">
                    <Image src={selectedImage} alt="Selected" width={400} height={200} />
                </div>
            )}
        </div>
    );
};

export default GroupsPostForm;
