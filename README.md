## What?

> [!NOTE]
>
> Encouraged: **use the original [virtual-code-owners](https://github.com/sverweij/virtual-code-owners)**
> which was written in TypeScript and is still actively maintained.
>
> This is a re-implementation of [virtual-code-owners](https://github.com/sverweij/virtual-code-owners) 
> in go. I made this because I wanted to learn go and I thought this would be a 
> good project to start with. If you need the functionality of the original
> project - please use that one. If you'd rather have it as this go implementation,
> (and simply plonk an executable on a box that works) _please let me know_
> 
> Because this is a re-implementation and I've learned a bit more about the
> functionality there _are_ differences between the two, most notably that
> this one reads the virtual teams from json only.

This generates your `CODEOWNERS` file (_patterns_ x _users_) from

- a `VIRTUAL-CODEOWNERS.txt` (_patterns_ x _teams_)
- a `virtual-teams.json` (_teams_ x _users_)

... which makes it easier to keep `CODEOWNERS` in sync on multi-team mono repos
when you don't have (enough) 'real' GitHub or GitLab teams.

## Usage

- Rename `.github/CODEOWNERS` to `.github/VIRTUAL-CODEOWNERS.txt` and put team names in them.
- Define teams that don't (yet) exist on GitHub level in `.github/virtual-teams.json`
- Run this:

```
vcodeowners
# Wrote '.github/CODEOWNERS'
```

- :sparkles:

## Formats

### VIRTUAL-CODEOWNERS.txt

`VIRTUAL-CODEOWNERS.txt` sticks to the [CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners) format,
but adds the ability to include teams defined in `virtual-teams.json`.

For example a CODEOWNERS file can look like this:

```CODEOWNERS
#! comments that start with #! won't appear in the CODEOWNERS output
#!
#! this is not the CODEOWNERS file - to get that one run
#!   npx virtual-code-owners
#!
# catch-all to ensure there at least _is_ a code owner, even when
# it's _everyone_

* @cloud-heroes-all

# admin & ci stuff => transversal

.github/            @ch/transversal

# generic stuff

apps/framework/     @ch/transversal
apps/ux-portal/     @ch/ux @ch/transversal
libs/components/    @ch/ux

# specific functionality

libs/sales/         @ch/sales
libs/after-sales/   @ch/after-sales
libs/refund/        @ch/sales @ch/after-sales
libs/baarden/       @ch/mannen-met-baarden
```

... where only the @cloud-heroes-all is a 'real' team on GitHub level. The other
ones are defined in `virtual-teams.json`.

### virtual-teams.json

A valid JSON file that contains a list of teams and their members.
Update it whenever you have new team members and run `vcodeowners`
to keep CODEOWNERS current.

```json
{
  "ch/after-sales": [
    "john-doe-ch",
    "pete-peterson-ch",
    "john-galt-ch",
    "daisy-duck",
    "donald-duck"
  ],
  "ch/sales": [
    "gregory-gregson-ch",
    "jane-doe-ch",
    "abraham-ableton-ch",
    "dagny-taggert-ch",
    "karl-marx-ch"
  ],
  "ch/ux": [
    "davy-davidson-ch",
    "john-johnson-ch",
    "joe-dalton-ch",
    "koos-koets"
  ],
  "ch/transversal": [
    "mary-the-merry-ch",
    "luke-the-lucky-ch",
    "benjamin-franklin",
    "koos-koets",
    "abraham-lincoln"
  ],
  "ch/mannen-met-baarden": [
    "jan@example.com",
    "pier@example.com",
    "tjorus@example.com",
    "korneel@example.com"
  ]
}
```

### CODEOWNERS

Running `vcodeowners` will combine these into a CODEOWNERS file like this:

```CODEOWNERS
#
# DO NOT EDIT - this file is generated and your edits will be overwritten
#
# To make changes:
#
#   - edit .github/VIRTUAL-CODEOWNERS.txt
#   - and/ or add team members to .github/virtual-teams.json
#   - run 'vcodeowners'
#

# catch-all to ensure there at least _is_ a code owner, even when
# it's _everyone_

* @cloud-heroes-all

# admin & ci stuff => transversal

.github/            @abraham-lincoln @benjamin-franklin @koos-koets @luke-the-lucky-ch @mary-the-merry-ch

# generic stuff

apps/framework/     @abraham-lincoln @benjamin-franklin @koos-koets @luke-the-lucky-ch @mary-the-merry-ch
apps/ux-portal/     @abraham-lincoln @benjamin-franklin @davy-davidson-ch @joe-dalton-ch @john-johnson-ch @koos-koets @luke-the-lucky-ch @mary-the-merry-ch
libs/components/    @davy-davidson-ch @joe-dalton-ch @john-johnson-ch @koos-koets

# specific functionality

libs/sales/         @abraham-ableton-ch @dagny-taggert-ch @gregory-gregson-ch @jane-doe-ch @karl-marx-ch
libs/after-sales/   @daisy-duck @donald-duck @john-doe-ch @john-galt @pete-peterson-ch
libs/refund/        @abraham-ableton-ch @dagny-taggert-ch @daisy-duck @donald-duck @gregory-gregson-ch @jane-doe-ch @john-doe-ch @john-galt @karl-marx-ch @pete-peterson-ch
libs/baarden/       jan@example.com korneel@example.com pier@example.com tjorus@example.com
```

## FAQ

### Any gotcha's?

- It won't check whether the users or teams you entered exist.

### Do I have to run this each time I edit `VIRTUAL-CODEOWNERS.txt`?

Yes. But please automate this for your own sake.

You can for instance set up a rule for `lint-staged` in a `.lintstagedrc.json`
like this:

```json
{
  ".github/{VIRTUAL-CODEOWNERS.txt,virtual-teams.json}": [
    "vcodeowners",
    "git add ."
  ]
}
```

### Can I mix real and virtual teams in `VIRTUAL-CODEOWNERS.txt`?

Yes.

It might be you already have a team or two defined, but just want to use
_additional_ teams. In that case just don't specify the already-defined teams
in `virtual-teams.json` and _vcodeowners_ will leave them alone.

### Can I still use usernames in `VIRTUAL-CODEOWNERS.txt`?

Yes.

Just make sure there's no name clashes between the username and a (virtual)
team name and _vcodeowners_ will leave the real name alone.

### Can I automatically label PR's for virtual teams?

Yep.

Use [actions/labeler](https://github.com/actions/labeler) and tickle
`vcodeowners` to generate the labeler config file:

```sh
vcodeowners --emitLabeler
# Wrote .github/CODEOWNERS AND .github/labeler.json
```

If you have an alternate file location for the `labeler.yml` you can specify that
with vcodeowners's `--labelerLocation` parameter.

> [!NOTE]
> actions/labeler changed the labeler.yml format from v4 to v5.
>
> vcodeowners generates labeler.yml v5 format exclusively, if you're still
> on v4 you will want to upgrade anyway. see [actions/labeler#v5](https://github.com/actions/labeler/blob/8558fd74291d67161a8a78ce36a881fa63b766a9/README.md)
> for details.


### What validations does vcodeowners perform?

vcodeowners checks for basic CODEOWNERS format errors and invalid
user/team names but doesn't verify their existence in the project.

- valid user/team names start with an `@` or are an e-mail address
- valid rules have a file pattern and at least one user/team name
- valid sections headings comply with the syntax described over at [GitLab](https://docs.gitlab.com/ee/user/project/codeowners/reference.html#sections)
  > different from GitLab's syntax the line `[bla @group` is not interpreted
  > as a rule, but as an erroneous section heading. This behaviour might change
  > to be the same as GitLab's in future releases without a major version bump.

### I want to specify different locations for the files (e.g. because I'm using GitLab)

Here you go:

```
vcodeowners \
  --virtualCodeOwners .gitlab/VIRTUAL-CODEOWNERS.txt \
  --virtualTeams      .gitlab/virtual-teams.json \
  --codeOwners        .gitlab/CODEOWNERS
```

### Can I just validate VIRTUAL-CODEOWNERS.txt & virtual-teams.yml without generating output?

Sure thing. Use `--dryRun`:

```
vcodeowners --dryRun
# Wrote '.github/CODEOWNERS' (dry run)
```

### Why the `.txt` extension?

It keeps editors and IDE's from messing up your formatting.

Various editors assume an ALL_CAPS file name with `#` characters on various lines
to be markdown, and will auto format them as such. This makes for either very ugly
or in worst cases invalid CODEOWNERS files. Usually such autoformatting is not
present on text files.

Apparently these editors know about CODEOWNERS, though, so they don't mess with
the formatting of _those_.

### Why does this exist at all? Why not just use GitHub teams?

You should _totally_ use GitHub teams! If you can.

Organizations sometimes have large mono repositories with many code owners.
They or their bureaucracy haven't landed on actually using GitHub teams to
demarcate that. Or you're working on a cross-functional team that doesn't follow
the organization chart (and hence the GitHub teams). Teams in those organizations
who want to have clear code ownership can either:

- Wrestle the bureaucracy.  
  Recommended! It might take a while, though - and even though there are good
  people on many levels in bureaucracies, it might eventually not pan out
  because #reasons.
- Maintain a CODEOWNERS file with code assigned to large lists of individuals.  
  An option, but laborious to maintain, even for smaller projects

This is where `virtual-code-owners` and `vcodeowners` come in.
