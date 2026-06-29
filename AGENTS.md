# AGENTS.md

Deno 2 CLI project (not Node/npm). Runtime + dependency map live in `deno.json`; do not reach for `package.json`/`node_modules`.

## Commands

Defined as `deno` tasks in `deno.json`:

- `deno task dev` — run CLI with `--watch` (auto-restart on edit)
- `deno task cli` — run CLI once
- `deno task compile:cli` — compile standalone binary to `dist/bin/ordo`
- `deno task start:web` — run the Hono JSX web UI (`web-ui/main.tsx`)

No tests, lint, typecheck, or CI are configured. Don't assume `deno test`/`deno lint` pass or matter; `deno fmt` is the only formatter available but is not enforced.

Always run with `--allow-all` (as the tasks do) — the CLI spawns IDE processes and uses `Deno.build.os`. `start:web` deliberately uses only `--allow-net`.

## Architecture

- **Entry:** `cli/main.ts` wires Commander commands from `cli/command/{project,ide,help}`.
- **Storage:** `lib/config.ts` exports a singleton `Conf` instance (`projectName: "ordo"`). All projects/IDEs persist via `npm:conf` to the user's OS config dir (e.g. `~/Library/Preferences/ordo/...` on macOS) — there is no repo-local data file. Mutations go through `store/project.store.ts` and `store/ide.store.ts`, never edit `Conf` directly.
- **Types:** `types/store.ts` (`Project`, `IDE`, `ProjectAndIDE`) and `types/config.ts` (`Config`) are the source of truth for the persisted shape.
- **Opening projects:** `utils/command.ts` builds the launch string. On darwin/linux an `open` prefix is prepended only when the IDE path is already absolute; relative paths/commands are used verbatim. Windows gets no prefix.
- **Web UI:** Separate Hono app (`web-ui/main.tsx`). `deno.json` sets `jsx: precompile` + `jsxImportSource: hono/jsx` — do not add React/Vite.

## Conventions

- Style: 4-space indent, double quotes for non-npm imports, `.ts`/`.tsx` extensions kept in import specifiers (Deno requires explicit extensions).
- User-facing messaging goes through `consola` (logs + interactive `consola.prompt` for IDE selection when a project has no default). Don't swap in `console.log`.
- UUIDs (`npm:uuid` v4) are minted only in `store/*.store.ts` `add*` functions; `id` is the canonical key, but `findOneProject`/`deleteProject` accept either `id` or `name`.
- Don't introduce npm-only packages without mapping them in `deno.json` `imports` first (e.g. `"foo": "npm:foo@x"`). `jsr:` and `npm:` specifiers are used inline in source via the import map, not `npm:` install.