# CI integration

`provenance` is a platform-agnostic binary — this is the documented pattern for wiring it into a CI pipeline as a required check, not a published GitHub Action.

## GitHub Actions

```yaml
name: provenance
on:
  pull_request:
    branches: [main]

jobs:
  provenance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0   # full history is required to diff against the base ref

      - name: Build provenance
        run: |
          git clone https://github.com/tolvi-labs/provenance /tmp/provenance
          cd /tmp/provenance && go build -o /usr/local/bin/provenance ./cmd/provenance

      - name: Run the capture gate
        id: check
        run: |
          provenance check \
            --base "origin/${{ github.base_ref }}" \
            --head "${{ github.sha }}" \
            --json-out report.json
        continue-on-error: true

      - name: Post the report to the PR
        if: always()
        run: |
          if [ -f report.json ]; then
            gh pr comment "${{ github.event.pull_request.number }}" --body-file report.json
          else
            gh pr comment "${{ github.event.pull_request.number }}" \
              --body "Provenance check produced no report — a malformed capture or an execution error stopped it before the report was written. See the job log."
          fi
        env:
          GH_TOKEN: ${{ github.token }}

      - name: Fail the check if blocked
        if: steps.check.outcome == 'failure'
        run: exit 1
```

Wire this job as a required status check in branch protection — that's what makes it un-skippable. The local pre-push hook (`provenance hook install`) is optional fast pre-flight; it is not the enforcement point.

## Other CI providers

The same three steps apply to any provider: check out with full history, run `provenance check --base <base> --head <head> --json-out report.json`, post `report.json` however that provider's PR/MR-comment mechanism works, and make the job a required check. Guard the post step on the file existing: a malformed capture or an execution error exits non-zero before `--json-out` is written, so an unguarded post fails with a confusing missing-file error exactly where the log matters most.
