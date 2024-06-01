# Changeset

This is honestly just a fun side project I'm working on with no intentions of using it in any production-like setting.

It's based on the existing [changesets](https://github.com/changesets/changesets) tool, but I wanted to make it more generic and not tied to a specific project.

## Usage

### Adding a changeset

```bash
changeset add --bump-type major --message "Added a new feature" # or simply `changeset add`
```

### Consuming changesets

```bash
changeset version
```

### Getting the current version

```bash
changeset get-version
```

A dry run can be performed by passing the `--dry-run` flag.

This will output the highest version type found in the `.changeset` directory and the changesets that were found.

## TODO

- [x] Add support for creating a new changeset
- [ ] Add support for publishing a changeset
- [ ] Add support for parsing the current version from one of the supported project files
- [ ] Add support for creating and amending a `CHANGELOG.md` file
- [ ] Add support for consuming changesets and updating the version in supported project files:
  - [ ] Unsupported project files (`.changeset/version` file)
  - [ ] pyproject.toml
  - [ ] package.json
  - [ ] Cargo.toml
  - [ ] Go.mod
- [ ] Add support for auto-committing changesets (via `--autocommit` flag for `changeset add`)
- [ ] Add support for tagging releases in git (via `--tag` flag for `changeset add`)
- [ ] Add support for an additional number in the version (e.g. `1.2.3.4`). This is for projects which are an add-on to existing projects.
