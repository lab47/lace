
# LACE - Lab47 Accelerated Computing Environment

Lace is a clojure environment forked from https://github.com/candid82/joker and embued with enhanced powers, such as:

* A Bytecode VM to speed up interpretation
* Extensive interact with Go types
* Detailed stacktraces that can alternate between clojure code and Go code

## Installation

For now, build it via `go install`.

## Usage

`lace` - launch REPL

`lace <filename>` - execute a script. Lace uses `.clj` filename extension. For example: `lace foo.clj`. Normally exits after executing the script, unless `--exit-to-repl` is specified before `--file <filename>`
in which case drops into the REPL after the script is (successfully) executed. (Note use of `--file` in this case, to ensure `<filename>` is not treated as a `<socket>` specification for the repl.)

`lace -` - execute a script on standard input (os.Stdin).

## Project goals

Lace is designed to be a dynamic glue language for Go packages. It leans fully into Greenspun's 10th rule:

> Any sufficiently complicated C or Fortran program contains an ad hoc, informally-specified, bug-ridden, slow implementation of half of Common Lisp.

In this case, we're picking Clojure rather than Common Lisp because it's the most commonly used lisp today.

The project began by forking Joker and beginning to add extensive reflection capabilities to allow clojure code to interact with Go types and code.

## Differences with Clojure

1. Primitive types are different due to a different host language and desire to simplify things. Scripting doesn't normally require all the integer and float types, for example. Here is a list of Lace's primitive types:

  | Lace type | Corresponding Go type |
  |------------|-----------------------|
  | BigFloat   | big.Float             |
  | BigInt     | big.Int               |
  | Boolean    | bool                  |
  | Char       | rune                  |
  | Double     | float64               |
  | Int        | int                   |
  | Keyword    | n/a                   |
  | Nil        | n/a                   |
  | Ratio      | big.Rat               |
  | Regex      | regexp.Regexp         |
  | String     | string                |
  | Symbol     | n/a                   |
  | Time       | time.Time             |

  Note that `Nil` is a type that has one value `nil`.

1. The set of persistent data structures is much smaller:

  | Lace type | Corresponding Clojure type |
  | ---------- | -------------------------- |
  | ArrayMap   | PersistentArrayMap         |
  | MapSet     | PersistentHashSet (or hypothetical PersistentArraySet, depending on which kind of underlying map is used) |
  | HashMap    | PersistentHashMap          |
  | List       | PersistentList             |
  | Vector     | PersistentVector           |

1. The following features are not implemented: protocols, records, structmaps, chunked seqs, transients, tagged literals, unchecked arithmetics, primitive arrays, custom data readers, transducers, validators and watch functions for vars and atoms, hierarchies, sorted maps and sets.
1. Unrelated to the features listed above, the following function from clojure.core namespace are not currently implemented but will probably be implemented in some form in the future: `subseq`, `iterator-seq`, `reduced?`, `reduced`, `mix-collection-hash`, `definline`, `re-groups`, `hash-ordered-coll`, `enumeration-seq`, `compare-and-set!`, `rationalize`, `load-reader`, `find-keyword`, `comparator`, `resultset-seq`, `file-seq`, `sorted?`, `ensure-reduced`, `rsubseq`, `pr-on`, `seque`, `alter-var-root`, `hash-unordered-coll`, `re-matcher`, `unreduced`.
1. Built-in namespaces have `lace` prefix. The core namespace is called `lace.core`. Other built-in namespaces include `lace.string`, `lace.json`, `lace.os`, `lace.base64` etc. See [standard library reference](https://candid82.github.io/lace/) for details.
1. Miscellaneous:
  - `case` is just a syntactic sugar on top of `condp` and doesn't require options to be constants. It scans all the options sequentially.
  - `slurp` only takes one argument - a filename (string). No options are supported.
  - `ifn?` is called `callable?`
  - Map entry is represented as a two-element vector.
  - resolving unbound var returns `nil`, not the value `Unbound`. You can still check if the var is bound with `bound?` function.

## Coding Guidelines

- Dashes (`-`) in namespaces are not converted to underscores (`_`) by Lace, so (unlike with Clojure) there's no need to name `.clj` files accordingly.
- Avoid `:refer :all` and the `use` function, as that reduces the effectiveness of linting.

## License


```
Copyright (c) Evan PHoenix. All rights reserved.
Copyright (c) Roman Bataev (from Joker). All rights reserved.
The use and distribution terms for this software are covered by the
Eclipse Public License 1.0 (http://opensource.org/licenses/eclipse-1.0.php)
which can be found in the LICENSE file.
```

Lace contains parts of Clojure source code (from `clojure.core` namespace). Clojure is licensed as follows:

```
Copyright (c) Rich Hickey. All rights reserved.
The use and distribution terms for this software are covered by the
Eclipse Public License 1.0 (http://opensource.org/licenses/eclipse-1.0.php)
which can be found in the file epl-v10.html at the root of this distribution.
By using this software in any fashion, you are agreeing to be bound by
the terms of this license.
You must not remove this notice, or any other, from this software.
```


