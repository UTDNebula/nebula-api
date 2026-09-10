Moving stuff off of README.md

### How to Contribute

Create your own fork by [forking this repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo#forking-a-repository)

[Clone](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo#cloning-your-forked-repository) your forked repository. (Don't forget to install Git if you haven't already)

Submit proposed changes via a [Pull Request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request)

## Building
### Requirements
- [Golang 1.26 or Higher](https://go.dev/dl/)

### Building for Windows

Setup Go Dependencies with:
```cmd
.\build.bat setup
```

Build with:
```cmd
.\build.bat build
```

This will create an executable named `rest-api.exe` in the root directory.

Run with:
```cmd
.\rest-api.exe
```
> Note: Some users have experienced issues with Windows Defender or other antivirus blocking `rest-api.exe` from reading files, editing files, or causing slowed performance. Consider adding an exception to your `nebula-api` folder.

### Building for macOS, Linux, and WSL

Setup Go dependencies with:
```bash
make setup
```

Build with:
```bash
make build
```

This will create an executable named `rest-api` in the root directory.

> Note: If Make fails with "swag: No such file or directory" or similar, you may need to add GOPATH/bin to your path. On Mac/Linux, use `echo 'export PATH=${PATH}:'$(go env GOPATH)'/bin' >> ~/.zshrc && source ~/.zshrc` (or `.bashrc`) to add it permanently.

Run with:
```bash
./rest-api
```

## Running API locally
Copy `.env.template` to `.env` with:
```bash
cp .env.template .env
```

Enter Nebula MongoDB URI in `.env` (ask for help in the [Discord](https://discord.utdnebula.com))

Run `rest-api`:
```bash
./rest-api
```

Check command output to see the route serving traffic. It's likely port 8080.

Visit `http://localhost:8080` to access nebula-api locally.

> Storage and email routes require additional environment variables. If you're working on these routes, ask for help in the [Discord](https://discord.utdnebula.com)
