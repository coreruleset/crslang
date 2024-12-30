// Copyright 2023 Felipe Zipitria
// SPDX-License-Identifier: Apache-2.0

package parsing

//go:generate java -Xmx500M -cp "../seclang_parser/lib/antlr-4.13.2-complete.jar" org.antlr.v4.Tool -Dlanguage=Go -no-visitor -package parsing ../seclang_parser/parser/SecLangLexer.g4 ../seclang_parser/parser/SecLangParser.g4 -Xexact-output-dir -o .
