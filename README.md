
# Nebula API

_A database for some really useful UTD data collected by our [tools](https://github.com/UTDNebula/api-tools)._

Project maintained by [Nebula Labs](https://about.utdnebula.com).

## Contributing

Contributions are welcome!

This project uses the MIT License.

Please visit our [Discord](https://discord.utdnebula.com) and talk to us if you'd like to contribute!

## Documentation

Documentation for the current production API can be found [here.](https://api.utdnebula.com/swagger/index.html)

## How to use

- Visit our [Discord](https://discord.utdnebula.com) and ask to be provisioned an API key (please provide details on your use case)
- Read the documentation listed above (and authenticate with your key for interactive demos)
- Make requests to `https://api.utdnebula.com` with your provisioned api key set as the `x-api-key` request header
- **Build cool stuff!**

## Contributing
Contributions are welcome!

This project uses the MIT License.

Please visit our [Discord](https://discord.utdnebula.com) and talk to us if you'd like to contribute!
### How to Contribute

Create your own fork by [forking this repository](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo#forking-a-repository)

[Clone](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo#cloning-your-forked-repository) your forked repository. (Don't forget to install Git if you haven't already)

Submit proposed changes via a [Pull Request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request)

## Building
### Requirements
- [Golang 1.23 or Higher](https://go.dev/dl/)


### Building for Windows
cd into `nebula-api\api`

Setup Go Dependencies with 
`.\build.bat setup`

Build with
`.\build.bat build`

Run with
`.\go-api.exe`

### Building for macOs, Linux, and WSL
cd into `nebula-api/api`

Setup Go dependencies with 
`make setup`

Build with
`make build`

Run with
`./go-api.exe`
## Running to API locally
Copy `.env.template` to `.env` with
`cp .env.template .env`

Enter Nebula MongoDB URI in `.env`

Run with
`./go-api.exe`

Check command output to see the route serving traffic. It's likely port 8080

Visit `http://localhost:8080` to access nebula-api locally
