// Code generated by argsgen.
// DO NOT EDIT!
package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
)

func (o *options) flagSet() *flag.FlagSet {
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.StringVar(&o.db, "db", o.db, "database name (sqlite)")
    flagSet.StringVar(&o.model, "model", o.model, "name of the model")
    return flagSet
}

// Parse parses the arguments in os.Args
func (o *options) Parse() error {
    flagSet := o.flagSet()
    args := os.Args[1:]
    for len(args) > 0 {
        if err := flagSet.Parse(args); err != nil {
            return err
        }

        if remaining := flagSet.NArg(); remaining > 0 {
            posIndex := len(args) - remaining
            args = args[posIndex+1:]
            continue
        }
        break
    }

    if o.db == "" {
        return errors.New("argument 'db' is required")
    }
    if o.model == "" {
        return errors.New("argument 'model' is required")
    }
    return nil
}

// MustParse parses the arguments in os.Args or exists on error
func (o *options) MustParse() {
    if err := o.Parse(); err != nil {
        o.flagSet().PrintDefaults()
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
