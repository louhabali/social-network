"use client";

import { useState, useRef, useEffect, useCallback } from "react";
import Post from "@/components/post";
import Image from "next/image";
import "../../styles/profile.css";

export default function ProfileClient({ profile, userId }) {
  const [isPublic, setIsPublic] = useState(profile.user.IsPublic);
  const [followStatus, setFollowStatus] = useState(profile.followStatus);
  const [posts, setPosts] = useState(profile.posts || []);
  const [activeTab, setActiveTab] = useState("posts");
  const [offset, setOffset] = useState(posts.length);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const observer = useRef(null);

  const isOwnProfile = profile.user.ID === profile.userid;
  const canViewProfile = isPublic || followStatus === "approved" || isOwnProfile;

  async function handleProfilePrivacyToggle() {
    const newStatus = !isPublic;
    const res = await fetch(`http://localhost:8080/api/profile/${userId}/privacy`, {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({ isPublic: newStatus }),
    });

    if (res.ok) {
      setIsPublic(newStatus);
    }
  }

  const fetchPosts = async () => {
    if (loading || !hasMore) return;

    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/api/posts?offset=${offset}&userId=${userId}`, {
        credentials: "include",
      });

      if (!response.ok) {
        setHasMore(false);
        return;
      }

      const data = await response.json();
      if (!data || data.length === 0) {
        setHasMore(false);
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
    const method = followStatus !== "" ? "DELETE" : "POST";
    const res = await fetch(`http://localhost:8080/api/profile/${userId}/follow`, {
      method,
      credentials: "include",
    });

    if (res.ok) {
      setFollowStatus(followStatus !== "" ? "" : "pending");
    }
  }

  return (
    <div className="profile-container">
      <div className="profile-header">
        <div className="profile-cover-photo"></div>
        <div className="profile-info">
          <div className="profile-avatar">
            {profile.user.Avatar ? (
              <Image src={`data:image/png;base64,${profile.user.Avatar}`} alt="Avatar" width={120} height={120} />
            ) : (
              <Image src="/image/no-avatar.png" alt="No Avatar" width={120} height={120} />
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
                posts.map((post) => <Post key={post.Pid} post={post}/>)
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
                    <div key={follower.ID} className="follower-item user-item">
                      <a href={`/profile/${follower.ID}`} className="user-link">
                        {follower.Firstname} {follower.Lastname}
                      </a>
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
                    <div key={followingUser.ID} className="following-item user-item">
                      <a href={`/profile/${followingUser.ID}`} className="user-link">
                        {followingUser.Firstname} {followingUser.Lastname}
                      </a>
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
