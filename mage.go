//go:build mage

package main

import (
	"os"

	"github.com/magefile/mage/sh"
)


func Run() error {
    
    if err := os.Chdir("seclang_parser/parser"); err != nil {
        return err
    }
	if err := sh.Run("java", "-Xmx500M", "-cp" ,"../lib/antlr-4.13.0-complete.jar:$CLASSPATH", "org.antlr.v4.Tool", "-Dlanguage=Go", "-no-visitor", "-package", "parsing", "-o", "../../parsing", "SecLangParser.g4", "SecLangLexer.g4"); err != nil {
        return err
    }
    if err := os.Chdir("../.."); err != nil {
        return err
    }
    return sh.Run("go", "run", ".")
}

func Build() error {
    if err := os.Chdir("seclang_parser/parser"); err != nil {
        return err
    }
	if err := sh.Run("java", "-Xmx500M", "-cp" ,"../lib/antlr-4.13.0-complete.jar:$CLASSPATH", "org.antlr.v4.Tool", "-Dlanguage=Go", "-no-visitor", "-package", "parsing", "-o", "../../parsing", "SecLangParser.g4", "SecLangLexer.g4"); err != nil {
        return err
    }
    if err := os.Chdir("../.."); err != nil {
        return err
    }
    return sh.Run("go", "build", ".")
}