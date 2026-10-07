import type { APIRoute } from 'astro';
import { getCollection } from 'astro:content';
import { locales, type LanguageLinks, type Locale } from '@/lib/i18n';
import { postRoute, publishedPosts, translationLinks } from '@/lib/posts';
import { homePath, postPath, postsPath, projectsPath, seriesPagePath, topicPath } from '@/lib/routes';

interface SitemapEntry {
  path: string;
  lastmod?: Date;
  languageLinks: LanguageLinks;
}

function localizedPages(pathForLocale: (locale: Locale) => string): SitemapEntry[] {
  const languageLinks = Object.fromEntries(locales.map((locale) => [locale, pathForLocale(locale)]));
  return locales.map((locale) => ({ path: pathForLocale(locale), languageLinks }));
}

function escapeXml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&apos;');
}

export const GET: APIRoute = async ({ site }) => {
  if (!site) {
    throw new Error('sitemap.xml을 생성하려면 astro.config.mjs의 site 설정이 필요합니다.');
  }

  const [posts, topics, series] = await Promise.all([
    publishedPosts(),
    getCollection('topics'),
    getCollection('series'),
  ]);
  const usedTopicIds = new Set(posts.flatMap((post) => post.data.topics.map((topic) => topic.id)));
  const usedSeriesIds = new Set(posts.flatMap((post) => (post.data.series ? [post.data.series.id] : [])));

  // Use localized URLs instead of the duplicate / and /projects pages.
  const entries: SitemapEntry[] = [
    ...localizedPages(homePath),
    ...localizedPages(postsPath),
    ...localizedPages(projectsPath),
    ...posts.map((post) => {
      const { locale, slug } = postRoute(post);
      return {
        path: postPath(locale, slug),
        lastmod: post.data.updatedAt ?? post.data.publishedAt,
        languageLinks: translationLinks(post, posts),
      };
    }),
    ...topics
      .filter((topic) => usedTopicIds.has(topic.id))
      .flatMap((topic) => localizedPages((locale) => topicPath(locale, topic.data.slug))),
    ...series
      .filter((item) => usedSeriesIds.has(item.id))
      .flatMap((item) => localizedPages((locale) => seriesPagePath(locale, item.data.slug))),
  ];

  const absoluteUrl = (path: string) => escapeXml(new URL(path, site).href);
  const urls = entries
    .sort((a, b) => a.path.localeCompare(b.path))
    .map(({ path, lastmod, languageLinks }) => {
      const alternates = Object.entries(languageLinks).map(
        ([locale, href]) => `    <xhtml:link rel="alternate" hreflang="${locale}" href="${absoluteUrl(href)}" />`,
      );
      return [
        '  <url>',
        `    <loc>${absoluteUrl(path)}</loc>`,
        ...(lastmod ? [`    <lastmod>${lastmod.toISOString()}</lastmod>`] : []),
        ...alternates,
        '  </url>',
      ].join('\n');
    });
  const xml = [
    '<?xml version="1.0" encoding="UTF-8"?>',
    '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">',
    ...urls,
    '</urlset>',
    '',
  ].join('\n');

  return new Response(xml, { headers: { 'Content-Type': 'application/xml; charset=utf-8' } });
};
