# 🧠 Onigma Shell

**Windows terminal reimagined** – a custom interactive shell written in Go with all the shortcuts, colors, and conveniences you expect, plus extras.

![License](https://img.shields.io/badge/license-MIT-blue)
![Platform](https://img.shields.io/badge/platform-Windows-0078d7)

Onigma Shell is a drop‑in replacement for `cmd.exe` that brings:

- 🎨 **Modern ANSI colored prompt** (works in Windows Terminal, VS Code, ConPTY)
- 🔁 **Command history** with persistent storage (↑/↓ arrows)
- 📁 **Tab completion** for commands, aliases, and file/directory names
- ⌨️ **Full line editing** – Home/End, Ctrl+A/E, Ctrl+U/K, Ctrl+W
- 🛡️ **Ctrl+C cancels** the current command (does **not** exit the shell)
- 📦 **Built‑in commands**: `cd`, `exit`, `cls`/`clear`, `neofetch`
- ⚡ **Aliases** – `ls` → `dir`, `cat` → `type`, `man` → `help`, easily extendable
- 💪 **Runs any Windows command** – Python, Node, Git, SSH, Docker, batch files, etc.
- 🪟 **Pure Windows native** – uses `cmd /C`, so pipes (`|`) and redirection (`>`) work perfectly

> **Yes, you can use Onigma Shell for your daily work.** Run Python scripts, `git status`, `npm install`, `ssh`, or even start the Python REPL – everything just works.

---

## ✨ Features at a glance

| Feature | How it works |
|---------|--------------|
| **Command history** | Up/Down arrows – saved to `.onigma_history` |
| **Tab completion** | Commands, aliases, and file/folder names |
| **Quoted paths** | `cd "My Documents"` works |
| **Aliases** | `ls` → `dir`, plus custom aliases you can add |
| **Built‑ins** | `cd`, `exit`, `cls`, `clear`, `neofetch` |
| **External commands** | `python`, `git`, `node`, `winget`, `ssh`, … |
| **Pipes & redirection** | `dir \| find ".go"`, `echo hello > file.txt` |
| **Ctrl+C** | Cancels current line – safe |
| **Ctrl+D / Ctrl+Z** | Exits the shell |
| **Home / End** | Jump to start/end of line |
| **Ctrl+U / Ctrl+K** | Delete from cursor to start / end |
| **Ctrl+W** | Delete previous word |
| **F7** (optional) | Interactive history menu (with readline) |

---

## 🚀 Getting started

### Prerequisites
- **Go** 1.18+ (only for building from source)
- **Windows** 10/11 (or Windows Server 2019+ with ConPTY support)

### Option 1: Download pre‑built binary
[Download the latest `onigma.exe` from Releases](https://github.com/AlexanderXinarxZenDev/Onigma-Shell/releases/) – just run it.

### Option 2: Build from source

```bash
git clone https://github.com/OnigmaShell/onigma
cd onigma
go mod init onigma              # only if no go.mod exists
go get github.com/chzyer/readline
go build -o onigma.exe main.go
```

Then run:  
```bash
.\onigma.exe
```

### Option 3: Run without building (go run)

```bash
go run main.go
```

> *Note:* The readline library is required for history and line editing. If you prefer zero dependencies, use the `bufio` fallback (remove readline imports) – but you'll lose arrow keys and completion.

---

## 🎮 Usage examples

```text
   ██████╗ ███╗   ██╗██╗ ██████╗ ███╗   ███╗ █████╗
  ██╔═══██╗████╗  ██║██║██╔════╝ ████╗ ████║██╔══██╗
  ██║   ██║██╔██╗ ██║██║██║  ███╗██╔████╔██║███████║
  ██║   ██║██║╚██╗██║██║██║   ██║██║╚██╔╝██║██╔══██║
  ╚██████╔╝██║ ╚████║██║╚██████╔╝██║ ╚═╝ ██║██║  ██║
   ╚═════╝ ╚═╝  ╚═══╝╚═╝ ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝

          Onigma Shell – Now with Windows Shortcuts!

➜ ONIGMA (C:\Users\John) > cd "Desktop\My Project"
➜ ONIGMA (~\Desktop\My Project) > python --version
Python 3.12.0

➜ ONIGMA (~\Desktop\My Project) > ls
onigma.exe   main.go   README.md   downloads\

➜ ONIGMA (~\Desktop\My Project) > git status
On branch main
Your branch is up to date with 'origin/main'.

➜ ONIGMA (~\Desktop\My Project) > neofetch
┌─[ John@DESKTOP-ABC ]
└─[ ~\Desktop\My Project ]
  ➜  Shell: Onigma Shell (windows)

➜ ONIGMA (~\Desktop\My Project) > echo "Hello, pipes!" | find "pipes"
"Hello, pipes!"

➜ ONIGMA (~\Desktop\My Project) > exit
```

---

## ⌨️ Shortcuts & key bindings

| Key                     | Action |
|-------------------------|--------|
| **↑ / ↓**               | History navigation |
| **← / →**               | Move cursor left/right |
| **Tab**                 | Auto‑complete command, alias, file/folder |
| **Home / End**          | Jump to start/end of line |
| **Ctrl + A / Ctrl + E** | Same as Home/End |
| **Ctrl + U**            | Clear from cursor to beginning of line |
| **Ctrl + K**            | Clear from cursor to end of line |
| **Ctrl + W**            | Delete previous word |
| **Ctrl + L**            | Clear screen (same as `cls`) |
| **Ctrl + C**            | Cancel current line – **does not exit** |
| **Ctrl + D / Ctrl + Z** | Exit shell |
| **F7**                  | Interactive history menu (if terminal supports it) |

---

## 🛠️ Built‑in commands

| Command       | Description |
|---------------|-------------|
| `cd [path]`   | Change directory – supports spaces and quotes. `cd` alone goes to home. |
| `exit`        | Terminates the shell. |
| `cls` / `clear` | Clears the screen. |
| `neofetch`    | Display system info: user, host, current dir, shell name, OS. |

---

## 🔧 Customization

### Adding new aliases
Edit the `aliases` map in `main.go`:

```go
var aliases = map[string]string{
    "ls":     "dir",
    "cat":    "type",
    "man":    "help",
    "gst":    "git status",   // add your own
    "cl":     "clear",
}
```

### Adding more built‑ins
Add a case to the `executeCommand` switch:

```go
case "mycmd":
    // your logic here
    return nil
```

### Changing prompt colors
Modify the `green`, `cyan`, `blue` constants at the top of the file.  
Use standard ANSI codes: `\033[31m` (red), `\033[33m` (yellow), etc.

---

## 📦 Dependencies

- [github.com/chzyer/readline](https://github.com/chzyer/readline) – for advanced line editing, history, and completion.  
- The standard library (`os`, `os/exec`, `strings`, etc.) – no other external dependencies.

---

## 🤝 Contributing

Pull requests are welcome! Feel free to open an issue for:

- New built‑in commands (like `pwd`, `echo`, `mkdir`)
- Persistent aliases (load from a config file)
- Git branch in prompt
- Better environment variable support (PowerShell style)
- Background jobs (`&`)

---

## 📄 License

This project is licensed under the **MIT License** – see the [LICENSE](LICENSE) file for details.

---

## ❤️ Acknowledgements

- [chzyer/readline](https://github.com/chzyer/readline) for making the shell actually usable.
- The Go standard library for making this possible in a few hundred lines.

---

**Built with ❤️ for Windows users who want a terminal that feels like home.**
