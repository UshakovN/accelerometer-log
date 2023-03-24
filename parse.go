package main

import (
  "fmt"
  "io"
  "log"
  "os"
  "strconv"
  "strings"
)

type ParsedAccelLog struct {
  Meta    *MetaInfo
  Metrics []*Metric
}

type MetaInfo struct {
  TotalMatchedLines int
  TotalMatchedParts int
  TotalErrorsCount  int
}

type Metric struct {
  SourceLog string
  AxesCrd   *AxesCoordinates
}

func ParseAccelLog(filePath string, enableLogOpt ...bool) (*ParsedAccelLog, error) {
  logger := log.Default()
  if len(enableLogOpt) == 0 || !enableLogOpt[0] {
    logger.SetOutput(io.Discard)
  }
  defer logger.SetOutput(os.Stdout)

  file, err := os.ReadFile(filePath)
  if err != nil {
    return nil, fmt.Errorf("cannot read file: %v", err)
  }
  content := stripLineBreaks(string(file))

  logLineMatches := regLogLine.FindAllString(content, matchAll)

  var (
    totalMatchedLines int
    totalMatchedParts int
    totalErrCount     int
  )

  totalMatchedLines = len(logLineMatches)
  logger.Printf("matched '%d' lines from logger", totalMatchedLines)

  var (
    axisParts          []string
    decAxesCoordinates []int64
    axisValues         []*Metric
    parseCrdErrCount   int
  )

  for _, logLineMatch := range logLineMatches {
    logLineMatch = stripSpaces(logLineMatch)

    originalLineMatch := logLineMatch
    logger.Printf("matched line: '%s'", originalLineMatch)

    logLineMatch = strings.TrimPrefix(logLineMatch, prefixPart)
    logLineMatch = strings.TrimSuffix(logLineMatch, suffixPart)

    logPartMatches := regLogPart.FindAllString(logLineMatch, matchAll)

    totalMatchedParts += len(logPartMatches)

    for _, logPart := range logPartMatches {
      logPart = stripSpaces(logPart)
      logger.Printf("line: '%s' extracted part: '%s'", originalLineMatch, logPart)

      axisParts = append(axisParts, logPart)
      if len(axisParts) < partsInAxisCrd {
        continue
      }
      reverseSlice(axisParts)
      axisCrd := strings.Join(axisParts, "")

      axisParts = axisParts[:0]

      decAxisCrd, err := strconv.ParseInt(axisCrd, hexBase, bitSize)
      if err != nil {
        parseCrdErrCount++
        err = fmt.Errorf("cannot parse hex axis coordinate: %v", err)

        if parseCrdErrCount <= possibleParseCrdErrCount {
          logger.Println(err)
          continue
        }
        return nil, err
      }

      decAxesCoordinates = append(decAxesCoordinates, decAxisCrd)
      if len(decAxesCoordinates) < axisCount {
        continue
      }

      axisValues = append(axisValues, &Metric{
        SourceLog: originalLineMatch,
        AxesCrd: &AxesCoordinates{
          X: decAxesCoordinates[idxAxisX],
          Y: decAxesCoordinates[idxAxisY],
          Z: decAxesCoordinates[idxAxisZ],
        },
      })

      decAxesCoordinates = decAxesCoordinates[:0]
    }
  }
  totalErrCount = parseCrdErrCount

  return &ParsedAccelLog{
    Metrics: axisValues,
    Meta: &MetaInfo{
      TotalMatchedLines: totalMatchedLines,
      TotalMatchedParts: totalMatchedParts,
      TotalErrorsCount:  totalErrCount,
    },
  }, nil
}
