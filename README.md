# Wave

Terminal Radio Player

## Requirements

- [mpv](https://mpv.io)

## Install

Downloads the CLI based on your OS/arch and puts it in `/usr/local/bin`.
Needed [jq](https://stedolan.github.io/jq/download) instaled.

```sh
curl -fsSL https://raw.githubusercontent.com/leonardobiffi/wave/master/scripts/install.sh | sh
```

## Configuration

Edit file `~/.wave/stations.yaml` to add new seach parameters stations.

Ex.: `Station Name, City State, Country`

```yaml
stations:
  search:
    - Super Rádio Tupi, Rio de Janeiro RJ, Brazil
    - Kiss FM, São Paulo SP, Brazil
```

## Stations Search Source

This application use the [RadioGarden API](http://radio.garden/) to find radio stations configured on the file
