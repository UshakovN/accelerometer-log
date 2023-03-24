package main

import "strings"

func reverseSlice[T any](s []T) {
  for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    s[i], s[j] = s[j], s[i]
  }
}

func stripLineBreaks(s string) string {
  return regLineBreaks.ReplaceAllLiteralString(strings.TrimSpace(s), "")
}

func stripSpaces(s string) string {
  return strings.TrimSpace(s)
}
