# Examples

These files are copyable permission policies and verification helpers. They are
intentionally small enough to review, and many tool-specific examples overlap
with the broader operational templates.

Use the broader templates when you want a starting posture. Use the focused
examples when you want to copy one tool's rule shape into an existing policy.

## Operational templates

- [`personal-cautious.yml`](personal-cautious.yml): individual default posture
  for local Git inspection, Git writes, and narrow local inspection patterns.
- [`infra-cautious.yml`](infra-cautious.yml): combined AWS, Kubernetes, Helm,
  Helmfile, Terraform, Docker, and Argo CD posture.
- [`team-baseline.yml`](team-baseline.yml): small-team baseline combining Git
  and Kubernetes rules with end-to-end tests.
- [`ci-verify-policy.sh`](ci-verify-policy.sh): minimal verification script for
  CI or pre-commit style checks.

## Focused examples

These examples demonstrate a single parser or narrow workflow. Some are
deliberately redundant with `infra-cautious.yml` so each tool can be copied and
reviewed in isolation.

- [`git-status-semantic.yml`](git-status-semantic.yml): minimal Git semantic
  matching, including wrapped command forms.
- [`git-safe-readonly.yml`](git-safe-readonly.yml): Git read-only inspection
  with force push denied.
- [`aws-identity.yml`](aws-identity.yml): AWS caller identity with profile and
  environment matching.
- [`kubectl-readonly.yml`](kubectl-readonly.yml): Kubernetes read-only verbs and
  mutating verbs behind confirmation.
- [`gws-readonly.yml`](gws-readonly.yml): Google Workspace read-only operations
  with destructive and unmasked credential export operations denied.
- [`argocd-app-delete-deny.yml`](argocd-app-delete-deny.yml): Argo CD app get,
  sync, and delete posture.
- [`helm-readonly-upgrade.yml`](helm-readonly-upgrade.yml): Helm read-only,
  install/upgrade, plugin changes, and uninstall posture.
- [`helmfile-diff-apply.yml`](helmfile-diff-apply.yml): Helmfile diff allowed
  and apply behind confirmation.
- [`terraform-readonly.yml`](terraform-readonly.yml): Terraform read-only
  commands, apply confirmation, and auto-approved destroy denial.
- [`terraform-strict.yml`](terraform-strict.yml): stricter Terraform posture
  where plan, apply, state mutation, and workspace delete require confirmation.
- [`docker-readonly.yml`](docker-readonly.yml): Docker inspection commands.
- [`docker-strict.yml`](docker-strict.yml): Docker run/compose confirmation and
  high-risk Docker operations denied.
- [`docker-compose-safe.yml`](docker-compose-safe.yml): Docker Compose
  read-only operations scoped to one compose file and project name.

## Notes on overlap

- `infra-cautious.yml` intentionally includes reduced versions of the
  Kubernetes, Helm, Helmfile, Terraform, Docker, AWS, and Argo CD examples.
- `personal-cautious.yml`, `team-baseline.yml`, and `git-safe-readonly.yml`
  all include Git rules, but target different adoption scopes.
- Tool-specific files are kept separate so users can copy one reviewed section
  without adopting the full infrastructure template.
