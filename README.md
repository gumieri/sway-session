# sway-session

Tool for saving the state of the [Sway WM](https://swaywm.org) session and restoring it.

**At the moment, it is a [PoC](https://en.wikipedia.org/wiki/Proof_of_concept)!**

[![Go Report Card](https://goreportcard.com/badge/github.com/gumieri/note)](https://goreportcard.com/report/github.com/gumieri/note)

## Usage
For saving the running programs and its workspace disposition run the given command:
```bash
sway-session save
```
It will create a file in creates a file at `$XDG_DATA_HOME/sway-session/`

To restore simply use:
```bash
sway-session restore
```
The recomendation would be to place at the sway config file something like that:
```config
exec sway-session restore
```

## Supported programs
Considering that a lot of programs have different ways of retrieving it state and restoring it to the desired state,
the `sway-session` can only offer a generic approach for all the ecosystem and for more specific programs (like terminal-emulators)
to offer some rules with more capabilities.

### → [alacritty](https://github.com/jwilm/alacritty)
 * current working directory
