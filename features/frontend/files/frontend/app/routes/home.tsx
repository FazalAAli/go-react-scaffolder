import { useState } from "react";
import { client } from "~/lib/client";

export default function Home() {
  const [name, setName] = useState("");
  const [greeting, setGreeting] = useState("");
  const [error, setError] = useState("");

  async function onGreet() {
    setError("");
    try {
      const res = await client.greet({ name });
      setGreeting(res.greeting);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-4 p-8">
      <h1 className="text-2xl font-semibold">ConnectRPC Boilerplate</h1>
      <div className="flex gap-2">
        <input
          className="rounded border px-3 py-2"
          placeholder="your name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <button className="rounded bg-black px-4 py-2 text-white" onClick={onGreet}>
          Greet
        </button>
      </div>
      {greeting && <p className="text-lg">{greeting}</p>}
      {error && <p className="text-red-600">{error}</p>}
    </main>
  );
}
