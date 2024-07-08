
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

