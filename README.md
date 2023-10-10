# Chronomerge

**Automerge with crontab.** Does your organization disable automerge? Run this in your own crontab to automerge your PRs.

## Usage

Build the binary.

```bash
go build ./main.go
```

Then add it to your crontab.

```bash
$ crontab -e
# then add
0 * * * * /path/to/chronomerge/main <your github api token>
```

To automerge a PR, add it to `~/chronomerge.txt`. For example,

```bash
echo "https://github.com/torvalds/linux/pulls/1" >> ~/chronomerge.txt
```

And when it can be merged, Chronomerge will merge it for you.
