# Getting Started with Nebula API

Here's how to set up your local development environment.

## Prerequisites

Ensure you have the following installed:

- Git - For version control
  - If you've never used git, need a refresher, or need help setting it up, please check out [Nebula's Git Workshop](https://github.com/UTDNebula/git-workshop).
- Go - For building and running
  - You can install from the [Go website](https://go.dev/dl/), or from [Homebrew](https://brew.sh/) or another package manager for automatic updates
- Make _(Linux/macOS)_ - For build automation
  - Pre-installed on macOS (via `xcode-select --install`) and most Linux distributions (`build-essential`)
- Docker _(Optional)_ - For running containerized runners locally

**Make** is a build automation tool. Feel free to check out [`Makefile`](../Makefile) to see exactly what's being run.

If you're using Windows, instead of using `make`, you can also use our [`build.bat`](../build.bat) file. When you see any command starting with `make`, you can instead use `.\build.bat`.
For example instead of `make setup`, you can run `.\build.bat setup`.

## Local Setup

### Clone the Repository

Clone the repository with `git clone`, and then, `cd` into the project directory or open it in your preferred code editor.

**HTTPS:**

```bash
git clone https://github.com/UTDNebula/nebula-api.git
cd nebula-api
```

or

**SSH:**

```bash
git clone git@github.com:UTDNebula/nebula-api.git
cd nebula-api
```

Now open the cloned project in your code editor!

### Install Development Tooling

Install the Go static analysis and formatting tools (`staticcheck` and `goimports`):

```bash
make setup
```

or

```cmd
.\build.bat setup
```

If you installed **Go** with [**Homebrew**](https://brew.sh/), you need to add Go tools to your path for tools like `staticcheck` and others to work.

- **For zsh** (the default shell on MacOS)

  ```bash
  echo 'export PATH=${PATH}:`go env GOPATH`/bin' >> ~/.zshrc && source ~/.zshrc
  ```

- **For bash** (the default shell on most Linux distributions)

  ```bash
  echo 'export PATH=${PATH}:`go env GOPATH`/bin' >> ~/.bashrc && source ~/.bashrc
  ```

- **For fish**

  ```bash
  echo 'fish_add_path (go env GOPATH)/bin' >> ~/.config/fish/config.fish
  ```

### Configure Environment Variables

First, make a file called `.env` at the root of the project. 

Then, copy the contents of `.env.template` into it. Some parts in `nebula-api` require certain environment variables, which you can fill in `.env`.

If you're not sure what to put, please try contacting one of our members through [Discord](https://discord.utdnebula.com).

> [!IMPORTANT] Do NOT put any environment variables into `.env.template`.

### Run Code Verification & Formatting

Check your code with:

```bash
make check
```

or

```cmd
.\build.bat check
```

You'll want to run this frequently while developing.

### Build the CLI Executable

Compile the Go source code into a runnable binary (`./rest-api` on Linux/macOS, `rest-api.exe` on Windows):

```bash
make build
```

or

```cmd
.\build.bat build
```

### Verify with Automated Tests

Run the test suite to confirm your environment is ready:

```bash
make test
```

or

```cmd
.\build.bat test
```

If all tests pass, congratulations! Your local development environment is ready.

### Running the API locally

To run `rest-api` locally, use:

```bash
./rest-api
```

Check command output to see the route serving traffic. It's likely port 8080.

Visit `http://localhost:8080` to access nebula-api locally.

> [!NOTE] Some users have experienced issues with Windows Defender or other antivirus blocking `rest-api.exe` from reading files, editing files, or causing slowed performance. Consider adding an exception to your `nebula-api` folder.

## Next Step

Now that your environment is running, it's time to dive deeper. Check out [Project-Architecture.md](Project-Architecture.md)