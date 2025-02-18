

export default async function Home() {
    const f = await fetch('/api/cookie', {
      method: 'POST',}
    );
    const res = f.json();
    console.log(res);
    

  return <div>Home Page</div>;
}
