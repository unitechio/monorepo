import { Block } from "@/types/page";
import Hero from "./blocks/Hero";
import Features from "./blocks/Features"; // We will create this next

type BlockComponent = React.ComponentType<{ data: Block }>;

const BLOCK_REGISTRY: Record<string, BlockComponent> = {
    hero: Hero,
    features: Features,
};

export default function BlockRenderer({ block }: { block: Block }) {
    const Component = BLOCK_REGISTRY[block.type];

    if (!Component) {
        console.warn(`No component found for block type: ${block.type}`);
        return null;
    }

    return <Component data={block} />;
}
