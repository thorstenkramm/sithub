import pluginVue from 'eslint-plugin-vue'
import tseslint from 'typescript-eslint'
import vueParser from 'vue-eslint-parser'

// typescript-eslint's recommended configs, with its TypeScript-only file globs extended to .vue files.
const tsRecommended = tseslint.configs.recommended.map((config) =>
  config.files ? { ...config, files: [...config.files, '**/*.vue'] } : config,
)

export default [
  {
    name: 'app/files-to-lint',
    files: ['**/*.{ts,mts,tsx,vue}'],
  },
  {
    name: 'app/files-to-ignore',
    ignores: ['**/dist/**', '**/dist-ssr/**', '**/coverage/**'],
  },
  ...pluginVue.configs['flat/essential'],
  ...tsRecommended,
  // tseslint sets its own parser globally, so .vue files need vue-eslint-parser restored afterwards.
  ...pluginVue.configs['flat/base'],
  {
    name: 'app/vue-typescript-setup',
    files: ['**/*.vue'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: { js: 'espree', jsx: 'espree', ts: tseslint.parser, tsx: tseslint.parser },
        // vue-eslint-parser bundles espree 9, which supports ES2024 at most.
        ecmaVersion: 2024,
        ecmaFeatures: { jsx: false },
        extraFileExtensions: ['.vue'],
      },
    },
    rules: {
      'vue/block-lang': ['error', { script: { lang: ['ts'], allowNoLang: false } }],
    },
  },
  {
    name: 'app/cypress-overrides',
    files: ['cypress/**/*.ts'],
    rules: {
      '@typescript-eslint/no-namespace': 'off',
      '@typescript-eslint/no-unused-expressions': 'off',
    },
  },
  {
    name: 'app/vue-rules',
    files: ['**/*.vue'],
    rules: {
      'vue/multi-word-component-names': [
        'error',
        { ignores: ['Alert', 'Breadcrumb', 'Button', 'Card', 'Checkbox', 'Table', 'Tooltip'] },
      ],
    },
  },
]
