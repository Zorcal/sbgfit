# sbgfit

A fitness application focused on CrossFit and Hyrox training, with bodybuilding support.

## Core Principles

**Speed is the primary UX factor.** All feature development and decisions prioritize user-perceived speed and simplicity:

- **Simple, fast UX/UI** - Minimize friction in user interactions
- **Clean, performant code** - Optimize for speed in implementation and execution
- **CrossFit & Hyrox focused** - Primary training methodologies with bodybuilding as secondary support

When deciding on features, speed and simplicity guide all development choices.

## Development Guidelines

All Go development must follow guidelines in [`backend/AGENTS.md`](./backend/AGENTS.md).

## Development Requirements

Add the following to your /etc/hosts for communication with app and docker containers to work correctly:

```sh
127.0.0.1	sbgfit-postgres
```

Run:

```sh
cd ./infra && docker compose up -d
```

Install `moq`:

```sh
cd ./backend && go install github.com/matryer/moq@latest
```

## Local Debugging

### Observability Stack

The application includes distributed tracing via OpenTelemetry. To access the observability tools:

1. **Start the observability stack**: `make up`
2. **Access Grafana**: Open http://localhost:3000
   - **Username**: `admin`
   - **Password**: `admin`
3. **View traces**: Navigate to "Explore" and select the "Tempo" datasource to search and view distributed traces
