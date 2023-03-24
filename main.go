package main

import (
  "flag"
  "log"
)

func main() {
  path := flag.String("path", "sample.log", "path to accelerometer log")
  verbose := flag.Bool("verbose", false, "verbose for log parsing")
  threshold := flag.Int64("threshold", 1000, "threshold for calibration matrix")
  residual := flag.Int64("residual", 10, "calibration residual")
  flag.Parse()

  if *path == "" {
    log.Fatalf("file path not specified")
  }
  if *threshold <= 0 {
    log.Fatalf("threshold must be positive integer")
  }
  if *residual < 0 || *residual >= 1000 {
    log.Fatalf("invalid residual value")
  }

  parsedLog, err := ParseAccelLog(*path, *verbose)
  if err != nil {
    log.Fatalf("cannot parse accel log: %v", err)
  }
  calibrationMatrix := parsedLog.FindCalibrationMatrix(*threshold, *residual)

  formedInfo, err := calibrationMatrix.FormInfo()
  if err != nil {
    log.Fatalf("cannot form info for calibration matrix: %v", err)
  }
  log.Println(formedInfo)
}
