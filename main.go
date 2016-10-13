/*
 * Copyright (c) 2015 Sung Pae <self@sungpae.com>
 */

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

func abort(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(1)
}

func withReadonlyBoltDB(boltpath string, f func(db *bolt.DB) error) error {
	if _, err := os.Stat(boltpath); os.IsNotExist(err) {
		return err
	}

	opts := bolt.DefaultOptions
	opts.ReadOnly = true
	db, err := bolt.Open(boltpath, 0600, opts)
	if err != nil {
		return err
	}

	return f(db)
}

func dumpJSON(db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		fmt.Print("{")
		dumpBucket(tx.Cursor().Bucket())
		fmt.Print("}")
		return nil
	})
}

func dumpBucket(b *bolt.Bucket) {
	c := b.Cursor()
	k, v := c.First()
	i := 0

	for {
		if k == nil {
			return
		}

		if i > 0 {
			fmt.Print(",")
		}

		if v == nil {
			fmt.Printf("%#v:{", string(k))
			dumpBucket(b.Bucket(k))
			fmt.Print("}")
		} else {
			fmt.Printf("%#v:%#v", string(k), string(v))
		}

		k, v = c.Next()
		i++
	}
}

func main() {
	if len(os.Args) != 2 {
		abort(errors.New("usage: bolt2json dbpath"))
	}

	if err := withReadonlyBoltDB(os.Args[1], dumpJSON); err != nil {
		abort(err)
	}
}
