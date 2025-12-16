import { Page } from "@/types/page";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

export async function fetchPage(slug: string, lang: string = "en"): Promise<Page | null> {
    try {
        const res = await fetch(`${API_URL}/pages/${slug}?lang=${lang}`, {
            cache: 'no-store' // Ensure fresh data since we are fetching client-side
        });

        if (!res.ok) {
            console.error("Failed to fetch page:", res.statusText);
            return null;
        }

        return await res.json();
    } catch (error) {
        console.error("Error fetching page:", error);
        return null;
    }
}
