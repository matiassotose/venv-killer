# Release Process

To release a new version of `venv-killer`, follow these steps:

## Prerequisites

1.  **Install Goreleaser**:
    ```bash
    brew install goreleaser/tap/goreleaser
    ```
    Or follow instructions at [goreleaser.com](https://goreleaser.com/install/).

2.  **GitHub Token**:
    Ensure you have a `GITHUB_TOKEN` with `repo` permissions exported in your environment.

## Creating a Release

1.  **Tag the release**:
    ```bash
    git tag -a v0.1.0 -m "First release"
    git push origin v0.1.0
    ```

2.  **Run Goreleaser**:
    ```bash
    goreleaser release --clean
    ```

## Testing Locally

To test the build process without publishing:

```bash
goreleaser release --snapshot --clean
```

Artifacts will be in the `dist/` folder.
