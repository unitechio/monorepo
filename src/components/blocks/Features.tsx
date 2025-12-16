import { Block } from "@/types/page";

export default function Features({ data }: { data: Block }) {
    return (
        <section className="py-16 bg-white dark:bg-black">
            <div className="container mx-auto px-4">
                {data.title && (
                    <h2 className="text-3xl font-bold text-center mb-12">{data.title}</h2>
                )}

                <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
                    {data.items?.map((item) => (
                        <div key={item.id} className="p-6 border rounded-lg hover:shadow-lg transition">
                            <h3 className="text-xl font-semibold mb-2">{item.title}</h3>
                            <p className="text-gray-600 dark:text-gray-400">{item.subtitle}</p>
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
}
