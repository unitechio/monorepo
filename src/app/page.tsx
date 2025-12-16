"use client";

import { useEffect, useState } from "react";
import { fetchPage } from "@/lib/api";
import { Page } from "@/types/page";
import BlockRenderer from "@/components/BlockRenderer";
import { useSearchParams, useRouter } from "next/navigation";

export default function Home() {
  const [data, setData] = useState<Page | null>(null);
  const [loading, setLoading] = useState(true);

  const searchParams = useSearchParams();
  const router = useRouter();

  const lang = searchParams.get("lang") || "en";

  useEffect(() => {
    async function loadData() {
      setLoading(true);
      const pageData = await fetchPage("home", lang);
      setData(pageData);
      setLoading(false);
    }
    loadData();
  }, [lang]);

  const switchLang = (newLang: string) => {
    router.push(`/?lang=${newLang}`);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="flex items-center justify-center min-h-screen text-red-500">
        Failed to load page content.
      </div>
    );
  }

  return (
    <main className="min-h-screen">
      {/* Language Switcher */}
      <div className="fixed top-4 right-4 z-50 bg-white dark:bg-gray-800 p-2 rounded shadow-md flex gap-2">
        <button
          onClick={() => switchLang("en")}
          className={`px-3 py-1 rounded ${lang === "en" ? "bg-blue-600 text-white" : "bg-gray-200 dark:bg-gray-700 hover:bg-gray-300"}`}
        >
          English
        </button>
        <button
          onClick={() => switchLang("vi")}
          className={`px-3 py-1 rounded ${lang === "vi" ? "bg-blue-600 text-white" : "bg-gray-200 dark:bg-gray-700 hover:bg-gray-300"}`}
        >
          Tiếng Việt
        </button>
      </div>

      {data.blocks.map((block) => (
        <BlockRenderer key={block.id} block={block} />
      ))}
    </main>
  );
}
