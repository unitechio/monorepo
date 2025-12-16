import { Block } from "@/types/page";
import Image from "next/image";
import Link from "next/link";

export default function Hero({ data }: { data: Block }) {
    const item = data.items?.[0];

    if (!item) return null;

    return (
        <section className="bg-gray-100 dark:bg-gray-900 py-20 text-center">
            <div className="container mx-auto px-4">
                <h1 className="text-5xl font-bold mb-4">{data.title || item.title}</h1>
                <p className="text-xl mb-8 text-gray-600 dark:text-gray-300">{item.subtitle}</p>

                {item.image && (
                    <div className="relative w-full max-w-2xl mx-auto h-64 mb-8">
                        {/* Using standard img for simplicity in dev/export if remote images not config */}
                        {/* eslint-disable-next-line @next/next/no-img-element */}
                        <img src={item.image} alt={item.title || "Hero"} className="object-cover rounded-lg w-full h-full" />
                    </div>
                )}

                {item.link && (
                    <Link
                        href={item.link}
                        className="inline-block bg-blue-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-blue-700 transition"
                    >
                        {item.metadata?.buttonText || "Learn More"}
                    </Link>
                )}
            </div>
        </section>
    );
}
