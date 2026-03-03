# nag

A LazyGit-style terminal UI for [Apple Reminders](https://support.apple.com/guide/reminders/welcome/mac). Browse lists, view reminders, create, complete, and delete — all from the terminal.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![macOS](https://img.shields.io/badge/macOS-only-000000?logo=apple&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue)

## Features

- **Two-panel layout** — lists sidebar + reminders
- **Smart lists** — Today (includes overdue) and Scheduled views
- **Vim-style navigation** — `j`/`k`, `g`/`G`, `Ctrl-d`/`Ctrl-u`
- **Create reminders** — title, due date with time, priority
- **Complete/uncomplete** — toggle with Space, 2s grace period to undo
- **Delete** — with confirmation prompt
- **Open in Reminders** — jump to the reminder in Apple Reminders
- **Show/hide completed** — toggle visibility with `c`
- **Auto-refresh** — polls every 10 seconds for external changes
- **Filter/search** — fuzzy search across titles and notes
- **Mouse support** — click to focus panels, scroll to navigate

## Install

Requires [Go 1.25+](https://go.dev/dl/) and **macOS** (uses native EventKit via cgo):

```bash
go install github.com/oronbz/nag@latest
```

Or build from source:

```bash
git clone https://github.com/oronbz/nag.git
cd nag
make build
```

## Usage

```bash
nag
```

On first run, macOS will prompt for Reminders access. You can manage this in **System Settings > Privacy & Security > Reminders**.

## Key Bindings

### Navigation

| Key | Action |
|-----|--------|
| `j` / `k`, `↑` / `↓` | Navigate |
| `Enter` | Select list |
| `Tab` / `Shift-Tab` | Switch panel |
| `g` / `G` | Jump to top / bottom |
| `Ctrl-d` / `Ctrl-u` | Page down / page up |
| `/` | Filter / search |
| Left click | Focus panel |
| Mouse wheel | Scroll |

### Actions

| Key | Action |
|-----|--------|
| `Space` / `x` | Toggle reminder complete |
| `n` | Create new reminder |
| `d` | Delete reminder |
| `o` | Open in Apple Reminders |
| `c` | Toggle show completed |
| `r` | Refresh |

### General

| Key | Action |
|-----|--------|
| `?` | Toggle help overlay |
| `Esc` | Close dialog / overlay |
| `q` / `Ctrl-C` | Quit |

## Layout

```
┌──────────┬─────────────────────┐
│  Lists   │     Reminders       │
│          │                     │
│ ◉ Today  │ ☐ Buy groceries !!! │
│ ▦ Sched. │   today  Get milk.. │
│ ───────  │ ☑ Call dentist      │
│ Personal │   yesterday         │
│ Work     │                     │
└──────────┴─────────────────────┘
 status bar
```

## Built With

- [go-eventkit](https://github.com/BRO3886/go-eventkit) — Native macOS EventKit bindings (cgo)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework (Elm Architecture)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Styling and layout
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components (list, viewport, textinput, spinner)

## License

MIT
