## Development

When starting the dev server, use background mode:

```
astro dev --background
```

Manage the background server with `astro dev stop`, `astro dev status`, and `astro dev logs`.

## Documentation

When editing Korean (`src/content/posts/ko`) or Japanese (`src/content/posts/ja`) MDX,
run `pnpm mdx:check` before finishing to catch bold delimiters that render as literal text.
Do not require spaces around valid bold text or inline code: `**성능 관점**에서` and adjacent inline code are valid.
Use `pnpm mdx:fix` to add only the spaces needed by CommonMark's delimiter rules.
Do not change code blocks, frontmatter, or link destinations to satisfy the check.

Full documentation: https://docs.astro.build

Consult these guides before working on related tasks:

- [Adding pages, dynamic routes, or middleware](https://docs.astro.build/en/guides/routing/)
- [Working with Astro components](https://docs.astro.build/en/basics/astro-components/)
- [Using React, Vue, Svelte, or other framework components](https://docs.astro.build/en/guides/framework-components/)
- [Adding or managing content](https://docs.astro.build/en/guides/content-collections/)
- [Adding styles or using Tailwind](https://docs.astro.build/en/guides/styling/)
- [Supporting multiple languages](https://docs.astro.build/en/guides/internationalization/)
