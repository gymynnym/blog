import type { Locale } from '@/lib/i18n';

export interface Project {
  name: string;
  type: string;
  archived: boolean;
  description: Record<Locale, string>;
  href: string;
  topics: Array<{
    slug: string;
    name: string;
  }>;
}

export const projects: Project[] = [
  {
    name: 'loghub-me',
    type: 'webservice',
    archived: true,
    description: {
      ko: 'LogHub는 개발자들이 지식을 공유하고, 서로의 경험을 나누는 플랫폼입니다.',
      en: 'LogHub is a platform for developers to share knowledge and experiences with each other.',
      ja: 'LogHubは、開発者が知識を共有し、互いの経験を分かち合うためのプラットフォームです。',
    },
    href: 'https://github.com/loghub-me',
    topics: [
      { slug: 'typescript', name: 'TypeScript' },
      { slug: 'nextjs', name: 'Next.js' },
      { slug: 'kotlin', name: 'Kotlin' },
      { slug: 'spring-boot', name: 'Spring Boot' },
      { slug: 'elysia', name: 'ElysiaJS' },
    ],
  },
  {
    name: 'vscode-gitui',
    type: 'vscode extension',
    archived: false,
    description: {
      ko: 'VS Code 통합 터미널에서 GitUI 또는 Lazygit을 실행하는 확장 프로그램입니다.',
      en: 'A VS Code extension that opens GitUI or Lazygit in the integrated terminal.',
      ja: 'VS Code の統合ターミナルで GitUI または Lazygit を開く拡張機能です。',
    },
    href: 'https://github.com/gymynnym/vscode-gitui',
    topics: [
      { slug: 'typescript', name: 'TypeScript' },
      { slug: 'vscode', name: 'VS Code' },
    ],
  },
];
