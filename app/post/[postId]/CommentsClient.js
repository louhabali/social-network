'use client';

import { useState, useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFileArrowUp, faShare } from '@fortawesome/free-solid-svg-icons';
import Image from "next/image";

export default function ClientSideComments({ postId,lenComments,initialComments  }) {
    const [comments, setComments] = useState(initialComments ||[]);
    const [content, setContent] = useState('');
    const [offset, setOffset] = useState(initialComments?.length);
    const [lenC, setlenC] = useState(lenComments);
    const [hasMore, setHasMore] = useState(lenComments > initialComments?.length);
    const [selectedImage, setSelectedImage] = useState(null);
    const [sendImage, setSendImage] = useState(null);

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

    async function loadMoreComments() {
        const newOffset = offset + 4;
        const response = await fetch(`http://localhost:8080/api/posts/${postId}/comments?offset=${offset}`, {
            credentials: "include",
        });

        if (!response.ok) return;

        const newComments = await response.json();
        
        if (newComments) {
            setComments(prev => [...prev, ...newComments]);
            setOffset(newOffset);
            setHasMore(lenC > newOffset)
            console.log(lenC,newOffset);
            
        }else{
            setHasMore(false)
        }
    }

    const handleCommentSubmit = async () => {
        if (!content.trim() && !selectedImage) return;
        const formData = new FormData();
        formData.append('content', content);
        formData.append('image', sendImage);
        const response = await fetch(`http://localhost:8080/api/posts/${postId}/comments`, {
            method: 'POST',
            body: formData,
            credentials: "include",
        });

        if (!response.ok) return;

        const newComment = await response.json();
        setComments(prev => prev ? [newComment, ...prev] : [newComment]);
        setOffset(prevOffset => prevOffset + 1);
        setContent('');
        setSelectedImage(null);
        setSendImage(null);
        document.getElementById('fileInput').value = ""
        setlenC(lenComments+1)
    }

    return (
        <div>
            <input
                type="text"
                placeholder="Write a comment..."
                value={content}
                onChange={(e) => setContent(e.target.value)}
            />
            <div style={{ cursor: 'pointer' }} onClick={triggerFileInput}>
                <FontAwesomeIcon icon={faFileArrowUp} />
            </div>
            <input
                id="fileInput"
                type="file"
                accept="image/*"
                onChange={handleImageChange}
                style={{ display: 'none' }}
            />
            {selectedImage && <Image src={selectedImage} alt="Selected" width={100} height={100} />}
            <button onClick={handleCommentSubmit}>Post</button>

            <div className="comments-list">
                {comments.map(comment => (
                    <div key={comment.Id} className="comment">
                        {comment.Avatar ? (
                            <Image src={`data:image/png;base64,${comment.Avatar}`} alt="" width={32} height={32} />
                        ) : (
                            <Image src="/image/no-avatar.png" alt="" width={32} height={32} />
                        )}
                        <strong>{comment.Username}</strong> <span>{comment.CreatedAt}</span>
                        <div>{comment.Content}</div>
                        {comment.Image && <Image src={"http://localhost:8080/api/images/"+comment.Image} alt="Comment Image" width={300} height={300} />}
                    </div>
                ))}
            </div>

            {hasMore && (
                <button onClick={loadMoreComments}>
                    Load More Comments
                </button>
            )}
        </div>
    );
}
