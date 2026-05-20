# Homebrew TUI Design Specification

## Overview
A terminal user interface (TUI) for managing Homebrew packages, built in Go using the Charm ecosystem (Bubble Tea, Lip Gloss, Bubbles, Huh, Glamour). It supports full management (search, install, update, uninstall) with a tabbed interface, inline progress logs, and a Catppuccin Mocha color theme with gradient headers/footers.

## Core Features
1.  **Tabbed Navigation**: Switch between "Installed", "Search", and "Updates" views.
2.  **Package Management**: Install, uninstall, and update Homebrew formulae and casks.
3.  **Details View**: Render rich markdown package information (description, dependencies, caveats) using Glamour.
4.  **Inline Progress**: Display a spinner and streaming logs when executing long-running `brew` commands.
5.  **Confirmation Dialogs**: Use `huh` for confirming destructive actions (uninstalls).

## Architecture & State
-   **Paradigm**: The Elm architecture via Charm `bubbletea`.
-   **Concurrency Model**: Synchronous/Action-Driven. The UI remains responsive (spinner animates, logs scroll) during a command, but only one action (e.g., an install) can run at a time. Other inputs are locked until completion.
-   **CLI Interaction**: Execute native `brew` commands (e.g., `brew info --json=v2`, `brew search`, `brew upgrade`) using `exec.Command` wrapped in `tea.Cmd`. Outputs are parsed and dispatched as `tea.Msg`.

## UI Components
1.  **Layout**:
    *   **Header**: Tab bar showing current context. Styled with a Catppuccin Mocha gradient.
    *   **Main Body**: Split view (or single view depending on context).
        *   Left/Main: `bubbles/list` for packages.
        *   Right/Overlay: Package details (rendered by Glamour) or search input (`bubbles/textinput`).
    *   **Footer**: Help menu and log/status output. Styled with a Catppuccin Mocha gradient.
2.  **Components**:
    *   `list`: For displaying Installed, Search Results, and Outdated packages.
    *   `textinput`: For the Search tab query.
    *   `spinner` & `viewport`: For the execution log pane.
    *   `huh`: For confirmation forms (e.g., "Confirm Uninstall?").

## Visual Design
-   **Theme**: [Catppuccin Mocha](https://catppuccin.com/palette).
-   **Gradients**: Header and Footer components will utilize lipgloss gradient rendering using Mocha colors (e.g., transitioning across Mauve, Pink, Flamingo, etc.).
-   **Package Info**: Markdown styling optimized for dark terminal backgrounds.

## Error Handling
-   Command failures (e.g., a failed install) will stop the spinner, display the error in the log pane, and allow the user to dismiss the error and resume normal navigation.