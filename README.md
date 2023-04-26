# [esplus](https://github.com/kpym/esplus)

A cli helper tool for [espanso](https://espanso.org/).

It has two commands:
- `template` that allows you to use [golang templates](https://pkg.go.dev/text/template) as espanso variables,
- `run` that allows you to run a command without waiting for it to finish (returns empty string).

## Installation

Dowload it from the [releases page](https://github.com/kpym/esplus/releases) and put it in your path.
Or build it yourself:

```bash
go install github.com/kpym/esplus@latest
```

## Usage

### help

When you run `esplus`, it displays the following help message

```
> esplus
esplus is a helper cli for espanso.
Version: 0.3.1
Usage: esplus <command> <args>

Commands:
  template <file> <args> : if file exists, use it as template with args (using {{ and }} as delimiters)
  template <template string> <args> : execute a template with args (using [[ and ]] as delimiters)
  run [milliseconds] <cmd> <args> : run a command (with delay) without waiting for it to finish

Examples:
  esplus template 'Hello [[.|upper]]' 'World'
  esplus template 'Hello [[range .]][[.|upper|printf "%s\n"]][[end]]' 'World' 'and' 'Earth'
  esplus template 'file.template.txt' 'World'
  esplus run 200 code .

Project repository:
  https://github.com/kpym/esplus
```

### template

The templates are executed with the [text/template](https://pkg.go.dev/text/template) golang package. The [sprig](github.com/Masterminds/sprig) functions are available.
If there is a single parameter it is passed as a string, else the parameters are passed as array of strings.

To see how go templates work, you can check [hashicorp's help](https://developer.hashicorp.com/nomad/tutorials/templates/go-template-syntax).

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

## Configuration

The file `~/.config/esplus/config.toml` (if it exists) is used to configure `esplus`. For now, it is only used to provide aliases to some commands. The reasons is that under MacOS `espanso` runs scripts with some minimal PATH, so programs that are in the path could not be found. The aliases are used to provide the full path to that programs. For example, the following `config.toml` file will allow to use `code` and `subl` as aliases for `Visual Studio Code` and `Sublime Text`:

```toml
[aliases]
code = "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code"
subl = "/Applications/Sublime Text.app/Contents/SharedSupport/bin/subl"
```

## Note

I hope this tool to become usless one day (see [espanso#1449](https://github.com/espanso/espanso/discussions/1449) and [espanso#1415](https://github.com/espanso/espanso/discussions/1415)).

## License

[MIT](LICENSE) License
