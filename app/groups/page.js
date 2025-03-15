"use client";
import Link from "next/link"
import React, { useState, useEffect, useRef } from "react";
import styles from "./Groups.module.css";

export default function GroupsPage() {

    const [groups, setGroups] = useState([]);
    const [isOpen, setIsOpen] = useState(false);
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [offset, setOffset] = useState(groups?.length);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const observer = useRef(null);

    function openModal() {
        setIsOpen(true);
    }

    function closeModal() {
        setIsOpen(false);
    }

    const fetchGroups = async () => {
        if (loading || !hasMore) return;

        setLoading(true);
        const res = await fetch(`http://localhost:8080/api/groups?offset=${offset}`, {
            credentials: "include",
        });
        if (!res.ok) return;
        const data = await res.json();
        if (data === null) {
            setHasMore(false);
            return
        }
        if (data.length < 4) {
            setHasMore(false);
        };
        setGroups((prevGroups) => [...prevGroups, ...data]);
        setOffset((prevOffset) => prevOffset + data?.length);
        setLoading(false);
    }

    async function handleSubmit(e) {
        e.preventDefault();
        const response = await fetch("http://localhost:8080/api/groups", {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({ title, description }),
        });

        const newGroup = await response.json();
        console.log(newGroup);
        setGroups((prevGroups) => (prevGroups ? [newGroup, ...prevGroups] : [newGroup]));
        setTitle("");
        setDescription("");
        setOffset((prevOffset) => prevOffset + 1);
        setIsOpen(false);
    }


    useEffect(() => {
        if (loading || !hasMore) return;

        observer.current = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting) {
                    fetchGroups();
                }
            },
            { threshold: 1.0 }
        );

        if (observer.current && document.getElementById("load-more")) {
            observer.current.observe(document.getElementById("load-more"));
        }

        return () => observer.current?.disconnect();
    }, [fetchGroups, loading, hasMore]);


    async function handleJoin(groupId, currentStatus) {
        let newStatus = currentStatus === "pending" ? "" : "pending";
        let method = currentStatus === "pending" ? "DELETE" : "POST";

        const response = await fetch(`http://localhost:8080/api/groups/${groupId}/join`, {
            method,
            credentials: "include",
        });

        if (response.ok) {
            setGroups((prevGroups) =>
                prevGroups.map((group) =>
                    group.id === groupId ? { ...group, status: newStatus } : group
                )
            );
        }

    }
        
    return (
        <div className={styles.container}>
            <div className={styles.pageHeader}>
                <h1>Groups</h1>
                <button className={styles.createGroupBtn} onClick={openModal}>Create Group</button>
            </div>

            <ul className={styles.groupsList}>
                {groups?.map((group) => (
                    <li key={group.id} className={styles.groupItem}>
                        <h3 className={styles.groupTitle}>{group.title}</h3>
                        <p className={styles.groupDescription}>{group.description}</p>
                        
                        <div className={styles.groupActions}>
                            {group.status === "approved" || group.status === "creator" ? (
                                <Link href={`/groups/${group.id}`}>
                                    <button className={styles.openBtn}>Open</button>
                                </Link>
                            ) : group.status === "pending" ? (
                                <div className={styles.pendingStatus}>
                                    <span>Request Pending</span>
                                    <button 
                                        className={styles.cancelBtn}
                                        onClick={() => handleJoin(group.id, "pending")}
                                    >
                                        Cancel Request
                                    </button>
                                </div>
                            ) : (
                                <button 
                                    className={styles.joinBtn} 
                                    onClick={() => handleJoin(group.id)}
                                >
                                    Join Group
                                </button>
                            )}
                        </div>
                    </li>
                ))}
            </ul>

            {hasMore && <div id="load-more" className={styles.loadMore} />}
            {loading && <p className={styles.loadingIndicator}>Loading...</p>}

            {isOpen && (
                <div className={styles.modalOverlay} onClick={closeModal}>
                    <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
                        <h2>Create a New Group</h2>
                        <form onSubmit={handleSubmit}>
                            <div className={styles.formGroup}>
                                <label htmlFor="groupTitle">Group Title</label>
                                <input
                                    id="groupTitle"
                                    type="text"
                                    placeholder="Enter a name for your group"
                                    value={title}
                                    onChange={(e) => setTitle(e.target.value)}
                                    required
                                />
                            </div>
                            
                            <div className={styles.formGroup}>
                                <label htmlFor="groupDescription">Description</label>
                                <textarea
                                    id="groupDescription"
                                    placeholder="What's your group about?"
                                    value={description}
                                    onChange={(e) => setDescription(e.target.value)}
                                    required
                                ></textarea>
                            </div>
                            
                            <div className={styles.modalButtons}>
                                <button type="button" onClick={closeModal}>Cancel</button>
                                <button type="submit">Create</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}