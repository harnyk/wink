# wink

## What is it?

Wink is a simple, lightweight, and fast PeopleHR API command line client.

It has a very limited feature set, including:

* Check-in
* Check-out
* Get current timesheet

## Installation

Download the latest release from the [releases page](https://github.com/harnyk/wink/releases) and extract it to a directory in your PATH.

## Usage

```
wink [command]

Commands:
  wink in     Check in
  wink out    Check out
  wink ls     Get current timesheet
  wink init   Initialize wink
```

## Configuration

Wink uses a configuration file located at `~/.wink/secrets`. You can create this file by running `wink init`.

You will be prompted for your PeopleHR API key and your PeopleHR user ID.

After that you will be prompted for the password, which will be the encryption key for your secrets file, containing your API key and user ID.

## License

WTFPL

