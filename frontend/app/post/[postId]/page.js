import { cookies } from 'next/headers';
import Post from "@/components/post";
import CommentsClient from "./CommentsClient";
import ErrorPage from "@/components/ErrorPage";

// import "./post.css"

async function getPost(postId) {
  const cookieStore = await cookies();
  const sessionToken = cookieStore.get('Token')?.value;

  const res = await fetch(`http://localhost:8080/api/posts/${postId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Cookie: sessionToken ? `Token=${sessionToken}` : "",
    },
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error("Internal server error");
  }
  return res.json();
}

async function getComments(postId, offset = 0) {
  const cookieStore = await cookies();
  const sessionToken = cookieStore.get('Token')?.value;

  const res = await fetch(`http://localhost:8080/api/posts/${postId}/comments?offset=${offset}`, {
    headers: {
      "Content-Type": "application/json",
      Cookie: sessionToken ? `Token=${sessionToken}` : "",
    },
  });

  if (!res.ok) {
    throw new Error("Internal server error");
  }
  return res.json();
}

export default async function PostPage({ params: paramsPromise }) {
  const params = await paramsPromise;

  try {
    const post = await getPost(params.postId);
    const initialComments = await getComments(params.postId);

    // console.log(post);

    return (

      <div className='comment-div'>
        <div className='post-comment'>
          <Post post={post} />
        </div>
        <CommentsClient postId={params.postId} lenComments={post.NbComment} initialComments={initialComments} />
      </div>
    );
  } catch (error) {
    return <ErrorPage error={error} />;

  }
}
