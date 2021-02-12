# Test

## Usage

Deploy Rabbitmq
```bash
docker-compose up -d
```

Configure Rabbitmq
```bash
./rabbitmq-setup.sh
```

Publish to Rabbitmq
```bash
./rabbitmq-publish.sh
```

Build and run Stanza CLI
```bash
./stanza.sh
```

## Load Generation

The package in `load/` can be used to send messages to Rabbitmq. The source can
be tweaked to increase number of goroutines, or number of messages sent per goroutine.
