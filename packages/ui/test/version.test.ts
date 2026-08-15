import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import { createRequire } from 'node:module';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import test from 'node:test';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const binPath = path.join(__dirname, '..', 'bin', 'ask.mjs');
const require = createRequire(import.meta.url);
const { version } = require('../package.json');

test('--version prints the package version and exits 0', () => {
  const result = spawnSync(process.execPath, [binPath, '--version'], { encoding: 'utf8' });
  assert.equal(result.status, 0, `exit code should be 0, got ${result.status}\nstderr: ${result.stderr}`);
  assert.equal(result.stdout.trim(), version);
  assert.equal(result.stderr, '');
});

test('-v prints the package version and exits 0', () => {
  const result = spawnSync(process.execPath, [binPath, '-v'], { encoding: 'utf8' });
  assert.equal(result.status, 0, `exit code should be 0, got ${result.status}\nstderr: ${result.stderr}`);
  assert.equal(result.stdout.trim(), version);
  assert.equal(result.stderr, '');
});
