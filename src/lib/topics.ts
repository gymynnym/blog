import { getEntries } from 'astro:content';
import type { CollectionEntry } from 'astro:content';
import type { Post } from '@/lib/posts';

export interface Topic {
  slug: string;
  name: string;
}

export function topicFromEntry({ data: { slug, name } }: CollectionEntry<'topics'>): Topic {
  return { slug, name };
}

export async function resolveTopics(references: Post['data']['topics']): Promise<Topic[]> {
  return (await getEntries(references)).map(topicFromEntry);
}
