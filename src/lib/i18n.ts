export const locales = ['ko', 'en', 'ja'] as const;
export type Locale = (typeof locales)[number];
export type LanguageLinks = Partial<Record<Locale, string>>;

export const defaultLocale: Locale = 'ko';

export const ui = {
  ko: {
    blogDescription: '소프트웨어와 도구, 그리고 무언가를 만들며 배운 것들을 기록합니다.',
    clearFilter: '필터 해제',
    feedTitle: 'gymynnym의 블로그',
    projects: '프로젝트',
    contact: '연락처',
    topics: '토픽',
    series: '시리즈',
    publishedAt: '게시일',
    updatedAt: '수정일',
    noPosts: '아직 작성된 포스트가 없습니다.',
    noFilterResults: '선택한 조건에 해당하는 항목이 없습니다.',
    previousPost: '이전 글',
    nextPost: '다음 글',
    languages: '언어',
    contents: '목차',
    writtenBy: '작성자:',
    comments: '댓글',
  },
  en: {
    blogDescription: 'Notes on software, tools, and the things I learn while building.',
    clearFilter: 'clear',
    feedTitle: "gymynnym's blog",
    projects: 'projects',
    contact: 'contact',
    topics: 'topics',
    series: 'series',
    publishedAt: 'published',
    updatedAt: 'updated',
    noPosts: 'No posts have been published yet.',
    noFilterResults: 'No items match the selected filters.',
    previousPost: 'previous',
    nextPost: 'next',
    languages: 'languages',
    contents: 'contents',
    writtenBy: 'Written by:',
    comments: 'comments',
  },
  ja: {
    blogDescription: 'ソフトウェアやツール、ものづくりを通して学んだことを記録しています。',
    clearFilter: '解除',
    feedTitle: 'gymynnymのブログ',
    projects: 'プロジェクト',
    contact: 'お問い合わせ',
    topics: 'トピック',
    series: 'シリーズ',
    publishedAt: '公開日',
    updatedAt: '更新日',
    noPosts: 'まだ投稿がありません。',
    noFilterResults: '選択した条件に該当する項目はありません。',
    previousPost: '前の記事',
    nextPost: '次の記事',
    languages: '言語',
    contents: '目次',
    writtenBy: '著者:',
    comments: 'コメント',
  },
} as const;

const dateLocales: Record<Locale, string> = {
  ko: 'ko-KR',
  en: 'en-US',
  ja: 'ja-JP',
};

export function isLocale(value: string | undefined): value is Locale {
  return locales.includes(value as Locale);
}

export function localeFromPath(pathname: string): Locale {
  const locale = pathname.split('/')[1];
  return isLocale(locale) ? locale : defaultLocale;
}

export function formatDate(value: Date, locale: Locale): string {
  return new Intl.DateTimeFormat(dateLocales[locale], {
    dateStyle: 'long',
  }).format(value);
}
