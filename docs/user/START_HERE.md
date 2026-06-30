# Start Here

Use this page as the docs entry point.

## Choose A Path

| Goal | Read | Run |
| --- | --- | --- |
| First local setup | `docs/user/QUICKSTART.md` | `cc-bash-guard setup` |
| Pick a starter policy posture | `docs/user/OPERATIONAL_TEMPLATES.md` | `cc-bash-guard init --list-profiles --verbose` |
| Let an agent help author policy | `docs/user/AGENTIC_POLICY_AUTHORING.md` | `cc-bash-guard suggest "git status"` |
| Understand a decision | `docs/user/EXPLAIN.md` | `cc-bash-guard explain --why-not allow "git status"` |
| Write or review rules | `docs/user/PERMISSION_SCHEMA.md` | `cc-bash-guard verify` |
| Inspect semantic support | `docs/user/SEMANTIC_SCHEMAS.md` | `cc-bash-guard semantic-schema docker --examples` |
| Debug a broken setup | `docs/user/TROUBLESHOOTING.md` | `cc-bash-guard doctor` |
| Review safety boundaries | `docs/user/THREAT_MODEL.md` | `cc-bash-guard verify --all-failures` |

## Quick Commands

```sh
cc-bash-guard setup
cc-bash-guard init --list-profiles
cc-bash-guard init --list-profiles --verbose
cc-bash-guard init --profile git-safe
cc-bash-guard verify
cc-bash-guard doctor
cc-bash-guard explain "git status"
cc-bash-guard help permission
cc-bash-guard help semantic
```

Other useful references:

- `docs/user/SEMANTIC_COVERAGE.md`
- `docs/user/EXAMPLES.md`
- `docs/user/AWS_GUIDELINES.md`
