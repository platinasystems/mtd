// Copyright © 2020 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by the GPL-2 license described in the
// LICENSE file.

// Package mtd provides utility service for Linux Memory Technology Devices
package mtd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MTD struct {
	unit int
}

var nameToUnit map[string]*MTD

var ErrMTDNotFound = errors.New("MTD device not found")

// NameToUnit takes a MTD partition name and returns a unit number
func NameToUnit(name string) (unit int, err error) {
	if nameToUnit == nil {
		f, err := os.Open("/proc/mtd")
		if err != nil {
			return 0, fmt.Errorf("Error opening /proc/mtd: %w", err)
		}
		defer f.Close()

		nameToUnit = make(map[string]*MTD)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fields := strings.SplitAfterN(scanner.Text(), " ", 4)
			if !strings.HasPrefix(fields[0], "mtd") {
				continue
			}
			name := strings.TrimRight(fields[3], " ")
			if strings.HasPrefix(name, `"`) && strings.HasSuffix(name, `"`) {
				name = strings.TrimSuffix(strings.TrimPrefix(name, `"`),
					`"`)
			}
			unit, err := strconv.Atoi(strings.TrimSuffix(
				strings.TrimPrefix(
					strings.TrimRight(fields[0], " "),
					"mtd"),
				":"))
			if err != nil {
				return 0, fmt.Errorf("Error converting %s to unit: %w",
					fields[0], err)
			}
			mtd := &MTD{unit: unit}
			nameToUnit[name] = mtd
		}
		if scanner.Err() != nil {
			return 0, fmt.Errorf("Error scanning /proc/mtd: %w",
				scanner.Err())
		}
	}
	mtd := nameToUnit[name]
	if mtd == nil {
		return 0, ErrMTDNotFound
	}
	return mtd.unit, nil
}
