import { getCollection } from 'astro:content';
import type { CollectionEntry } from 'astro:content';
import { isLocale, type LanguageLinks, type Locale } from '@/lib/i18n';
import { postPath } from '@/lib/routes';

export type Post = CollectionEntry<'posts'>;

export function postRoute(post: Post): { locale: Locale; slug: string } {
  const [locale, ...slugParts] = post.id.split('/');
  const slug = slugParts.join('/');

  if (!isLocale(locale) || !slug) {
    throw new Error(`포스트 "${post.id}"는 src/content/posts/{ko|en|ja}/{slug}.mdx 경로에 있어야 합니다.`);
  }

  return { locale, slug };
}

export async function publishedPosts(locale?: Locale): Promise<Post[]> {
  const posts = await getCollection('posts', ({ data }) => !data.draft);
  const seenTranslations = new Set<string>();

  for (const post of posts) {
    const { locale, slug } = postRoute(post);
    const translation = `${slug}:${locale}`;

    if (seenTranslations.has(translation)) {
      throw new Error(`slug "${slug}"의 ${locale} 포스트가 중복되었습니다.`);
    }
    seenTranslations.add(translation);
  }

  return posts
    .filter((post) => !locale || postRoute(post).locale === locale)
    .sort((a, b) => b.data.publishedAt.getTime() - a.data.publishedAt.getTime());
}

export function translationLinks(post: Post, posts: Post[]): LanguageLinks {
  const { slug: currentSlug } = postRoute(post);

  return Object.fromEntries(
    posts
      .filter((candidate) => postRoute(candidate).slug === currentSlug)
      .map((candidate) => {
        const { locale, slug } = postRoute(candidate);
        return [locale, postPath(locale, slug)];
      }),
  );
}
