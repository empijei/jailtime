/*
 * jailtime version 0.5
 * Copyright (c)2015-2017 Christian Blichmann
 *
 * Dynamic loader configuration parsing
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package loader // import "blichmann.eu/code/jailtime/loader"

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Default loader config path
const loaderConfig = "/etc/ld.so.conf"

var LdSearchPaths []string = ParseLdConfig(loaderConfig)

func ParseLdConfig(conf string) (paths []string) {
	paths = []string{}
	f, err := os.Open(conf)
	if err != nil {
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var line string
	for {
		line, err = r.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		comp := strings.SplitN(line, " ", 2)
		if len(comp) == 2 && comp[0] == "include" {
			m, err := filepath.Glob(strings.TrimSpace(comp[1]))
			if err != nil {
				continue
			}
			for _, p := range m {
				if newPaths := ParseLdConfig(p); len(newPaths) > 0 {
					paths = append(paths, newPaths...)
				}
			}
		} else {
			paths = append(paths, comp[0])
		}
	}
	return
}

func FindLibraryFunc(basename string, paths []string,
	usable func(path string) bool) string {
	for _, p := range paths {
		full := filepath.Join(p, basename)
		if _, err := os.Stat(full); err == nil && usable(full) {
			return full
		}
	}
	return ""
}
