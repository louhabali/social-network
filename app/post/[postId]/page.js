import { cookies } from 'next/headers';
import Post from "@/components/post";
import CommentsClient from "./CommentsClient";
// import "./post.css"

async function getPost(postId) {
  const cookieStore = cookies();
  const sessionToken = cookieStore.get('Token')?.value;
  
  const res = await fetch(`http://localhost:8080/api/posts/${postId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Cookie: sessionToken ? `Token=${sessionToken}` : "",
    },
    cache: "no-store",
  });
  if (!res.ok) return null;
  return res.json();
}

async function getComments(postId, offset = 0) {
  const cookieStore = cookies();
  const sessionToken = cookieStore.get('Token')?.value;
  
  const res = await fetch(`http://localhost:8080/api/posts/${postId}/comments?offset=${offset}`, {
    headers: {
      "Content-Type": "application/json",
      Cookie: sessionToken ? `Token=${sessionToken}` : "",
    },
  });

  if (!res.ok) return [];
  return res.json();
}

export default async function PostPage({ params }) {

  const post = await getPost(params.postId);
  const initialComments = await getComments(params.postId);

  // console.log(post);
  
  return (
    <div>
      <Post post={post} />
      
      <CommentsClient postId={params.postId} lenComments={post.NbComment} initialComments={initialComments} />
    </div>
  );
}
