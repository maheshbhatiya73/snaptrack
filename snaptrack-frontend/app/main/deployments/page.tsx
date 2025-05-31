"use client"
import { useState, useEffect } from "react";

export default function Home() {
  const [deployments, setDeployments] = useState<any[]>([]);
  const [form, setForm] = useState({ appName: "", userName: "", deployPath: "" });
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
  fetch("http://localhost:8000/api/deployments")
    .then(res => res.json())
    .then(data => {
      if (Array.isArray(data)) {
        setDeployments(data);
      } else {
        setDeployments([]);
      }
    })
    .catch(() => setDeployments([])); 
}, []);


  const handleSubmit = async (e: any) => {
    e.preventDefault();
    if (!file) return;

    const formData = new FormData();
    formData.append("appName", form.appName);
    formData.append("userName", form.userName);
    formData.append("deployPath", form.deployPath);
    formData.append("file", file);

    setLoading(true);
    const res = await fetch("http://localhost:8000/api/deployments", {
      method: "POST",
      body: formData,
    });

    if (res.ok) {
      const newDep = await res.json();
      setDeployments(prev => [newDep, ...prev]);
      setForm({ appName: "", userName: "", deployPath: "" });
      setFile(null);
    } else {
      const err = await res.text();
      alert("Failed: " + err);
    }
    setLoading(false);
  };

  return (
    <main className="max-w-4xl mx-auto py-8 px-4">
      <h1 className="text-3xl font-bold mb-6">🚀 Deployment Manager</h1>

      <form onSubmit={handleSubmit} className="bg-white p-4 rounded shadow space-y-4">
        <input
          className="w-full border p-2 rounded"
          placeholder="App Name"
          value={form.appName}
          onChange={e => setForm({ ...form, appName: e.target.value })}
          required
        />
        <input
          className="w-full border p-2 rounded"
          placeholder="User Name"
          value={form.userName}
          onChange={e => setForm({ ...form, userName: e.target.value })}
          required
        />
        <input
          className="w-full border p-2 rounded"
          placeholder="Deploy Path (must exist on server)"
          value={form.deployPath}
          onChange={e => setForm({ ...form, deployPath: e.target.value })}
          required
        />
        <input
          type="file"
          className="w-full"
          onChange={e => setFile(e.target.files?.[0] || null)}
          required
        />
        <button
          disabled={loading}
          type="submit"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          {loading ? "Deploying..." : "Deploy File"}
        </button>
      </form>

      <h2 className="text-xl font-semibold mt-8 mb-4">📋 Deployed Apps</h2>
      <div className="grid gap-4">
        {deployments.map(dep => (
          <div key={dep.id} className="bg-gray-100 p-4 rounded shadow">
            <h3 className="font-bold text-lg">{dep.appName}</h3>
            <p>User: {dep.userName}</p>
            <p>File: {dep.fileName}</p>
            <p>Path: {dep.deployPath}</p>
            <p className="text-sm text-gray-500">
              {new Date(dep.createdAt).toLocaleString()}
            </p>
          </div>
        ))}
        {deployments.length === 0 && (
          <p className="text-gray-500">No deployments yet.</p>
        )}
      </div>
    </main>
  );
}
