"use client";

import { useState } from "react";
import Post from "@/components/post";
import Image from "next/image";

export default function ProfileClient({ profile, userId }) {
  const [isPublic, setIsPublic] = useState(profile.user.IsPublic);
  const [followStatus, setFollowStatus] = useState(profile.followStatus);
  const [posts, setPosts] = useState(profile.posts || []);
  const [activeTab, setActiveTab] = useState("posts");
  const isOwnProfile = profile.user.ID === profile.userid;
  const canViewProfile = isPublic || profile.followStatus === "approved" || isOwnProfile;

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

  async function onFollowToggle() {
    const method = followStatus != "" ? "DELETE" : "POST";
    const res = await fetch(`http://localhost:8080/api/profile/${userId}/follow`, {
      method,
      credentials: "include",
    });

    if (res.ok) {
      if (followStatus != "") {
        setFollowStatus("");
      } else {
        setFollowStatus("pending");
      }
    }
  }

  const handleLike = async (e) => {

    const pd = e.target.closest('section').id
    const likeIcon = e.target

    const postIndex = posts.findIndex(post => post.Pid.toString() === pd)
    const uPosts = [...posts]

    let lk = uPosts[postIndex].UserInteraction === 1 ? 0 : 1



    const response = await fetch(`http://localhost:8080/api/like`, {
      method: 'POST',
      body: JSON.stringify({ pd, lk }),
      credentials: "include",
    });

    if (response.ok) {

      uPosts[postIndex].UserInteraction = lk
      uPosts[postIndex].Likes += lk ? 1 : -1
      likeIcon.style.color = lk ? 'red' : 'white'

      setPosts(uPosts)
    }
  }

  console.log(profile);
  

  return (
    <div>
      <div className="profile-header">
        <div className="profile-info">
          {profile.user.Avatar ? (
            <Image src={`data:image/png;base64,${profile.user.Avatar}`} alt="Avatar" width={64} height={64} />
          ) : (
            <Image src="/image/no-avatar.png" alt="No Avatar" width={64} height={64} />
          )}
          <h1>{profile.user.Firstname} {profile.user.Lastname}</h1>
          {!canViewProfile ? (
            <p>This profile is private. Follow to see more.</p>
          ) : (
            <>
              <p><strong>Email:</strong> {profile.user.Email}</p>
              <p>{profile.user.AboutMe || "No bio available"}</p>
              <p><strong>Date of Birth:</strong> {profile.user.DateOfBirth}</p>
            </>
          )}
          <p><strong>Followers:</strong> {profile.followers?.length || 0}</p>
          <p><strong>Following:</strong> {profile.following?.length || 0}</p>
          <p><strong>Posts:</strong> {profile.posts?.length || 0}</p>

          {isOwnProfile ? (
            <button onClick={handleProfilePrivacyToggle}>
              {isPublic ? "Switch to Private Profile" : "Switch to Public Profile"}
            </button>
          ) : (
            followStatus === "pending" ? (
              <>
                <span>Request Pending </span>
                <button onClick={onFollowToggle} className="cancel-button">
                  Cancel Request
                </button>
              </>
            ) : followStatus === "approved" ? (
              <button onClick={onFollowToggle} className="follow-button">Unfollow</button>
            ) : (
              <button onClick={onFollowToggle} className="follow-button">Follow</button>
            )
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
              {(profile.posts?.length || 0) === 0 ? <p>No posts yet!</p> : profile.posts.map((post) => <Post key={post.Pid} post={post} handleLike={handleLike} />)}
            </div>
          )}

          {activeTab === "followers" && (
            <div className="profile-followers">
              <h3>Followers</h3>
              {profile.followers?.map((follower) => (
                <div key={follower.ID}>{follower.Firstname} {follower.Lastname}</div>
              ))}
            </div>
          )}

          {activeTab === "following" && (
            <div className="profile-following">
              <h3>Following</h3>
              {profile.following?.map((followingUser) => (
                <div key={followingUser.ID}>{followingUser.Firstname} {followingUser.Lastname}</div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}