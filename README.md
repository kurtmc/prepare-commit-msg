prepare-commit-msg
==================

prepare-commit-msg hook for git, this will prefix your commits with the branch
name if it matches the pattern:

```
[A-Z]+-[0-9]+
```

## Building

```
make
```

## .gitconfig

If you put the binary here: `~/.git_template/hooks/prepare-commit-msg` and have
this in your `.gitconfig`:

```
[init]
	templatedir = ~/.git_template
```

then when you create new repos or re-init an existing repo it will add this
hook.
