# esplus

A cli helper tool for [espanso](https://espanso.org/).

It has two commands:
- `template` that allows you to use golang templates as espanso variables,
- `run` that allows you to run a command without waiting for it to finish (returns empty string).

## Installation

Dowload it from the [releases page](https://github.com/kpym/esplus/releases) and put it in your path.
Or build it yourself:

```bash
go install github.com/kpym/esplus@latest
```

## Usage

### template

If there is a single parameter it is passed as a string, else the parameters are passed as array of strings.

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

The following espanso trigger will replace `!snippet` with the `snippet.txt` file content used as a template. 

```yaml
  - trigger: "!snippet"
    replace: "{{output}}"
    vars:
      - name: "ask"
        type: form
        params:
          layout: "Name [[name]], Age [[age]]"
      - name: output
        type: script
        params:
          args:
            - esplus
            - template
            - "full/path/to/snippet.txt"
            - "{{ask.name}}"
            - "{{ask.age}}"
```

The file `snippet.txt` could looks like this:

```txt
The name is {{index . 0}} and the age is {{index . 1}}.
```

If the file is not found, it is interpreted as a template string, so probably it will be returned as is.

### run

The following espanso trigger will :
- immediately return an empty string,
- wait for 210 ms, the time for espanso to remove `!edit` (replace it with the empty string),
- then open the espanso config folder in vscode.

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

**NOTE:** I hope this tool to become usless one day (see [espanso#1449](https://github.com/espanso/espanso/discussions/1449) and [espanso#1415](https://github.com/espanso/espanso/discussions/1415)).
