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
- [x] Plugin support for reading/writing the version to/from a file
  - [x] Maybe allow the cli itself to install & manage plugins? `changeset plugin install <plugin-name>`
- [x] Add support for publishing a changeset
- [x] Add support for parsing the current version from one of the supported project files
- [ ] Add support for creating and amending a `CHANGELOG.md` file
- [ ] Add a command to preview the `CHANGELOG.md` file prefix before publishing. `changeset preview`
- [ ] Add support for consuming changesets and updating the version in supported project files:
  - [x] Unsupported project files (`.changeset/version` file)
  - [ ] pyproject.toml
  - [ ] package.json
  - [ ] Cargo.toml
  - [ ] Go.mod
- [ ] Add support for auto-committing changesets (via `--autocommit` flag for `changeset add`)
- [ ] Add support for tagging releases in git (via `--tag` flag for `changeset add`)
- [ ] Add support for an additional number in the version (e.g. `1.2.3.4`). This is for projects which are an add-on to existing projects.
- [ ] Side-car repo for bot to manage releases via a pull request

## Plugins

### VersionedFile

This plugin is used to read/write the version to/from a plain file.

The file must be a plain text file with the following format:

```text
1.2.3
```

Example config file:

```json
{
  "plugin": {
    "name": "versionfile",
    "sha256": "beef1de60035053ad01eff83875999dc9918a65e1cffc006fca95c3bfbe55d70",
    "url": "https://github.com/alex-way/changesets-go-versionfile-plugin/releases/download/0.0.2/versionfile.wasm",
    "versionedFile": ".changeset/version"
  }
}
```

## Implementing your own plugin
