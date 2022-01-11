# Reddit DL

This is a tool that downloads videos from a line break seperated list of reddit post URLs.

## Dependencies

If you want to download videos that contain audio, you'll need to install
`ffmpeg` and make sure it is in your systems `PATH` environment variable.

## Build

Download the latest golang version and run:

```
go build -o redditdl .
```

## Usage

### As a CLI

Run the following in your terminal (powershell if you are on windows):

```
cat urls.txt | redditdl
```

### As a library

TODO
