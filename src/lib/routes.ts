import type { Locale } from '@/lib/i18n';

export const homePath = (locale: Locale) => `/${locale}`;
export const postsPath = (locale: Locale) => `/${locale}/posts`;
export const projectsPath = (locale: Locale) => `/${locale}/projects`;
export const feedPath = (locale: Locale) => `/${locale}/feed`;
export const postPath = (locale: Locale, slug: string) => `${postsPath(locale)}/${slug}`;
export const topicPath = (locale: Locale, slug: string) => `${postsPath(locale)}/topics/${slug}`;
export const seriesPagePath = (locale: Locale, slug: string) => `${postsPath(locale)}/series/${slug}`;
export const seriesPath = (locale: Locale, slug: string) => `${postsPath(locale)}?series=${encodeURIComponent(slug)}`;
