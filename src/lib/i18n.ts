export const locales = ['ko', 'en', 'ja'] as const;
export type Locale = (typeof locales)[number];
export type LanguageLinks = Partial<Record<Locale, string>>;

export const defaultLocale: Locale = 'ko';

export const ui = {
  ko: {
    posts: '글',
    homeTitle: '블로그',
    latestPosts: '최근 글',
    allPosts: '모든 글 보기',
    allProjects: '모든 프로젝트 보기',
    readPost: '글 읽기',
    backToPosts: '글 목록으로',
    skipToContent: '본문으로 건너뛰기',
    blogDescription: '소프트웨어, 도구, 그리고 개발하면서 배우는 것들에 대한 기록입니다.',
    promptPreset: '예시 질문 선택',
    promptQuestion: '선택한 질문',
    promptActions: '질문 복사 또는 AI에서 열기',
    openInAI: '{platform}에서 질문하기 (새 탭)',
    platformPromptCopied: '질문을 복사했어요. 입력칸이 비어 있으면 붙여넣어 주세요.',
    agentPrompts: [
      { page: 'home', label: '블로그 소개', message: '이 블로그를 읽고, 운영자가 어떤 사람인지 알려줘.' },
      { page: 'posts', label: '관심사', message: '글들을 살펴보고, 이 사람의 주요 관심사를 정리해 줘.' },
      { page: 'projects', label: '프로젝트', message: '프로젝트들을 살펴보고, 이 사람이 무엇을 만드는지 알려줘.' },
      { page: 'posts', label: '추천 글', message: '처음 방문한 사람에게 추천할 글 3개와 그 이유를 알려줘.' },
    ],
    copyPrompt: '복사',
    promptCopySuccess: '질문과 주소를 복사했어요.',
    promptCopyError: '복사하지 못했어요. 질문을 직접 복사해 주세요.',
    clearFilter: '필터 해제',
    feedTitle: 'gymynnym의 블로그',
    projects: '프로젝트',
    contact: '연락처',
    topics: '토픽',
    browseByTopic: '토픽별 글 보기',
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
    posts: 'posts',
    homeTitle: 'Blog',
    latestPosts: 'Recent posts',
    allPosts: 'View all posts',
    allProjects: 'View all projects',
    readPost: 'Read post',
    backToPosts: 'Back to posts',
    skipToContent: 'Skip to content',
    blogDescription: 'Notes on software, tools, and the things I learn while building.',
    promptPreset: 'Choose an example question',
    promptQuestion: 'Selected question',
    promptActions: 'Copy question or open in AI',
    openInAI: 'Ask {platform} (opens in a new tab)',
    platformPromptCopied: 'Question copied. If the prompt field is empty, paste it there.',
    agentPrompts: [
      { page: 'home', label: 'About the blog', message: 'Read this blog and tell me about the person behind it.' },
      { page: 'posts', label: 'Interests', message: 'Browse these posts and summarize the author’s main interests.' },
      { page: 'projects', label: 'Projects', message: 'Explore these projects and tell me what this person builds.' },
      {
        page: 'posts',
        label: 'Reading list',
        message: 'Recommend three posts for a first-time visitor and explain why.',
      },
    ],
    copyPrompt: 'Copy',
    promptCopySuccess: 'Question and URL copied.',
    promptCopyError: 'Could not copy. Copy the question manually.',
    clearFilter: 'clear',
    feedTitle: "gymynnym's blog",
    projects: 'projects',
    contact: 'contact',
    topics: 'topics',
    browseByTopic: 'Browse by topic',
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
    posts: '記事',
    homeTitle: 'ブログ',
    latestPosts: '最近の記事',
    allPosts: 'すべての記事',
    allProjects: 'すべてのプロジェクト',
    readPost: '記事を読む',
    backToPosts: '記事一覧へ',
    skipToContent: '本文へスキップ',
    blogDescription: 'ソフトウェアやツール、ものづくりを通して学んだことを記録しています。',
    promptPreset: '質問例を選ぶ',
    promptQuestion: '選択した質問',
    promptActions: '質問をコピー、またはAIで開く',
    openInAI: '{platform}に聞く（新しいタブ）',
    platformPromptCopied: '質問をコピーしました。入力欄が空の場合は貼り付けてください。',
    agentPrompts: [
      { page: 'home', label: 'ブログ紹介', message: 'このブログを読んで、運営者がどんな人なのか教えて。' },
      { page: 'posts', label: '関心のあるテーマ', message: '記事を読んで、この人が関心を持っているテーマをまとめて。' },
      { page: 'projects', label: 'プロジェクト', message: 'プロジェクトを見て、この人が何を作っているのか教えて。' },
      { page: 'posts', label: 'おすすめ記事', message: '初めて訪れた人におすすめの記事を3つ、理由とともに教えて。' },
    ],
    copyPrompt: 'コピー',
    promptCopySuccess: '質問とURLをコピーしました。',
    promptCopyError: 'コピーできませんでした。質問を直接コピーしてください。',
    clearFilter: '解除',
    feedTitle: 'gymynnymのブログ',
    projects: 'プロジェクト',
    contact: 'お問い合わせ',
    topics: 'トピック',
    browseByTopic: 'トピックから記事を探す',
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
    timeZone: 'Asia/Seoul',
  }).format(value);
}
