# wink

## What is it?

Wink is a simple, lightweight, and fast PeopleHR API command line client.

It has a very limited feature set, including:

-   Check-in
-   Check-out
-   Get current timesheet

## Installation

### From source

```
go install github.com/harnyk/wink/cmd/wink
```

### From binary

Download the latest release from the [releases page](https://github.com/harnyk/wink/releases) and extract it to a directory in your PATH.

### Using [eget](https://github.com/zyedidia/eget) (recommended)

Read the [eget](https://github.com/zyedidia/eget) manual for more information.

In order to take all advantages of using `eget` package manager,
you should have some bin directory in your PATH
and specify it in .eget.toml file. For example:

```toml
# File: ~/.eget.toml

[global]
target = "~/bin"
```

Now you can install `wink` using `eget`:

```sh
eget harnyk/wink
```

## Usage

```
Usage:
  wink ls
  wink in [<time>]
  wink out [<time>]
  wink init
  wink report [--start=<start>] [--end=<end>]
  wink --version

Commands:
  ls   - list all my check-ins
  in   - check in to work
  out  - check out of work
  init - setup the API key, and employee ID. Encrypt them using a password
  report - generate a report for the current month

```

## Configuration

Wink uses a configuration file located at `~/.wink/secrets`. You can create this file by running `wink init`.

You will be prompted for your PeopleHR API key and your PeopleHR user ID.

After that you will be prompted for the password, which will be the encryption key for your secrets file, containing your API key and user ID.

## Report

You can generate a report for the current month by running `wink report`.

You can also specify a start and end date using the `--start` and `--end` flags.

Also, you can specify `--output=<path/to/file.json>` to export a report in JSON format.

JSON report is a list of records with the following structure:

```json
[
  {
    "date": "2023-03-01",
    "hours": 9.9,
    "is_complete": true,
    "is_invalid_sequence": false
  },
  {
    "date": "2023-03-02",
    "hours": 7.1,
    "is_complete": true,
    "is_invalid_sequence": false
  },
  ...
]
```

  - `date` - date of the record in `YYYY-MM-DD` format
  - `hours` - number of hours worked on this day
  - `is_complete` - `true` if the record is complete, `false` otherwise. A record is complete if it has both check-in and check-out.
  - `is_invalid_sequence` - `true` if the record has invalid check-in/check-out sequence, `false` otherwise. For example, if you check-in at 10:00 and check-out at 9:00, the record will be invalid.

## License

WTFPL
