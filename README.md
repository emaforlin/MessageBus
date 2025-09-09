# MessageBus

A lightweight message bus system with CLI client for publishing and subscribing to topics.

## Installation

### Using Go Install

If you have Go installed:

```bash
go install github.com/emaforlin/messagebus/cmd/cli@latest
```

### Binary Releases

Download the latest binary for your platform from the [releases page](https://github.com/emaforlin/messagebus/releases).

#### Linux only (for now)

```bash
# Download server
curl -L -o messagebus-server https://github.com/emaforlin/MessageBus/releases/download/latest/mbus-server
# Make executable
chmod +x messagebus-server
# Move to a directory in your PATH
sudo mv messagebus-server /usr/local/bin/
```

```bash
# Download client
curl -L -o mbus-cli https://github.com/emaforlin/MessageBus/releases/download/latest/mbus-cli
# Make executable
chmod +x mbus-cli
# Move to a directory in your PATH
sudo mv mbus-cli /usr/local/bin/
```

## CLI Usage

### Subscribe to a Topic

Subscribe to messages on a specific topic:

```bash
# Subscribe to a single topic
messagebus subscribe mytopic

# Subscribe with custom server address
messagebus subscribe mytopic --server localhost:50051
```

### Publish Messages

Publish messages to a topic:

```bash
# Publish a single message
messagebus publish mytopic "Hello, World!"

# Publish with custom server address
messagebus publish mytopic "Hello!" --server localhost:50051

# Publish from standard input
echo "My message" | messagebus publish mytopic

```

### Examples

#### Basic Publishing

```bash
# Publish a simple message
messagebus publish notifications "New user registered"
```

#### Using with Pipes

```bash
# Pipe a file content as a message
cat data.json | messagebus publish data-updates

# Process log files
tail -f app.log | messagebus publish logs
```

### Server Configuration

By default, the CLI connects to `localhost:50051`. You can specify a different server:

### Global Flags

- `--server`: Server address (default: "localhost:50051")
- `--help`: Show help information
- `--version`: Show version information

## Building from Source

```bash
# Clone the repository
git clone https://github.com/emaforlin/messagebus.git

# Navigate to the project directory
cd messagebus

# Build the CLI client
go build -o messagebus cmd/cli/main.go
```
