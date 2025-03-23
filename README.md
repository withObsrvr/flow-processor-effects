# Flow Effects Processor

This is a Flow processor plugin that transforms operations into detailed effect records, providing a comprehensive view of blockchain changes.

## Overview

The Effects Processor analyzes transaction operations and generates effect records that describe the changes made to the ledger. It adapts the Stellar effects processing logic to work within the Flow plugin architecture, providing detailed information about operations such as:

- Account creation, crediting, and debiting
- Trustline manipulation
- Offer creation and trading
- Data entry changes
- Claimable balance operations
- Liquidity pool interactions
- Contract operations

## Configuration

The plugin requires the following configuration parameters:

```json
{
  "network_passphrase": "Public Global Stellar Network ; September 2015"
}
```

| Parameter | Required | Description |
|-----------|----------|-------------|
| network_passphrase | Yes | The network passphrase used for cryptographic operations |

## Usage

### Building the Plugin

You can build the plugin using Nix or standard Go tools:

#### Using Nix (recommended for reproducible builds)

```bash
# Initialize development environment
nix develop

# Or build directly
nix build
```

The plugin will be available at `./result/lib/flow-processor-effects.so`.

#### Using Go directly

```bash
go build -buildmode=plugin -o flow-processor-effects.so .
```

### Plugin Integration

To use the plugin with Flow, add it to your Flow configuration:

```yaml
plugins:
  processors:
    - path: /path/to/flow-processor-effects.so
      config:
        network_passphrase: "Public Global Stellar Network ; September 2015"
```

## Input and Output

### Input

The plugin expects transaction data with the following fields:
- `envelope_xdr`: Base64-encoded transaction envelope XDR
- `result_xdr`: Base64-encoded transaction result XDR
- `meta_xdr`: Base64-encoded transaction meta XDR
- `ledger_sequence`: Ledger sequence number
- `ledger_close_time`: ISO8601 formatted ledger close time

### Output

For each operation in the transaction that generates effects, the plugin outputs effect records with the following structure:

```json
{
  "address": "GA...",
  "address_muxed": "MA...",
  "operation_id": 12345,
  "details": {
    "asset_type": "native",
    "amount": "100.0000000"
  },
  "type": 2,
  "type_string": "account_credited",
  "closed_at": "2023-03-23T12:34:56Z",
  "ledger_sequence": 42,
  "index": 0,
  "id": "12345-0"
}
```

## Development

To set up a development environment:

```bash
# Clone the repository
git clone https://github.com/withObsrvr/flow-processor-effects.git
cd flow-processor-effects

# Set up development environment with Nix
nix develop

# Or manually install dependencies
go mod tidy
go mod vendor
```

## License

This project is licensed under the [MIT License](LICENSE).
