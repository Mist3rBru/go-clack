# `go-clack/core`

Clack contains low-level primitives for implementing your own command-line applications.

Currently exposes `Prompt` as well as:
- `TextPrompt`
- `PasswordPrompt`
- `PathPrompt`
- `ConfirmPrompt`
- `SelectPrompt`
- `MultiSelectPrompt`
- `GroupMultiSelectPrompt`
- `SelectPathPrompt`
- `MultiSelectPathPrompt`
- `SelectKeyPrompt` 

Each `Prompt` accepts a `Render` function.

```go
p := core.NewTextPrompt(core.TextPromptParams{
  Render: func(p *core.TextPrompt) string {
    return fmt.Sprintf("What's your name?\n%s", p.ValueWithCursor)
  },
})

name, err := p.Run()
if (err != nil) {
  os.Exit(0)
}
```
