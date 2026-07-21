# User Documentation

`cc-bash-guard` is declarative, testable Bash permission policy for Claude
Code. It is policy-as-code for Bash permissions: semantic rules, examples as
tests, `verify`, `explain`, and verified artifacts for hook execution.

Start with `docs/user/START_HERE.md`. It points to the right next document for
first setup, policy authoring, command explanation, semantic schemas, examples,
and troubleshooting.

Common paths:

- `docs/user/START_HERE.md`
- `docs/user/QUICKSTART.md`
- `docs/user/OPERATIONAL_TEMPLATES.md`
- `docs/user/AGENTIC_POLICY_AUTHORING.md`
- `docs/user/EXPLAIN.md`
- `docs/user/PERMISSION_SCHEMA.md`
- `docs/user/THREAT_MODEL.md`
- `docs/user/SEMANTIC_SCHEMAS.md`
- `docs/user/SEMANTIC_COVERAGE.md`
- `docs/user/EXAMPLES.md`
- `docs/user/TROUBLESHOOTING.md`
- `docs/user/AWS_GUIDELINES.md`

Useful CLI help:

```sh
cc-bash-guard setup
cc-bash-guard help
cc-bash-guard help permission
cc-bash-guard help semantic
cc-bash-guard semantic-schema docker --examples
cc-bash-guard help examples
cc-bash-guard help troubleshoot
```
