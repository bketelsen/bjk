# bjk.fyi

## Usage

Download a release of the command line tool.

To use it, place the 'bjk' or 'bjk.exe' binary somewhere in your path.

```
Usage:
  bjk URL SHORTCODE [flags]
  bjk [command]

Available Commands:
  config      create a config file in your home directory
  help        Help about any command
  serve       Run the BJK URL Shortening service

Flags:
  -a, --alias string     BJK Alias
  -c, --channel string   channel
      --config string    config file (default is $HOME/.bjk.yaml)
  -e, --event string     event
  -h, --help             help for bjk
      --server string    URL Shortening server (default "http://bjk.fyi")

Use "bjk [command] --help" for more information about a command.
```

## Getting Started

Run `bjk config`.  This will create a file called `.bjk.yaml` in your $HOME directory.

Edit this file, and replace the alias configuration value with your Microsoft alias.

## Creating a Shortened URL

```
bjk -c twitter -e ignite https://docs.microsoft.com/azure/x?WT.mc_id=opensource-0000-brketels mycode
```

This returns:

```
Using config file: /home/bketelsen/.bjk.yaml
Submitting to  http://bjk.fyi
bjk.fyi/mycode
```

It also copies the shortened URL into your clipboard.  Because I love you guys.

## Details

The shortener uses three flags or config settings to create the tracking link.

```
    -a --alias   :  Your bjk alias  ex: brketels
    -c --channel :  The channel/medium  ex: twitter
    -e --event   :  The event name  ex: ignite
```

If the configuration file at `$HOME/.bjk.yaml` has any of these values they will be defaulted for you and you may exclude them.

For example, if my config file looks like this:

```
Alias: brketels
Channel: twitter
```
Then I can use the `bjk` command and only specify the event with the `-e` flag, the channel and alias will be read from the config file.

You may fill in all three values in the config file, or none.  Any missing values will cause the program to fail with an error:
```
$ bjk -a brketels -e ignite https://microsoft.com/?WT.mc_id=opensource-0000-brketels ms1

> Using config file: /home/bketelsen/.bjk.yaml
> Channel is required.  Set with -c or in config file.
```

