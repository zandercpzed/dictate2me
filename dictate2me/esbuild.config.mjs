import { build } from 'esbuild';

build({
  entryPoints: ['src/main.ts'],
  bundle: true,
  outfile: 'dist/main.js',
  format: 'cjs',
  platform: 'node',
  external: ['obsidian'],
  sourcemap: true,
  logLevel: 'info',
}).catch(() => process.exit(1));
