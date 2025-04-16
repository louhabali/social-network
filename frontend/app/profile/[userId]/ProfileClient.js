"use client";

import { useState, useRef, useEffect, useContext } from "react";
import Post from "@/components/post";
import Image from "next/image";
import Link from "next/link"
import "../../styles/profile.css";
import { WebSocketContext } from "@/context/WebSocketContext";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function ProfileClient({ profile, userId }) {
  const [isPublic, setIsPublic] = useState(profile.user.IsPublic);
  const [followStatus, setFollowStatus] = useState(profile.followStatus);
  const [posts, setPosts] = useState(profile.posts || []);
  const [activeTab, setActiveTab] = useState("posts");
  const [offset, setOffset] = useState(posts.length);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const observer = useRef(null);
  const { showNotification } = useContext(WebSocketContext);

  const isOwnProfile = profile.user.ID === profile.userid;
  const canViewProfile = isPublic || followStatus === "approved" || isOwnProfile;
  console.log(profile);
  
  async function handleProfilePrivacyToggle() {
    const newStatus = !isPublic;
    const res = await fetch(`${apiUrl}/api/profile/${userId}/privacy`, {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({ isPublic: newStatus }),
    });

    if (!res.ok) {
      showNotification("Error updating privacy settings", "error");
      return;
    }
    setIsPublic(newStatus);

  }

  const fetchPosts = async () => {
    if (loading || !hasMore) return;

    setLoading(true);
    try {
      const response = await fetch(`${apiUrl}/api/posts?offset=${offset}&userId=${userId}`, {
        credentials: "include",
      });

      if (!response.ok) {
        showNotification("Error loading posts", "error");
        return;
      }

      const data = await response.json();
      if (!data || data.length === 0) {
        setHasMore(false);
        setLoading(false);
        return;
      }

      setPosts((prevPosts) => [...prevPosts, ...data]);
      setOffset((prevOffset) => prevOffset + data.length);
      if (data.length < 4) setHasMore(false);
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

    const loadMoreElement = document.getElementById("load-more");
    if (observer.current && loadMoreElement) {
      observer.current.observe(loadMoreElement);
    }

    return () => {
      observer.current?.disconnect();
    };
  }, [fetchPosts, loading, hasMore]);

  async function onFollowToggle() {
    const method = followStatus ? "DELETE" : "POST";
    const res = await fetch(`${apiUrl}/api/profile/${userId}/follow`, {
      method,
      credentials: "include",
    });

    if (!res.ok) {
      showNotification("Error updating follow status", "error");
      return;
    }
    if (method === "DELETE") {
      if (followStatus === 'approved') {
        profile.followers = profile.followers.filter(f => f.id !== profile.userid);
      }
      setFollowStatus("");
    } else if (method === "POST") {
      const p = await res.json().catch(() => null);
      if (p) {
        profile.followers = profile.followers ? [...profile.followers, {
          id: p.ID,
          fname: p.Firstname,
          lname: p.Lastname
        }] : [{
          id: p.ID,
          fname: p.Firstname,
          lname: p.Lastname
        }];
      }
      setFollowStatus(p ? "approved" : "pending");
    }

  }

  return (
    <div className="profile-container">
      <div className="profile-header">
        <div className="profile-cover-photo"></div>
        <div className="profile-info">
          <div className="profile-avatar">
            {profile.user.Avatar ? (
              <img src={`${apiUrl}/api/images/${profile.user.Avatar}`} alt="Avatar" width={120} height={120} />
            ) : (
              <img src="/image/no-avatar.png" alt="No Avatar" width={120} height={120} />
            )}
          </div>
          <h1 className="profile-name">
            {profile.user.Firstname} {profile.user.Lastname}
          </h1>
          {!canViewProfile ? (
            <p className="profile-private-notice">This profile is private. Follow to see more.</p>
          ) : (
            <div className="profile-details">
              <p className="profile-email">
                <strong>Email:</strong> {profile.user.Email}
              </p>
              <p className="profile-bio">{profile.user.AboutMe || "No bio available"}</p>
              <p className="profile-dob">
                <strong>Date of Birth:</strong> {profile.user.DateOfBirth}
              </p>
            </div>
          )}

          <div className="profile-stats">
            <div className="stat-item">
              <span className="stat-count">{profile.followers?.length || 0}</span>
              <span className="stat-label">Followers</span>
            </div>
            <div className="stat-item">
              <span className="stat-count">{profile.following?.length || 0}</span>
              <span className="stat-label">Following</span>
            </div>
            <div className="stat-item">
              <span className="stat-count">{profile.totalPosts}</span>
              <span className="stat-label">Posts</span>
            </div>
          </div>

          <div className="profile-actions">
            {isOwnProfile ? (
              <button onClick={handleProfilePrivacyToggle} className="profile-button privacy-button">
                {isPublic ? "Switch to Private Profile" : "Switch to Public Profile"}
              </button>
            ) : followStatus === "pending" ? (
              <div className="follow-actions">
                <span className="pending-notice">Request Pending </span>
                <button onClick={onFollowToggle} className="profile-button cancel-button">
                  Cancel Request
                </button>
              </div>
            ) : followStatus === "approved" ? (
              <button onClick={onFollowToggle} className="profile-button follow-button unfollow">
                Unfollow
              </button>
            ) : (
              <button onClick={onFollowToggle} className="profile-button follow-button">
                Follow
              </button>
            )}
          </div>
        </div>

        {canViewProfile && (
          <div className="profile-tabs">
            <button className={`tab-button ${activeTab === "posts" ? "active" : ""}`} onClick={() => setActiveTab("posts")}>
              Posts
            </button>
            <button className={`tab-button ${activeTab === "followers" ? "active" : ""}`} onClick={() => setActiveTab("followers")}>
              Followers
            </button>
            <button className={`tab-button ${activeTab === "following" ? "active" : ""}`} onClick={() => setActiveTab("following")}>
              Following
            </button>
          </div>
        )}
      </div>

      {canViewProfile && (
        <div className="profile-content">
          {activeTab === "posts" && (
            <div className="profile-posts">
              {posts.length === 0 ? (
                <p className="no-content">No posts yet!</p>
              ) : (
                posts.map((post) => <Post key={post.Pid} post={post} />)
              )}
              {hasMore && <div id="load-more" className="load-more" />}
              {loading && <p className="loading">Loading...</p>}
            </div>
          )}

          {activeTab === "followers" && (
            <div className="followers-section">
              <h3 className="section-title">Followers</h3>
              {profile.followers?.length > 0 ? (
                <div className="user-list">
                  {profile.followers.map((follower) => (
                    <div key={follower.id} className="follower-item user-item">
                      <Link href={`/profile/${follower.id}`} className="user-link">
                        {follower.fname} {follower.lname}
                      </Link>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="no-content">No followers yet.</p>
              )}
            </div>
          )}

          {activeTab === "following" && (
            <div className="following-section">
              <h3 className="section-title">Following</h3>
              {profile.following?.length > 0 ? (
                <div className="user-list">
                  {profile.following.map((followingUser) => (
                    <div key={followingUser.id} className="following-item user-item">
                      <Link href={`/profile/${followingUser.id}`} className="user-link">
                        {followingUser.fname} {followingUser.lname}
                      </Link>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="no-content">Not following anyone.</p>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
