'use client';

import { useState, useRef, useEffect } from 'react';
import Post from "@/components/post";
import PostForm from "@/components/postForm";
import InviteUsers from "./InviteUsers";

export default function ClientSideGroups({ groupId, status, initialPosts, userInfo }) {
    const [posts, setPosts] = useState(initialPosts || []);
    const [events, setEvents] = useState([]);
    const [offset, setOffset] = useState(0);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const [activeTab, setActiveTab] = useState("posts");
    const observer = useRef(null);
    const [showEventForm, setShowEventForm] = useState(false);
    const [newEvent, setNewEvent] = useState({
        title: "",
        description: "",
        date: "",
    });
    const [memberStatus,setMemberStatus] = useState(status)

    const fetchPosts = async () => {
        if (loading || !hasMore) return;

        setLoading(true);
        try {
            const response = await fetch(`http://localhost:8080/api/posts?offset=${offset}&groupId=${groupId}`, {
                credentials: "include",
            });
            if (!response.ok) {
                router.push("/login");
                return;
            }
            const data = await response.json();
            if (data === null) {
                setHasMore(false);
                setLoading(false)
                return
            }
            if (data.length < 4) {
                setHasMore(false);
            };
            setPosts((prevPosts) => [...prevPosts, ...data]);
            setOffset((prevOffset) => prevOffset + data?.length);
        } catch (error) {
            console.error("Error fetching posts:", error);
        } finally {
            setLoading(false);
        }
    };


    useEffect(() => {
        if (loading || !hasMore) return;

        observer.current = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting) {
                    fetchPosts();
                }
            },
            { threshold: 1.0 }
        );

        if (observer.current && document.getElementById("load-more")) {
            observer.current.observe(document.getElementById("load-more"));
        }

        return () => observer.current?.disconnect();
    }, [fetchPosts, loading, hasMore]);


    const handleLike = async (e) => {
        const postId = e.target.closest("section").id;
        const likeIcon = e.target;

        const postIndex = posts.findIndex((post) => post.Pid.toString() === postId);
        const updatedPosts = [...posts];

        let likeStatus = updatedPosts[postIndex].UserInteraction === 1 ? 0 : 1;

        const response = await fetch("http://localhost:8080/api/like", {
            method: "POST",
            body: JSON.stringify({ postId, likeStatus }),
            credentials: "include",
        });

        if (response.ok) {
            updatedPosts[postIndex].UserInteraction = likeStatus
            updatedPosts[postIndex].Likes += likeStatus ? 1 : -1
            likeIcon.style.color = likeStatus ? 'red' : 'white'
            setPosts(updatedPosts);
        }
    };



    const fetchEvents = async () => {
        try {
            const response = await fetch(`http://localhost:8080/api/groups/${groupId}/events`, {
                credentials: "include",
            });
            if (response.ok) {
                const data = await response.json();
                setEvents(data);
            }
        } catch (error) {
            console.error("Error fetching events:", error);
        }
    };

    useEffect(() => {
        if (activeTab === "events") {
            fetchEvents();
        }
    }, [activeTab]);


    const handleCreateEvent = async () => {
        try {
            const response = await fetch(`http://localhost:8080/api/groups/${groupId}/events`, {
                method: "POST",
                body: JSON.stringify(newEvent),
                credentials: "include",
            });

            if (response.ok) {
                const createdEvent = await response.json();
                setEvents(prevEvents => prevEvents ? [createdEvent, ...prevEvents] : []);
                setShowEventForm(false);
                setNewEvent({ title: "", description: "", date: "" });
            } else {
                console.error("Failed to create event");
            }
        } catch (error) {
            console.error("Error creating event:", error);
        }
    };


    const handleGoing = async (eventId, status) => {
        try {
            const response = await fetch(`http://localhost:8080/api/groups/${groupId}/events/going`, {
                method: "POST",
                body: JSON.stringify({ eventId, status }),
                credentials: "include",
            });

            if (response.ok) {
                setEvents((prevEvents) =>
                    prevEvents.map((event) =>
                        event.id === eventId ? { ...event, userStatus: status } : event
                    )
                );
            } else {
                console.error("Failed to update RSVP");
            }
        } catch (error) {
            console.error("Error updating RSVP:", error);
        }
    };

    async function handleJoinGroup() {
        let method = memberStatus === "pending" ? "DELETE" : "POST";
    
        const response = await fetch(`http://localhost:8080/api/groups/${groupId}/join`, {
            method,
            credentials: "include",
        });
    
        if (response.ok) {
            setMemberStatus(memberStatus === "pending" ? "" : "pending");
        }
    }
    
    if (memberStatus != "approved" && memberStatus != "creator") {
        return (
            <div>
                <p>This is a private group. Join to see posts and events.</p>
                {memberStatus === "pending" ? (
                    <>
                    <span>Request Pending </span>
                    <button onClick={handleJoinGroup}>Cancel Join Request</button>
                    </>
                ) : (
                    <button onClick={handleJoinGroup}>Join Group</button>
                )}
            </div>
        );
    }
    

    return (
        <div>
            <InviteUsers groupId={groupId} />
            <div className="group-tabs">
                <button className={activeTab === "posts" ? "active" : ""} onClick={() => setActiveTab("posts")}>
                    Posts
                </button>
                <button className={activeTab === "events" ? "active" : ""} onClick={() => setActiveTab("events")}>
                    events
                </button>
            </div>

            {activeTab === "posts" && (
                <>
                    <PostForm userInfo={userInfo} setPosts={setPosts} setOffset={setOffset} groupId={groupId} />

                    {posts.map((post) => (
                        <Post key={post.Pid} post={post} handleLike={handleLike} />
                    ))}

                    {hasMore && <div id="load-more" style={{ height: "10px" }} />}
                    {loading && <p>Loading...</p>}
                </>
            )}

            {activeTab === "events" && (
                <>
                    <button onClick={() => setShowEventForm(true)}>+ Add Event</button>

                    {events === null || events.length === 0 ? (
                        <p>No events available.</p>
                    ) : (
                        events.map((event) => (
                            <div key={event.id} className="event-card">
                                <h3>{event.title}</h3>
                                <p>{event.description}</p>
                                <p>Date: {new Date(event.date).toLocaleDateString()}</p>
                                <div>
                                    <button
                                        className={event.userStatus === "going" ? "active" : ""}
                                        onClick={() => handleGoing(event.id, "going")}>
                                        Going
                                    </button>
                                    <button
                                        className={event.userStatus === "not_going" ? "active" : ""}
                                        onClick={() => handleGoing(event.id, "not_going")}>
                                        Not Going
                                    </button>
                                </div>
                            </div>
                        ))
                    )}
                </>
            )}

            {showEventForm && (
                <div className="popup">
                    <div className="popup-content">
                        <h2>Create Event</h2>
                        <input
                            type="text"
                            placeholder="Title"
                            value={newEvent.title}
                            onChange={(e) => setNewEvent({ ...newEvent, title: e.target.value })}
                        />
                        <textarea
                            placeholder="Description"
                            value={newEvent.description}
                            onChange={(e) => setNewEvent({ ...newEvent, description: e.target.value })}
                        />
                        <input
                            type="datetime-local"
                            value={newEvent.date}
                            onChange={(e) => setNewEvent({ ...newEvent, date: e.target.value })}
                        />
                        <button onClick={handleCreateEvent}>Create</button>
                        <button onClick={() => setShowEventForm(false)}>Cancel</button>
                    </div>
                </div>
            )}
        </div>
    );
}