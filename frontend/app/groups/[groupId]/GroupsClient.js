'use client';

import { useState, useRef, useEffect, useContext } from 'react';
import Post from "@/components/post";
import PostForm from "@/components/groupsPostForm";
import { WebSocketContext } from "@/context/WebSocketContext";
import './GroupClient.css';


const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function ClientSideGroups({ groupId, status, initialPosts, userInfo }) {
    const [posts, setPosts] = useState(initialPosts || []);
    const [events, setEvents] = useState([]);
    const [offset, setOffset] = useState(initialPosts?.length);
    const [loading, setLoading] = useState(false);
    const [loadingEvents, setLoadingEvents] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const [activeTab, setActiveTab] = useState("posts");
    const observer = useRef(null);
    const [showEventForm, setShowEventForm] = useState(false);
    const [newEvent, setNewEvent] = useState({
        title: "",
        description: "",
        date: "",
        user_response: "",
    });
    const [memberStatus, setMemberStatus] = useState(status)
    const [eventOffset, setEventOffset] = useState(0);
    const [hasMoreEvents, setHasMoreEvents] = useState(true);
    const { showNotification } = useContext(WebSocketContext);


    const fetchPosts = async () => {
        if (loading || !hasMore) return;

        setLoading(true);
        const response = await fetch(`${apiUrl}/api/posts?offset=${offset}&groupId=${groupId}`, {
            credentials: "include",
        });
        if (!response.ok) {
            showNotification("Error fetching posts", "error");
            setLoading(false);
            return;
        }
        const data = await response.json();

        if (data === null) {
            setHasMore(false);
            setLoading(false);
            return
        }

        if (data.length < 4) {
            setHasMore(false);
        };
        setPosts((prevPosts) => [...prevPosts, ...data]);
        setOffset((prevOffset) => prevOffset + data.length);
        setLoading(false);

    };

    const fetchEvents = async () => {
        if (loadingEvents || !hasMoreEvents) return;

        setLoadingEvents(true);
        try {
            const response = await fetch(`${apiUrl}/api/groups/${groupId}/events?offset=${eventOffset}`, {
                credentials: "include",
            });
            if (!response.ok) {
                showNotification("Error fetching events", "error");
                setLoadingEvents(false);
                return;
            }
            const data = await response.json();

            if (data === null) {
                setHasMoreEvents(false);
                setLoadingEvents(false);
                return
            }

            if (data.length < 4) {
                setHasMoreEvents(false);
            };
            setEvents((prevEvents) => [...prevEvents, ...data]);
            setEventOffset((prevOffset) => prevOffset + data.length);
        } catch (error) {
            console.log("Error fetching events:", error);
        } finally {
            setLoadingEvents(false);
        }
    };

    useEffect(() => {

        if (activeTab === "posts" && (loading || !hasMore)) return;
        if (activeTab === "events" && (loadingEvents || !hasMoreEvents)) return;

        if (activeTab === "posts") {
            observer.current = new IntersectionObserver((entries) => {
                if (entries[0].isIntersecting && hasMore) {
                    fetchPosts();
                }
            }, { threshold: 1.0 });

            const loadMoreElement = document.getElementById("load-more");
            if (loadMoreElement) {
                observer.current.observe(loadMoreElement);
            }

            return () => observer.current?.disconnect();
        } else if (activeTab === "events" && hasMoreEvents) {
            observer.current = new IntersectionObserver((entries) => {
                if (entries[0].isIntersecting && hasMoreEvents) {
                    fetchEvents();
                }
            }, { threshold: 1.0 });

            const loadMoreElement = document.getElementById("load-more-events");
            if (loadMoreElement) {
                observer.current.observe(loadMoreElement);
            }

            return () => observer.current?.disconnect();
        }
    }, [activeTab, loading, loadingEvents, hasMore, hasMoreEvents, fetchPosts, fetchEvents]);

    const handleCreateEvent = async () => {

        if (!newEvent.title || !newEvent.description || !newEvent.date || !newEvent.user_response) {
            showNotification("Please fill in all fields", "error");
            return;
        };
        if (newEvent.title.length > 50) {
            showNotification("Title must be less than 50 characters.");
            return
        }
        if (newEvent.description.length > 200) {
            showNotification("Description must be less than 200 characters.");
            return
        }
        if (newEvent.user_response !== "going" && newEvent.user_response !== "not_going") {
            showNotification("Invalid user response. Must be 'going' or 'not_going'.");
            return
        }

        const eventDate = new Date(newEvent.date);
        if (isNaN(eventDate.getTime())) {
            showNotification("Invalid date format. Use YYYY-MM-DDTHH:MM.");
            return
        }
        if (eventDate < new Date()) {
            showNotification("The event date cannot be in the past.");
            return
        }

        const response = await fetch(`${apiUrl}/api/groups/${groupId}/events`, {
            method: "POST",
            body: JSON.stringify(newEvent),
            credentials: "include",
        });

        if (!response.ok) {
            showNotification("Error creating event", "error");
            return;
        }
        const createdEvent = await response.json();
        setEvents(prevEvents => prevEvents ? [createdEvent, ...prevEvents] : [createdEvent]);
        setEventOffset(prevOffset => prevOffset + 1);
        setShowEventForm(false);
        setNewEvent({ title: "", description: "", date: "", user_response: "" });

        showNotification("Event created successfully", "info");
    };

    const handleGoing = async (eventId, response) => {

        const res = await fetch(`${apiUrl}/api/groups/${groupId}/events/going`, {
            method: "POST",
            body: JSON.stringify({ eventId, response }),
            credentials: "include",
        });

        if (!res.ok) {
            showNotification("Error updating event response", "error");
            return;
        }

        setEvents((prevEvents) => {
            return prevEvents.map((event) => {
                if (event.id === eventId) {

                    let user_response = event.user_response === response ? null : response;

                    let going_count = event.going_count;
                    let not_going_count = event.not_going_count;

                    if (event.user_response === "going") {
                        going_count -= 1;
                    } else if (event.user_response === "not_going") {
                        not_going_count -= 1;
                    }

                    if (user_response === "going") {
                        going_count += 1;
                    } else if (user_response === "not_going") {
                        not_going_count += 1;
                    }

                    return {
                        ...event,
                        user_response,
                        going_count,
                        not_going_count,
                    };
                }
                return event;
            });
        });

    };

    async function handleJoinGroup() {
        let method = memberStatus === "pending" ? "DELETE" : "POST";

        const response = await fetch(`${apiUrl}/api/groups/${groupId}/join`, {
            method,
            credentials: "include",
        });

        if (!response.ok) {
            showNotification("Error join group", "error");
            return;
        }
        setMemberStatus(memberStatus === "pending" ? "" : "pending");

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
        <div className="group-container">

            <div className="group-tabs">
                <button className={activeTab === "posts" ? "active" : ""} onClick={() => setActiveTab("posts")}>
                    Posts
                </button>
                <button className={activeTab === "events" ? "active" : ""} onClick={() => setActiveTab("events")}>
                    Events
                </button>
            </div>

            {activeTab === "posts" && (
                <div className="posts-container">
                    <PostForm userInfo={userInfo} setPosts={setPosts} setOffset={setOffset} groupId={groupId} />



                    <div className="posts-responsive-container">
                        {posts.length > 0 ? posts.map((post) => (
                            <Post key={post.Pid} post={post} />
                        )) : (
                            <p>No posts yet.</p>
                        )}
                    </div>
                    {hasMore && <div id="load-more" style={{ height: "10px" }} />}
                    {loading && <p className="loading">Loading...</p>}
                </div>
            )}

            {activeTab === "events" && (
                <div className="events-container">
                    <button className="add-event-btn" onClick={() => setShowEventForm(true)}>+ Add Event</button>

                    {events === null || events.length === 0 ? (
                        <p className="no-events">No events available.</p>
                    ) : (
                        events.map((event) => (
                            <div key={event.id} className="event-card">
                                <h3>{event.title}</h3>
                                <p>{event.description}</p>
                                <p className="event-date">Date: {new Date(event.date).toLocaleDateString()}</p>
                                <p>Going: {event.going_count} | Not Going: {event.not_going_count}</p>
                                <div className="event-actions">
                                    <button
                                        className={event.user_response === "going" ? "active" : ""}
                                        onClick={() => handleGoing(event.id, "going")}>
                                        Going
                                    </button>
                                    <button
                                        className={event.user_response === "not_going" ? "active" : ""}
                                        onClick={() => handleGoing(event.id, "not_going")}>
                                        Not Going
                                    </button>
                                </div>
                            </div>
                        ))
                    )}
                    {hasMoreEvents && <div id="load-more-events" style={{ height: "10px" }} />}
                    {loadingEvents && <p className="loading">Loading events...</p>}
                </div>
            )}

            {showEventForm && (
                <div className="popup">
                    <div className="popup-content">
                        <h2>Create Event</h2>
                        <input type="text" maxLength={50} placeholder="Title" value={newEvent.title} onChange={(e) => setNewEvent({ ...newEvent, title: e.target.value })} />
                        <textarea placeholder="Description" maxLength={200} value={newEvent.description} onChange={(e) => setNewEvent({ ...newEvent, description: e.target.value })} />
                        <input type="datetime-local" min={new Date().toISOString().slice(0, 16)} value={newEvent.date} onChange={(e) => setNewEvent({ ...newEvent, date: e.target.value })} />
                        <div className="radio-group">

                            <input type="radio" id="going" value="going" checked={newEvent.user_response === "going"} onChange={() => setNewEvent({ ...newEvent, user_response: "going" })} />
                            <label htmlFor="going">Going</label>

                            <input type="radio" id="not-going" value="not_going" checked={newEvent.user_response === "not_going"} onChange={() => setNewEvent({ ...newEvent, user_response: "not_going" })} />
                            <label htmlFor="not-going">Not Going</label>
                        </div>

                        <div className="popup-actions">
                            <button onClick={handleCreateEvent}>Create</button>
                            <button onClick={() => setShowEventForm(false)}>Cancel</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}