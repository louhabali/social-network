import { cookies } from 'next/headers';
import Post from "@/components/post";
import CommentsClient from "./CommentsClient";

async function getPost(postId) {
  const cookieStore = cookies();
  const sessionToken = cookieStore.get('Token')?.value;
  
  const res = await fetch(`http://localhost:8080/api/posts/${postId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Cookie: sessionToken ? `Token=${sessionToken}` : "",
    },
  });
  if (!res.ok) return null;
  return res.json();
}


export default async function PostPage({ params }) {

  const post = await getPost(params.postId);
  
  return (
    <div>
      <Post post={post} />      
      <CommentsClient postId={params.postId} lenComments={post.NbComment}/>
    </div>
  );
}
