package gci

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type PkgType int

const (
	// pkg type: standard, remote, local
	standard PkgType = iota
	// 3rd-party packages
	remote
	local
)

const (
	commentFlag     = "//"
	importStartFlag = "\nimport (\n"
	importEndFlag   = "\n)\n"

	blank     = " "
	indent    = "\t"
	linebreak = "\n"
)

type FlagSet struct {
	LocalFlag       []string
	DoWrite, DoDiff *bool
}

type importSpec struct {
	alias   string
	path    string
	comment string
}

type importBlock struct {
	list          map[PkgType][]string
	aboveComment  map[string]string
	inlineComment map[string]string
	alias         map[string]string
}

// ParseLocalFlag takes a comma-separated list of
// package-name-prefixes (as passed to the "-local" flag), and splits
// it in to a list.  This is different than strings.Split in that it
// handles the empty string and empty entries in the list.
func ParseLocalFlag(str string) []string {
	return strings.FieldsFunc(str, func(c rune) bool { return c == ',' })
}

func newImportBlock(data [][]byte, localFlag []string) *importBlock {
	p := &importBlock{
		list:          make(map[PkgType][]string),
		aboveComment:  make(map[string]string),
		inlineComment: make(map[string]string),
		alias:         make(map[string]string),
	}

	formatData := make([]string, 0)
	// remove all empty lines
	for _, v := range data {
		if len(v) > 0 {
			formatData = append(formatData, strings.TrimSpace(string(v)))
		}
	}

	n := len(formatData)
	for i := n - 1; i >= 0; i-- {
		line := formatData[i]

		if strings.HasPrefix(line, commentFlag) {
			// comment in the last line is useless, ignore it
			if i+1 >= n {
				continue
			}
			spec := parseImportSpec(formatData[i+1])
			p.aboveComment[spec.path] = line
			continue
		}

		spec := parseImportSpec(line)

		if spec.alias != "" {
			p.alias[spec.path] = spec.alias
		}
		if spec.comment != "" {
			p.inlineComment[spec.path] = spec.comment
		}

		pkgType := getPkgType(spec.path, localFlag)
		p.list[pkgType] = append(p.list[pkgType], spec.path)
	}

	return p
}

// fmt format import pkgs as expected
func (p *importBlock) fmt() []byte {
	var lines []string

	for _, pkgType := range []PkgType{standard, remote, local} {
		if len(p.list[pkgType]) == 0 {
			continue
		}
		if len(lines) > 0 && lines[len(lines)-1] != "" {
			lines = append(lines, "")
		}
		sort.Strings(p.list[pkgType])
		for _, s := range p.list[pkgType] {
			if p.aboveComment[s] != "" {
				if len(lines) > 0 && lines[len(lines)-1] != "" {
					lines = append(lines, "")
				}
				lines = append(lines, indent+p.aboveComment[s])
			}

			line := s
			if p.alias[s] != "" {
				line = p.alias[s] + blank + s
			}
			if p.inlineComment[s] != "" {
				line += blank + p.inlineComment[s]
			}
			lines = append(lines, indent+line)
		}
	}

	return []byte(strings.Join(lines, linebreak) + linebreak)
}

// parseImportSpec assumes line is a import path, and returns the (path, alias, comment).
func parseImportSpec(line string) importSpec {
	if s := strings.SplitN(line, commentFlag, 2); len(s) > 1 {
		pkgArray := strings.Fields(s[0])
		if len(pkgArray) > 1 {
			return importSpec{
				path:    pkgArray[1],
				alias:   pkgArray[0],
				comment: commentFlag + s[1],
			}
		} else {
			return importSpec{
				path:    pkgArray[0],
				alias:   "",
				comment: commentFlag + s[1],
			}
		}
	} else {
		pkgArray := strings.Fields(line)
		if len(pkgArray) > 1 {
			return importSpec{
				path:    pkgArray[1],
				alias:   pkgArray[0],
				comment: "",
			}
		} else {
			return importSpec{
				path:    pkgArray[0],
				alias:   "",
				comment: "",
			}
		}
	}
}

func getPkgType(line string, localFlag []string) PkgType {
	pkgName := strings.Trim(line, "\"\\`")

	for _, localPkg := range localFlag {
		if strings.HasPrefix(pkgName, localPkg) {
			return local
		}
	}

	if isStandardPackage(pkgName) {
		return standard
	}

	return remote
}

func diff(b1, b2 []byte, filename string) (data []byte, err error) {
	f1, err := writeTempFile("", "gci", b1)
	if err != nil {
		return
	}
	defer os.Remove(f1)

	f2, err := writeTempFile("", "gci", b2)
	if err != nil {
		return
	}
	defer os.Remove(f2)

	cmd := "diff"

	data, err = exec.Command(cmd, "-u", f1, f2).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		return replaceTempFilename(data, filename)
	}
	return
}

