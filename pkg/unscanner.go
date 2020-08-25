package kgo

import (
    "bufio"
    "io"
    "log"
)

type unScanner struct {
    ok bool
    line string
    scanner *bufio.Scanner
}

func newUnScanner(r io.Reader) *unScanner {
    return &unScanner{
        ok: false,
        scanner: bufio.NewScanner(r),
    }
}

func (s *unScanner) getLine() (line string, ok bool, err error) {
    if s.ok {
        s.ok = false
        ok = true
        line = s.line
    } else {
        ok = s.scanner.Scan()
        if ok {
            line = s.scanner.Text()
        } else {
            err = s.scanner.Err()
        }
    }
    return
}

func (s *unScanner) unGetLine(line string) {
    if s.ok {
        log.Panic("unscanner already has next")
    }
    s.ok = true
    s.line = line
    return
}
