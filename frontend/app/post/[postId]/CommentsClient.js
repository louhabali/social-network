'use client';

import { useState, useContext } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFileArrowUp } from '@fortawesome/free-solid-svg-icons';
import { WebSocketContext } from "@/context/WebSocketContext";
import Image from "next/image";
import './comments.css'

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";


export default function ClientSideComments({ postId, lenComments, initialComments }) {
    const [comments, setComments] = useState(initialComments || []);
    const [content, setContent] = useState('');
    const [offset, setOffset] = useState(initialComments?.length);
    const [lenC, setlenC] = useState(lenComments);
    const [hasMore, setHasMore] = useState(lenComments > initialComments?.length);
    const [selectedImage, setSelectedImage] = useState(null);
    const [sendImage, setSendImage] = useState(null);
    const { showNotification } = useContext(WebSocketContext);


    const handleImageChange = (event) => {
        const file = event.target.files[0];
        if (file) {
            if (file.size > 4 * 1024 * 1024) {
                showNotification("Profile image must be less than 4MB", "error");
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
        document.getElementById('fileInput').click()
    }

    async function loadMoreComments() {
        const newOffset = offset + 4;
        const response = await fetch(`${apiUrl}/api/posts/${postId}/comments?offset=${offset}`, {
            credentials: "include",
        });

        if (!response.ok) {
            showNotification("Error loading more comments", "error");
            return;
        }

        const newComments = await response.json();

        if (newComments) {
            setComments(prev => [...prev, ...newComments]);
            setOffset(newOffset);
            setHasMore(lenC > newOffset)
        } else {
            setHasMore(false)
        }
    }

    const handleCommentSubmit = async () => {
        if (!content.trim() && !selectedImage) return;
        if (content.length > 200) {
            showNotification("Comment must be less than 200 characters", "error");
            return;
        }

        const formData = new FormData();
        formData.append('content', content);
        formData.append('image', sendImage);
        const response = await fetch(`${apiUrl}/api/posts/${postId}/comments`, {
            method: 'POST',
            body: formData,
            credentials: "include",
        });

        if (!response.ok) {
            showNotification("Error posting comment", "error");
            return;
        }

        const newComment = await response.json();
        setComments(prev => prev ? [newComment, ...prev] : [newComment]);
        setOffset(prevOffset => prevOffset + 1);
        setContent('');
        setSelectedImage(null);
        setSendImage(null);
        document.getElementById('fileInput').value = ""
        setlenC(lenComments + 1)
    }

    return (
        <div className="comments-container">
            <div className="comment-input-wrapper">
                <div className="comment-input-area">
                    <input
                        type="text"
                        placeholder="Write a comment..."
                        value={content}
                        maxLength={200}
                        onChange={(e) => setContent(e.target.value)}
                        className="comment-input"
                    />
                    <div className="comment-actions">
                        <div className="comment-attachment" style={{ cursor: 'pointer' }} onClick={triggerFileInput}>
                            <FontAwesomeIcon icon={faFileArrowUp} />
                        </div>
                        <button onClick={handleCommentSubmit} className="comment-submit-btn">
                            Post
                        </button>
                    </div>
                    <input
                        id="fileInput"
                        type="file"
                        accept="image/*"
                        onChange={handleImageChange}
                        style={{ display: 'none' }}
                    />
                </div>
                {selectedImage && (
                    <div className="selected-image-preview">
                        <img src={selectedImage} alt="Selected" width={100} height={100} />
                    </div>
                )}
            </div>

            <div className="comments-list">
                {comments.map(comment => (
                    <div key={comment.Id} className="comment-item">
                        <div className="comment-avatar">
                            {comment.Avatar ? (
                                <img src={`${apiUrl}/api/images/${comment.Avatar}`} alt="" width={32} height={32} />
                            ) : (
                                <img src="/image/no-avatar.png" alt="" width={32} height={32} />
                            )}
                        </div>
                        <div className="comment-content">
                            <div className="comment-header">
                                <span className="comment-username">{comment.Username}</span>
                                <span className="comment-time">{comment.CreatedAt}</span>
                            </div>
                            <div className="comment-text">{comment.Content}</div>
                            {comment.Image && (
                                <div className="comment-image">
                                    <img
                                        src={`${apiUrl}/api/images/` + comment.Image}
                                        alt="Comment Image"
                                        width={300}
                                        height={300}
                                    />
                                </div>
                            )}
                        </div>
                    </div>
                ))}
            </div>

            {hasMore && (
                <div className="load-more-container">
                    <button onClick={loadMoreComments} className="load-more-btn">
                        Load More Comments
                    </button>
                </div>
            )}
        </div>
    );
}
