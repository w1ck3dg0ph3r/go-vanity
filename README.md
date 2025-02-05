# go-vanity

Dead simple Go vanity URL server.

```bash
go run . -listen :5000 -vanities vanities.list
```

Example `vanities.list`:

```bash
# {module} {path} {ref}
github badgerdb hypermodeinc/badger main
# {module} {path} {ref}
gitlab gitlab-runner gitlab-org/gitlab-runner main
# {module} {host} {path} {ref}
forgejo forgejo-runner code.forgejo.org forgejo/runner main
```
