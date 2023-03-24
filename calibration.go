package main

import (
  "encoding/json"
  "fmt"
  "math"
  "strings"
)

type AxesCoordinates struct {
  X, Y, Z int64
}

type CalibrationMatrix struct {
  X, Y, Z *Metric
}

func (l *ParsedAccelLog) FindCalibrationMatrix(thresholdVal, residualVal int64) *CalibrationMatrix {
  if l == nil {
    return nil
  }
  calMtr := &CalibrationMatrix{}

  for _, metric := range l.Metrics {
    switch {
    case calMtr.X == nil && metric.AxesCrd.hasCalibrationAxis(calAxisX, thresholdVal, residualVal):
      calMtr.X = metric
    case calMtr.Y == nil && metric.AxesCrd.hasCalibrationAxis(calAxisY, thresholdVal, residualVal):
      calMtr.Y = metric
    case calMtr.Z == nil && metric.AxesCrd.hasCalibrationAxis(calAxisZ, thresholdVal, residualVal):
      calMtr.Z = metric
    }
    if calMtr.X != nil && calMtr.Y != nil && calMtr.Z != nil {
      break
    }
  }

  return calMtr
}

func (c *AxesCoordinates) hasCalibrationAxis(calibrationAxis int, calThresholdVal, calResidualVal int64) bool {
  accepted := func(crdVal, thresholdVal int64) bool {
    return math.Abs(float64(crdVal-thresholdVal)) <= float64(calResidualVal)
  }
  var (
    hasAccepted bool
  )
  switch calibrationAxis {
  case calAxisX:
    hasAccepted = accepted(c.X, calThresholdVal) && accepted(c.Y, calZeroVal) && accepted(c.Z, calZeroVal)
  case calAxisY:
    hasAccepted = accepted(c.X, calZeroVal) && accepted(c.Y, calThresholdVal) && accepted(c.Z, calZeroVal)
  case calAxisZ:
    hasAccepted = accepted(c.X, calZeroVal) && accepted(c.Y, calZeroVal) && accepted(c.Z, calThresholdVal)
  }
  return hasAccepted
}

func (m *CalibrationMatrix) toMap() (map[string]*Metric, error) {
  mtrMap := map[string]*Metric{}
  b, err := json.Marshal(m)
  if err != nil {
    return nil, err
  }
  if err = json.Unmarshal(b, &mtrMap); err != nil {
    return nil, err
  }
  return mtrMap, err
}

func (m *CalibrationMatrix) FormInfo() (string, error) {
  mtrMap, err := m.toMap()
  if err != nil {
    return "", fmt.Errorf("cannot convert calibration mtr struct to map: %v", err)
  }
  sb := strings.Builder{}
  sb.WriteString("\n[calibration matrix report]\n")

  for axis, metric := range mtrMap {
    if metric == nil || metric.AxesCrd == nil {
      sb.WriteString(fmt.Sprintf("calibration axis '%s' not found for log\n", axis))
      continue
    }
    sb.WriteString(
      fmt.Sprintf("found '%s' calibration axis. source log: '%s'. coordinates: [x,y,z]=[%d,%d,%d]\n",
        axis, metric.SourceLog, metric.AxesCrd.X, metric.AxesCrd.Y, metric.AxesCrd.Z,
      ))
  }
  sb.WriteString("[end of report]")

  return sb.String(), nil
}