func writeTempFile(dir, prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return "", err
	}
	_, err = file.Write(data)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

// replaceTempFilename replaces temporary filenames in diff with actual one.
//
// --- /tmp/gofmt316145376	2017-02-03 19:13:00.280468375 -0500
// +++ /tmp/gofmt617882815	2017-02-03 19:13:00.280468375 -0500
// ...
// ->
// --- path/to/file.go.orig	2017-02-03 19:13:00.280468375 -0500
// +++ path/to/file.go	2017-02-03 19:13:00.280468375 -0500
// ...
func replaceTempFilename(diff []byte, filename string) ([]byte, error) {
	bs := bytes.SplitN(diff, []byte{'\n'}, 3)
	if len(bs) < 3 {
		return nil, fmt.Errorf("got unexpected diff for %s", filename)
	}
	// Preserve timestamps.
	var t0, t1 []byte
	if i := bytes.LastIndexByte(bs[0], '\t'); i != -1 {
		t0 = bs[0][i:]
	}
	if i := bytes.LastIndexByte(bs[1], '\t'); i != -1 {
		t1 = bs[1][i:]
	}
	// Always print filepath with slash separator.
	f := filepath.ToSlash(filename)
	bs[0] = []byte(fmt.Sprintf("--- %s%s", f+".orig", t0))
	bs[1] = []byte(fmt.Sprintf("+++ %s%s", f, t1))
	return bytes.Join(bs, []byte{'\n'}), nil
}

func visitFile(set *FlagSet) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err == nil && isGoFile(f) {
			err = processFile(path, os.Stdout, set)
		}
		return err
	}
}

func WalkDir(path string, set *FlagSet) error {
	return filepath.Walk(path, visitFile(set))
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func ProcessFile(filename string, out io.Writer, set *FlagSet) error {
	return processFile(filename, out, set)
}

func processFile(filename string, out io.Writer, set *FlagSet) error {
	var err error

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	ori := make([]byte, len(src))
	copy(ori, src)
	start := bytes.Index(src, []byte(importStartFlag))
	// in case no importStartFlag or importStartFlag exist in the commentFlag
	if start < 0 {
		fmt.Printf("skip file %s since no import\n", filename)
		return nil
	}
	end := bytes.Index(src[start:], []byte(importEndFlag)) + start

	ret := bytes.Split(src[start+len(importStartFlag):end], []byte(linebreak))

	block := newImportBlock(ret, set.LocalFlag)

	res := append(src[:start+len(importStartFlag)], append(block.fmt(), src[end+1:]...)...)

	if !bytes.Equal(ori, res) {
		if *set.DoWrite {
			// On Windows, we need to re-set the permissions from the file. See golang/go#38225.
			var perms os.FileMode
			if fi, err := os.Stat(filename); err == nil {
				perms = fi.Mode() & os.ModePerm
			}
			err = ioutil.WriteFile(filename, res, perms)
			if err != nil {
				return err
			}
		}
		if *set.DoDiff {
			data, err := diff(ori, res, filename)
			if err != nil {
				return fmt.Errorf("failed to diff: %v", err)
			}
			fmt.Printf("diff -u %s %s\n", filepath.ToSlash(filename+".orig"), filepath.ToSlash(filename))
			if _, err := out.Write(data); err != nil {
				return fmt.Errorf("failed to write: %v", err)
			}
		}
	}
	if !*set.DoWrite && !*set.DoDiff {
		if _, err = out.Write(res); err != nil {
			return fmt.Errorf("failed to write: %v", err)
		}
	}

	return err
}

// Run return source and result in []byte if succeed
func Run(filename string, set *FlagSet) ([]byte, []byte, error) {
	var err error

	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	ori := make([]byte, len(src))
	copy(ori, src)
	start := bytes.Index(src, []byte(importStartFlag))
	// in case no importStartFlag or importStartFlag exist in the commentFlag
	if start < 0 {
		return nil, nil, nil
	}
	end := bytes.Index(src[start:], []byte(importEndFlag)) + start

	// in case import flags are part of a codegen template, or otherwise "wrong"
	if start+len(importStartFlag) > end {
		return nil, nil, nil
	}

	ret := bytes.Split(src[start+len(importStartFlag):end], []byte(linebreak))

	block := newImportBlock(ret, set.LocalFlag)

	res := append(src[:start+len(importStartFlag)], append(block.fmt(), src[end+1:]...)...)

	if bytes.Equal(ori, res) {
		return ori, nil, nil
	}

	return ori, res, nil
}
