/**
 * SecLang language support for CodeMirror 6
 */

import { LRLanguage, LanguageSupport } from "@codemirror/language";
import { styleTags, tags as t } from "@lezer/highlight";
import { parser } from "./seclang-parser.js";
import {
  DirectiveName,
  Variable,
  Operator,
  Action,
  QuotedString,
  Number,
  Comment,
  Word,
  LineBreak
} from "./seclang-parser.terms.js";

export const seclangLanguage = LRLanguage.define({
  parser: parser.configure({
    props: [
      styleTags({
        [DirectiveName]: t.keyword,
        [Variable]: t.variableName,
        [Operator]: t.operator,
        [Action]: t.function(t.propertyName),
        [QuotedString]: t.string,
        [Number]: t.number,
        [Comment]: t.lineComment,
        [Word]: t.atom,
        [LineBreak]: t.separator,
      })
    ]
  }),
  languageData: {
    commentTokens: { line: "#" },
    indentOnInput: /^\s*$/,
  }
});

export function seclang() {
  return new LanguageSupport(seclangLanguage);
} 