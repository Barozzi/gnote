# Gnote

Notes CLI for people with the first initial G, or like, anyone else who might find this useful.

## Get started

You need to download Go
You need to close the repo
You need to build the project locally `go build -o gnote`
Run CLI `./gnote`

My config file at `~/.config/gnote/gnote.yaml`

```yaml
---
## Vault Path is where my Obsidian vault is located
vault_path: /Users/gb0218/vaults/work
## Subpaths are sub-folders of my obsidian vault
day_subpath: "00-dev-log" ## Day is where my daily notes go.
## PARA Method: Projects, Areas, Resources, Archives are organized via the PARA method of note taking
## https://fortelabs.com/blog/para/
projects_subpath: "01-projects"
areas_subpath: "02-areas"
resources_subpath: "03-resources"
archives_subpath: "04-archives"
```

## Background

### Command: gnote day

I have been using Obsidian through neovim. I make a new daily note each morning to track what I'm doing, what needs to be done next, etc. I store these notes in my Obsidian vault so they can be searched later.

### Command: gnote ticket

We use Jira at work and each time I pull a new ticket, I make a note folder to track my investigation, things I've done, things I'm going to do, etc. Doing this helps me when I get interrupted mid-feature and then come back to the ticket. When I have good notes, I find it easier to deal with having lots of unfinished tickets.

I organize my notes using the PARA method.

```
$ gnote
 ▗▄▄▖▗▖  ▗▖ ▗▄▖▗▄▄▄▖▗▄▄▄▖
▐▌   ▐▛▚▖▐▌▐▌ ▐▌ █  ▐▌
▐▌▝▜▌▐▌ ▝▜▌▐▌ ▐▌ █  ▐▛▀▀▘
▝▚▄▞▘▐▌  ▐▌▝▚▄▞▘ █  ▐▙▄▄▖

  GNote helps you manage your notes, dev logs, and projects
by providing commands to quickly create and organize files.
Use 'gnote [command] --help' for more information about a specific command.

Usage:
  gnote [command]

Available Commands:
  archive     Archive a project
  completion  Generate the autocompletion script for the specified shell
  day         Create a new DevLog for the current day.
  help        Help about any command
  search      Search all DevLogs for a string match
  ticket      Create a new unit of work folder with a standard set of files.

Flags:
  -h, --help     help for gnote
  -t, --toggle   Help message for toggle

Use "gnote [command] --help" for more information about a command.
```
