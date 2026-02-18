# setup-mdtest (GitHub Action)

Composite GitHub Action to install the **mdtest** CLI and add it to
`PATH`.

This action installs `mdtest` using `go install` and places the binary
in a deterministic location inside the workflow workspace:
`${GITHUB_WORKSPACE}/.tools/bin`.

------------------------------------------------------------------------

## Usage

### Basic (recommended: pin a version)

``` yaml
- name: Setup mdtest
  uses: UpCloudLtd/mdtest/.github/actions/setup-mdtest@main
  with:
    version: v0.1.0
```

### Use `latest` (not reproducible)

``` yaml
- name: Setup mdtest
  uses: UpCloudLtd/mdtest/.github/actions/setup-mdtest@main
```

### If your workflow does not installs Go

``` yaml
- name: Setup mdtest (no Go install)
  uses: UpCloudLtd/mdtest/.github/actions/setup-mdtest@main
  with:
    install-go: "true"
    go-version: "1.24.x"
```

------------------------------------------------------------------------

## Inputs

  ---------------------------------------------------------------------------
  Name           Required           Default          Description
  -------------- ------------------ ---------------- ------------------------
  `version`      No                 `latest`         mdtest version to
                                                     install (e.g. `v0.1.0`).
                                                     Prefer pinning a tag for
                                                     reproducible builds.

  `install-go`   No                 `true`           Whether the action
                                                     should install Go using
                                                     `actions/setup-go`.

  `go-version`   No                 `1.24.x`         Go version to install
                                                     when `install-go=true`.
  ---------------------------------------------------------------------------

------------------------------------------------------------------------

## Outputs

  -----------------------------------------------------------------------
  Name                   Description
  ---------------------- ------------------------------------------------
  `gobin`                Directory where `mdtest` was installed (defaults
                         to `${GITHUB_WORKSPACE}/.tools/bin`).

  -----------------------------------------------------------------------

------------------------------------------------------------------------

## Notes

-   Go installed by this action (when `install-go=true`) is available
    for the rest of the **job**.

-   Each job runs on a fresh runner. If you have multiple jobs, each
    must install Go if needed.

-   This action uses:

        go install github.com/UpCloudLtd/mdtest@<version>

-   Avoid using `latest` in production workflows. Pin a version tag
    instead.