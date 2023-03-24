package main

import "regexp"

const (
  patternLogLine    = `( ?AA AA ?)(\w{2} ){6}( ?FF FF ?)`
  patternLogPart    = `\w{2}`
  patternLineBreaks = "\r\n"
)

const (
  prefixPart = "AA AA"
  suffixPart = "FF FF"
  filePath   = "log-sample.log"
)

const (
  matchAll                 = -1
  partsInAxisCrd           = 2
  hexBase                  = 16
  bitSize                  = 64
  axisCount                = 3
  idxAxisX                 = 0
  idxAxisY                 = 1
  idxAxisZ                 = 2
  possibleParseCrdErrCount = 10
)

var (
  regLogLine    = regexp.MustCompile(patternLogLine)
  regLogPart    = regexp.MustCompile(patternLogPart)
  regLineBreaks = regexp.MustCompile(patternLineBreaks)
)

const (
  calAxisX   = 0
  calAxisY   = 1
  calAxisZ   = 2
  calZeroVal = 0
)
