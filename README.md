# esplus

A cli helper tool for [espanso](https://espanso.org/).

It has two commands:
- `template` that allows you to use golang templates as espanso variables
- `run` that allows you to run a command without waiting for it to finish (returns empty string)

## Installation

Dowload it from the [releases page](https://github.com/kpym/esplus/releases) and put it in your path.
Or build it yourself:

```bash
go install github.com/kpym/esplus@latest
```

## Usage

### template

The following espanso trigger will replace `!lo` with the clipboard content in lowercase.

```yaml
  - trigger: "!lo"
    replace: "{{output}}"
    vars:
      - name: "clipboard"
        type: "clipboard"
      - name: output
        type: script
        params:
          args:
            - esplus
            - template
            - "[[.|lower]]"
            - "{{clipboard}}"
```

### run

The following espanso trigger will open the espanso config folder in vscode. It will wait for 210 ms, the time for espanso to remove `!edit`.

```yaml
  - trigger: "!edit"
    replace: "{{output}}"
    vars:
      - name: output
        type: script
        params:
          args:
            - esplus
            - run
            - "210"
            - code
            - '%CONFIG%'

```
