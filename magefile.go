//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/mg"
)


func Generate() error {
    err := sh.RunV("go", "generate", "./...");
    return err
}

func Run() error {
    mg.Deps(Generate)
    return sh.Run("go", "run", ".")
}

func Build() error {
    mg.Deps(Generate)
    return sh.Run("go", "build", ".")
}

func Test() error {
    mg.Deps(Generate)
    return sh.Run("go", "test", ".")
}

