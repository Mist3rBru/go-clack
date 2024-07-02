# `go-clack/prompts`

Effortlessly build beautiful command-line apps 🪄 [Try the demo](https://stackblitz.com/edit/clack-prompts?file=index.js)

![clack-prompt](https://github.com/Mist3rBru/go-clack/blob/main/.github/assets/clack-demo.gif)

---

`go-clack/prompts` is an opinionated, pre-styled wrapper around [`go-clack/core`](https://www.npmjs.com/package/go-clack/core).

- 🤏 80% smaller than other options
- 💎 Beautiful, minimal UI
- ✅ Simple API
- 🧱 Comes with `Text`, `Confirm`, `Select`, `MultiSelect`, and `Spinner` components

## Basics

### Setup

The `intro` and `outro` functions will print a message to begin or end a prompt session, respectively.

```go
prompts.Intro("create-my-app")
// Do stuff
prompts.Outro("You're all set!")
```

### Cancellation

An `error` is returned when a user cancels a prompt with `CTRL + C`. You should handle this situation for each prompt, optionally providing a nice cancellation message with the `Cancel` utility.

```go
value, err := prompts.Text(/* TODO */)

if (err != nil) {
  prompts.Cancel("Operation cancelled.")
  os.Exit(0)
}
```

## Components

### Text

The `Text` component accepts a single line of text.

```go
meaning, err := prompts.Text(prompts.TextParams{
  Message: "What is the meaning of life?",
  Placeholder: "Not sure",
  InitialValue: "42",
  Validate: func(value string) error {
    if value.length === 0 {
      return errors.New("Value is required!")
    }
    return nil
  },
})
```

### Password

The `Password` component accepts a password input, masking the characters.

```go
password, err := prompts.Password(prompts.PasswordParams{
  Message: "Enter your password:",
  Validate: func(value string) error {
    if value.length < 8 {
      return errors.New("Password must be at least 8 characters long!")
    }
    return nil
  },
})
```

### Path

The `Path` component accepts a file or directory path.

```go
path, err := prompts.Path(prompts.PathParams{
  Message: "Enter the file path:",
  Placeholder: "/path/to/file",
  Validate: func(value string) error {
    if !fileExists(value) {
      return errors.New("File does not exist!")
    }
    return nil
  },
})
```

### Confirm

The `Confirm` component accepts a yes or no answer. The result is a boolean value of `true` or `false`.

```go
confirmed, err := prompts.Confirm(prompts.ConfirmParams{
  Message: "Are you sure?",
  InitialValue: true,
})
```

### Select

The `Select` component allows the user to choose a single option from a list.

```go
project, err := prompts.Select(prompts.SelectParams{
  Message: "Pick a project type:",
  Options: []prompts.SelectOption[string]{
			{Label: "TypeScript", Value: "ts"},
			{Label: "JavaScript", Value: "js"},
			{Label: "CoffeeScript", Value: "coffee", Hint: "oh no"},
		},
})
```

### MultiSelect

The `MultiSelect`component allows the user to choose multiple options from a list.

```go
additionalTools, err := prompts.MultiSelect(prompts.MultiSelectParams{
  Message: "Select additional tools:",
  Options: []prompts.MultiSelectOption[string]{
    {Value: "eslint", Label: "ESLint", Hint: "recommended"},
    {Value: "prettier", Label: "Prettier"},
    {Value: "gh-action", Label: "GitHub Action"},
  },
})
```

### GroupMultiSelect

The `GroupMultiSelect` component allows the user to choose multiple options from grouped lists.

```go
groupChoices, err := prompts.GroupMultiSelect(prompts.GroupMultiSelectParams{
  Message: "Select additional tools:",
  Options: map[string][]prompts.MultiSelectOption[string]{
    "Group 1": {
      {Label: "Option 1", Value: "1"},
      {Label: "Option 2", Value: "2"},
    },
    "Group 2": {
      {Label: "Option 3", Value: "3"},
      {Label: "Option 4", Value: "4"},
    },
  },
})
```

### SelectPath

The `SelectPath` component allows the user to select a file/folder on a tree based select with free navigation by arrow keys.

```go
selectedPath, err := prompts.SelectPath(prompts.SelectPathParams{
  Message: "Select a path:",
})
```

### MultiSelectPath

The `MultiSelectPath` component allows the user to select multiple files/folders on a tree based select with free navigation by arrow keys.

```go
selectedPaths, err := prompts.MultiSelectPath(prompts.MultiSelectPathParams{
  Message: "Select paths",
})
```

### SelectKey

The `SelectKey` component allows the user to choose a option associated to a key.

```go
selectedKey, err := prompts.SelectKey(prompts.SelectKeyParams{
  Message: "Select a key",
  Options: []prompts.SelectKeyOption[string]{
    {Key: "f", Label: "Foo"},
    {Key: "b", Label: "Bar"},
    {Key: "Enter", Label: "Baz"},
  },
})
```

### Spinner

The spinner component surfaces a pending action, such as a long-running download or dependency installation.

```go
s := prompts.Spinner()
s.Start("Installing via npm")
// Do installation here
s.Stop("Installed via npm")
```

## Utilities

### Tasks

Execute multiple tasks in spinners.

```go
prompts.Tasks(
  context.Background(),
  []prompts.Task{
    {
      Title: "Installing via npm",
      Task: func(message func(msg string)) (string, error) {
        // Do installation here
        return "Installed via npm", nil
      }
    },
  },
  prompts.SpinnerOptions{}
)
```
