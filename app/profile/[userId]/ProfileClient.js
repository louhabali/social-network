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

  const handleLike = async (e) => {
    const pd = e.target.closest("section").id;
    const postIndex = posts.findIndex((post) => post.Pid.toString() === pd);
    if (postIndex === -1) return;

    let updatedPosts = [...posts];
    let newLikeStatus = updatedPosts[postIndex].UserInteraction === 1 ? 0 : 1;

    const response = await fetch(`http://localhost:8080/api/like`, {
      method: "POST",
      body: JSON.stringify({ pd, lk: newLikeStatus }),
      credentials: "include",
    });

    if (response.ok) {
      updatedPosts[postIndex].UserInteraction = newLikeStatus;
      updatedPosts[postIndex].Likes += newLikeStatus ? 1 : -1;
      setPosts(updatedPosts);
    }
  };

  return (
    <div>
      <div className="profile-header">
        <div className="profile-info">
          {profile.user.Avatar ? (
            <Image src={`data:image/png;base64,${profile.user.Avatar}`} alt="Avatar" width={64} height={64} />
          ) : (
            <Image src="/image/no-avatar.png" alt="No Avatar" width={64} height={64} />
          )}
          <h1>
            {profile.user.Firstname} {profile.user.Lastname}
          </h1>
          {!canViewProfile ? (
            <p>This profile is private. Follow to see more.</p>
          ) : (
            <>
              <p>
                <strong>Email:</strong> {profile.user.Email}
              </p>
              <p>{profile.user.AboutMe || "No bio available"}</p>
              <p>
                <strong>Date of Birth:</strong> {profile.user.DateOfBirth}
              </p>
            </>
          )}
          <p>
            <strong>Followers:</strong> {profile.followers?.length || 0}
          </p>
          <p>
            <strong>Following:</strong> {profile.following?.length || 0}
          </p>
          <p>
            <strong>Posts:</strong> {profile.totalPosts}
          </p>

          {isOwnProfile ? (
            <button onClick={handleProfilePrivacyToggle}>
              {isPublic ? "Switch to Private Profile" : "Switch to Public Profile"}
            </button>
          ) : followStatus === "pending" ? (
            <>
              <span>Request Pending </span>
              <button onClick={onFollowToggle} className="cancel-button">
                Cancel Request
              </button>
            </>
          ) : followStatus === "approved" ? (
            <button onClick={onFollowToggle} className="follow-button">
              Unfollow
            </button>
          ) : (
            <button onClick={onFollowToggle} className="follow-button">
              Follow
            </button>
          )}
        </div>

        {canViewProfile && (
          <div className="profile-tabs">
            <button className={activeTab === "posts" ? "active" : ""} onClick={() => setActiveTab("posts")}>
              Posts
            </button>
            <button className={activeTab === "followers" ? "active" : ""} onClick={() => setActiveTab("followers")}>
              Followers
            </button>
            <button className={activeTab === "following" ? "active" : ""} onClick={() => setActiveTab("following")}>
              Following
            </button>
          </div>
        )}
      </div>
      {canViewProfile && (
        <div className="profile-content">
          {activeTab === "posts" && (
            <div className="profile-posts">
              {posts.length === 0 ? <p>No posts yet!</p> : posts.map((post) => <Post key={post.Pid} post={post} handleLike={handleLike} />)}
              {hasMore && <div id="load-more" style={{ height: "10px" }} />}
              {loading && <p>Loading...</p>}
            </div>
          )}

          {activeTab === "followers" && (
            <div>
              <h3>Followers</h3>
              {profile.followers?.length > 0 ? (
                profile.followers.map((follower) => (
                  <div key={follower.ID} className="follower-item">
                    <a href={`/profile/${follower.ID}`}>
                      {follower.Firstname} {follower.Lastname}
                    </a>
                  </div>
                ))
              ) : (
                <p>No followers yet.</p>
              )}
            </div>
          )}

          {activeTab === "following" && (
            <div>
              <h3>Following</h3>
              {profile.following?.length > 0 ? (
                profile.following.map((followingUser) => (
                  <div key={followingUser.ID} className="following-item">
                    <a href={`/profile/${followingUser.ID}`}>
                      {followingUser.Firstname} {followingUser.Lastname}
                    </a>
                  </div>
                ))
              ) : (
                <p>Not following anyone.</p>
              )}
            </div>
          )}

        </div>
      )}
    </div>
  );
}
