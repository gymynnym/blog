import rss from '@astrojs/rss';
import type { APIRoute, GetStaticPaths } from 'astro';
import { isLocale, locales, ui, type Locale } from '@/lib/i18n';
import { postRoute, publishedPosts } from '@/lib/posts';
import { postPath } from '@/lib/routes';

const rssLanguages: Record<Locale, string> = {
  ko: 'ko-KR',
  en: 'en-US',
  ja: 'ja-JP',
};

export const getStaticPaths = (() => locales.map((lang) => ({ params: { lang } }))) satisfies GetStaticPaths;

export const GET: APIRoute = async ({ params, site }) => {
  if (!isLocale(params.lang)) {
    return new Response('Not found', { status: 404 });
  }

  const locale = params.lang;
  const posts = await publishedPosts(locale);

  return rss({
    title: ui[locale].feedTitle,
    description: ui[locale].blogDescription,
    site: site!,
    customData: `<language>${rssLanguages[locale]}</language>`,
    items: posts.map((post) => {
      const { slug } = postRoute(post);

      return {
        title: post.data.title,
        description: post.data.description,
        pubDate: post.data.publishedAt,
        link: postPath(locale, slug),
      };
    }),
  });
};
