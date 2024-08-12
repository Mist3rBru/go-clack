
<a name="v0.1.13"></a>
## [v0.1.13](https://github.com/Mist3rBru/go-clack/compare/v0.1.12...v0.1.13) (2024-08-12)

### 🚀 Features

* **prompts:** add error handling utils


<a name="v0.1.12"></a>
## [v0.1.12](https://github.com/Mist3rBru/go-clack/compare/v0.1.11...v0.1.12) (2024-08-09)

### 🚀 Features

* improve control over empty directories
* add async validation support
* **prompts:** add more methods to Workflow
* **prompts:** add workflow prompt

### 🩹 Fixes

* **prompts:** initial theme with cursor and placeholder

### 🧪 Tests

* **core:** refactor tests and add more tests for Prompt
* **prompts:** add theme tests

### 🛠️ Refactors

* simplify if statements
* **core:** split extra code from prompt.go file into dedicated files
* **core:** add and document available events
* **prompts:** remove generics from Workflow

### 🏡 Chore

* update change set example with workflow prompt
* update change set example with workflow prompt


<a name="v0.1.11"></a>
## [v0.1.11](https://github.com/Mist3rBru/go-clack/compare/v0.1.10...v0.1.11) (2024-08-01)

### 🚀 Features

* **core:** add validations for more types on WrapValidate


<a name="v0.1.10"></a>
## [v0.1.10](https://github.com/Mist3rBru/go-clack/compare/v0.1.9...v0.1.10) (2024-08-01)

### 🚀 Features

* add required options to select prompt
* add filter option to multi select prompt
* add filter option to select prompt
* add sort to path node children
* add filter option to multi select path prompt
* add filter option to select path prompt
* **core:** add IsEqual method to PathNode
* **core:** add OSFileSystem as default for PathNode.FileSystem
* **prompts:** add theme symbol color and bar color

### 🩹 Fixes

* multi select invalid option selection
* **prompts:** Synbol typo

### 🧪 Tests

* improve code coverage to 89.9/92.6

### 🛠️ Refactors

* **core:** turn PathNode.MapChildren into a mutator method
* **core:** add Flat method to PathNode
* **core:** make TrackKeyValue agnostic of Prompt
* **core:** merge WrapValidate functions
* **core:** add PathNode.IsDir indentifier field
* **core:** move OSFileSystem to internals package
* **prompts:** remove context from Spinner

### 🏡 Chore

* update changelog
* **core:** add go docs to Prompt


<a name="v0.1.9"></a>
## [v0.1.9](https://github.com/Mist3rBru/go-clack/compare/v0.1.8...v0.1.9) (2024-07-06)

### 🚀 Features

* **core:** add internal validation of essential params
* **prompts:** add internal validation of essential params

### 🛠️ Refactors

* move utils to dedicated modules
* **core:** simplify prompt constructors
* **core:** add WrapValidate helper function
* **core:** add WrapRender helper function
* **prompts:** connect note borders

### 🏡 Chore

* update changelog


<a name="v0.1.8"></a>
## [v0.1.8](https://github.com/Mist3rBru/go-clack/compare/v0.1.7...v0.1.8) (2024-07-03)

### 🚀 Features

* add DisabledGroups option to GroupMultiSelectPrompt
* add required option to prompts
* **prompts:** add SpacedGroups option to GroupMultiSelect

### 🏡 Chore

* update changelog


<a name="v0.1.7"></a>
## [v0.1.7](https://github.com/Mist3rBru/go-clack/compare/v0.1.6...v0.1.7) (2024-07-03)

### 🚀 Features

* add label as value to prompts
* **prompts:** add multi line support to log functions

### 📖 Documentation

* fix typos
* add readme

### 🏡 Chore

* add code examples
* update changelog


<a name="v0.1.6"></a>
## [v0.1.6](https://github.com/Mist3rBru/go-clack/compare/v0.1.5...v0.1.6) (2024-07-02)

### 🩹 Fixes

* **prompts:** useless Spinner's error

### 🏡 Chore

* add CHANGELOG


<a name="v0.1.5"></a>
## [v0.1.5](https://github.com/Mist3rBru/go-clack/compare/v0.1.4...v0.1.5) (2024-06-26)

### 🩹 Fixes

* **core:** MultiSelectPathPrompt initial value
* **core:** MultiSelectPrompt initial value


<a name="v0.1.4"></a>
## [v0.1.4](https://github.com/Mist3rBru/go-clack/compare/v0.1.3...v0.1.4) (2024-06-23)

### 🩹 Fixes

* Path.OnlyShowDir mapping


<a name="v0.1.3"></a>
## [v0.1.3](https://github.com/Mist3rBru/go-clack/compare/v0.1.2...v0.1.3) (2024-06-13)

### 🩹 Fixes

* **prompts:** add bar to log messages


<a name="v0.1.2"></a>
## [v0.1.2](https://github.com/Mist3rBru/go-clack/compare/v0.1.1...v0.1.2) (2024-06-07)


<a name="v0.1.1"></a>
## [v0.1.1](https://github.com/Mist3rBru/go-clack/compare/v0.1.0...v0.1.1) (2024-06-07)

### 🚀 Features

* **core:** add MultiSelectPathPrompt
* **prompts:** add MultiSelectPath prompt

### 🛠️ Refactors

* change arbitrary prompt state to prompt state contants
* move third_party packages to thid_party/package folder


<a name="v0.1.0"></a>
## v0.1.0 (2024-06-06)

### 🚀 Features

* add multi select prompt
* add confirm prompt
* add base prompt
* add key name literals
* add erase utils
* add utils
* add track cursor value
* add text prompt
* add prompt event name literals
* add prompt options
* add select prompt
* add password prompt
* add select path prompt
* add prompts setup
* TextPrompt placeholder completion
* add default prompt input and output
* format lines method
* add prompt state literals
* add generics to prompts
* add cursor utils
* add buggy limit lines function
* add validate method to prompts
* add select key prompt
* add group multi select prompt
* add path prompt
* **prompts:** add path prompt
* **prompts:** text prompt
* **prompts:** add log prompts
* **prompts:** add Note prompt
* **prompts:** add password prompt
* **prompts:** add MultiSelect prompt
* **prompts:** add select prompt
* **prompts:** add SelectPath prompt
* **prompts:** add Confirm prompt
* **prompts:** add GroupMultiSelect prompt
* **prompts:** add SelectKey prompt
* **prompts:** add Spinner prompt
* **prompts:** add Tasks prompt

### 🩹 Fixes

* extra whitespace on format lines
* format blank line with cursor
* resturn of canceled prompt
* limit lines function
* missing char validation
* close callback
* read reader buffer

### 🧪 Tests

* add test coverage 70%
* add test coverage of 50%
* add text prompt tests
* add base prompt tests

### 🛠️ Refactors

* prepare for external tests
* rename Valeu param to InitialValue
* rename verbose literals
* remove unnecessary mutex implementation
* make LimitLines use internal CursorIndex
* rename Arrow* keys to only arrow name
* use Key struct instead of primitive key
* move globals to globals file
* rename options to params
* add select option struct
* remove default constructors
* **core:** add IsSelected to MultiSelectOption

### 🏡 Chore

* update makefile to support test loop
* adapt to github import
* add config files

