# go-mock-gen

| Usage after `Method(t)`                   | Arguments       | Description                               |
|-------------------------------------------|-----------------|-------------------------------------------|
| `<empty>`                                 | -               | only expect called                        |
| `.Return(...)`                            | -               | stub return, ignore args                  |
| `.CalledWith(...)`                        | all, value      | match all args by value, return zero      |
| `.CalledWith(...).Return(...)`            | all, value      | match all args by value, stub return      |
| `.CalledWith[Arg](...)`                   | partial, value  | match ctx by value, return zero           |
| `.CalledWith[Arg](...).Return(...)`       | partial, value  | match single arg by value, stub return    |
| `.Match(func(...) bool)`                  | all, custom     | match all args by callback, return zero   |
| `.Match(func(...) bool).Return(...)`      | all, custom     | match all args by callback, stub return   |
| `.Match[Arg](func(...) bool)`             | partial, custom | match single arg by callback, return zero |
| `.Match[Arg](func(...) bool).Return(...)` | partial, custom | match single arg by callback, stub return |
