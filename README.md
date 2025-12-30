# HodlBook App Store

Umbrel Community App Store for [HodlBook](https://github.com/HodlBook/hodlbook) - a privacy-focused, self-hosted cryptocurrency portfolio tracker.

## Adding to Umbrel

To install the HodlBook App Store on your Umbrel:

1. Open your Umbrel dashboard
2. Go to **App Store** → **Community App Stores**
3. Add this repository URL:
   ```
   https://github.com/HodlBook-org/hodlbook-app-store
   ```
4. Find and install **HodlBook** from the app store

## Repository Structure

```
hodlbook-app-store/
├── umbrel-app-store.yml      # App store metadata
├── hodlbook/
│   ├── umbrel-app.yml        # App manifest
│   └── docker-compose.yml    # Docker configuration
├── placeholder/
│   ├── umbrel-app.yml        # Template with PLACEHOLDER_VERSION
│   └── docker-compose.yml    # Template with PLACEHOLDER_VERSION
├── scripts/
│   └── main.go               # sync-umbrel tool source
└── .github/workflows/
    └── sync.yml              # Auto-sync workflow
```

## Version Sync

This repo automatically syncs with HodlBook releases:

1. When HodlBook publishes a new release, it triggers a `repository_dispatch` event
2. The `sync.yml` workflow runs `sync-umbrel`
3. `sync-umbrel` fetches the latest version from HodlBook's `umbrel/umbrel-app.yml`
4. Templates in `placeholder/` are processed, replacing `PLACEHOLDER_VERSION` with the actual version
5. Updated files are committed and pushed

### Manual Sync

To sync manually:

```bash
# Build the tool
go build -o sync-umbrel ./scripts

# Run sync
./sync-umbrel
```

Or trigger the workflow manually from GitHub Actions.

## About HodlBook

HodlBook is a self-hosted cryptocurrency portfolio tracker that lets you:

- Track your crypto holdings privately
- Record deposits, withdrawals, and exchanges
- View portfolio allocation and performance metrics
- Get live price updates via Binance/CoinGecko APIs

All data stays on your machine - no external accounts or data sharing required.

**Main repository:** [github.com/HodlBook/hodlbook](https://github.com/HodlBook/hodlbook)

## License

MIT License - see [LICENSE](LICENSE) for details.
