# Telegram DNS Resolver Bot

A simple Telegram bot written in Go that allows users to resolve domain names using different DNS resolvers. It supports commands for looking up domain names, listing available DNS resolvers, and performing DNS lookups with specific resolvers. The bot also includes a feature for IP address lookup (`/lookup`).

## Features
- Lookup domain IP addresses using different DNS resolvers.
- Support for multiple resolvers like Google, Cloudflare, Quad9, OpenDNS, and others.
- Inline query support for quick DNS lookups directly in Telegram chats.
- `/resolver` command to list available resolvers.
- `/lookup` command to look up domain names.
- `/dig` command for IP address lookups.
- Supports `/start` and `/help` for user guidance.

## Prerequisites

- [Go](https://golang.org/doc/install) (1.20 or higher)
- [Telegram Bot Token](https://core.telegram.org/bots#botfather)
- Optional: [Docker](https://www.docker.com/get-started) for containerization

## Setup

### 1. Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/XigmaDev/IPSeekBot.git
cd IPSeekBot
```

### 2. Install Dependencies
If you’re using Go modules, install the dependencies with:

```bash
Copy code
go mod tidy
```
### 3. Set Up Environment Variables
Create a .env file in the project root directory and add your Telegram bot token:

```bash
BOT_TOKEN=your-telegram-bot-token
```
You can get your bot token from BotFather.

### 4. Run the Bot Locally
To run the bot locally, use the following command:

```bash
Copy code
go run .
```
This will start the bot, and it will listen for commands on Telegram.

### 5. Docker Setup (Optional)




### 6. Available Commands
- /start: Start interacting with the bot.
- /help: Get a list of available commands.
- /resolver: List available DNS resolvers.
- /lookup [resolver] [domain]: Lookup the IP address for a domain using a specified resolver. Example: /lookup Google example.com.

### 7. Dockerfile
If you prefer using Docker, the project comes with a Dockerfile that builds the bot and allows you to run it in a container. See the Docker Setup section for instructions.

## Development
If you want to contribute to the project or make changes, follow these steps:

### 1. Fork the Repository
Click the "Fork" button at the top-right corner of this repository to create your own copy.

### 2. Clone Your Fork
Clone your fork to your local machine:

```bash
Copy code
git clone https://github.com/XigmaDev/IPSeekBot.git
cd IPSeekBot
```
### 3. Make Changes
Make your changes in a new branch:

```bash
Copy code
git checkout -b feature-branch
```
### 4. Commit and Push
Commit your changes:

```bash
Copy code
git add .
git commit -m "Add new feature"
```
Push your branch:

```bash
Copy code
git push origin feature-branch
```
### 5. Create a Pull Request
Go to the GitHub repository and create a pull request for your changes.

## License
This project is open-source and available under the MIT License.

## Acknowledgments
- Go for being an awesome programming language.
- Telegram Bot API for providing an easy way to interact with Telegram users.
- golangci-lint for linting the Go codebase.