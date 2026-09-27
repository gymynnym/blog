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
];
