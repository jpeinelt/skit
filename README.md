# Skit
a basic presentation app for the command line

`→` next slide

`←` previous slide

`Ctrl-C` quit

![Screenshot with Text](/slide-text.png?raw=true)
![Screenshot with Image](/slide-image.png?raw=true)

## Table Of Contents

- [Table Of Contets](#table-of-contents)
- [Features](#features)
  - [Syntax](#syntax)
- [Contributing](#installing-and-contributing)
- [Inspiration](#inspiration)
- [License](#license)

## Features

### Syntax

- Each paragraph marks a slide.
- Text is just starts at the beginning of a line.
- If the first character is a `/` a control sequence is expected, except if another `/`, then it's just a `/` as text:
    - `#`  -> Slide title, required.
    - `@`  -> Path to image, optional. A slide shows either an      image or text.
    - `_`  -> Background color, between 0 and 255.
    - `^` -> Foreground color, between 0 and 255.
    - `!`  -> Comment, is parsed but ignored for now.
- The colors are ignored for images, since images are in ascii and need a black background.

See [example.skt](./example.skt).

### Images

This program supports the following file types:

- .jpg
- .png

## Installing and Contributing

Install with:

```bash
go get github.com/jpeinelt/skit
```

Contributions are always welcome

## Inspiration

### General Idea
[zeichma][zeichma] - A dear colleague of mine came up with the idea to have a
    really simple presentation software. Thank you Karl!

### Lexer and Parser
[Lexing in Go][rob-pike] - Rob Pike talk about lexing

[Ini Parser][ini-parser] - Simple parser, also inspired by Rob Pike's talk

### UI
[slack-term][slack-term] - Slack client for the CLI

## License

Copyright © 2018 Julius Peinelt

This software is distributed under the terms of the Apache License Version 2.0
(Apache-2.0). For details see [LICENSE](./LICENSE).



[zeichma]: https://github.com/fleischie/zeichma
[rob-pike]: https://www.youtube.com/watch?v=HxaD_trXwRE
[ini-parser]: https://github.com/adampresley/sample-ini-parser
[slack-term]: https://github.com/erroneousboat/slack-term
