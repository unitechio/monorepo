export interface Item {
    id: string;
    title?: string;
    subtitle?: string;
    image?: string;
    link?: string;
    metadata?: Record<string, any>;
}

export interface Block {
    id: string;
    type: string; // 'hero', 'features', etc.
    title?: string;
    items?: Item[];
    properties?: Record<string, any>;
}

export interface Page {
    slug: string;
    title: string;
    config?: Record<string, any>;
    blocks: Block[];
}
