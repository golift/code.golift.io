package main

import (
	"io/ioutil"
	"syscall"
	"testing"
)

func TestParseFlags(t *testing.T) {
	test := []string{"-l", "127.0.0.1:456", "-c", "config.file", "-v"}
	flags := parseFlags(test)
	if flags.listenAddr != test[1] {
		t.Errorf("test flag was not parsed properly: %v", flags.listenAddr)
	}
	if flags.configPath != test[3] {
		t.Errorf("test flag was not parsed properly: %v", flags.configPath)
	}
	if !flags.showVer {
		t.Errorf("test flag was not parsed properly: showVer=%v", flags.showVer)
	}
	flags = parseFlags([]string{})
	if flags.listenAddr != ":8080" {
		t.Errorf("default flag value not correct: %v", flags.listenAddr)
	}
	if flags.configPath != "./config.yaml" {
		t.Errorf("default flag value not correct: %v", flags.configPath)
	}
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		host   string
		config string
		title  string
		bdPath string
	}{{
		config: "title: FOO\n" +
			"bd_path: /bd\n" +
			"paths:\n" +
			"  /gopdf:\n" +
			"    repo: https://bitbucket.org/zombiezen/gopdf\n",
		title:  "FOO",
		host:   "",
		bdPath: "/bd/",
	}, {
		config: "host: foo.com\n" +
			"paths:\n" +
			"  /gopdf:\n" +
			"    repo: https://testbucket.org/zombiezen/gopdf\n" +
			"    vcs: git\n",
		title: "foo.com",
		host:  "foo.com",
	}}

	for _, test := range tests {
		f, err := ioutil.TempFile("", "*.conf")
		if err != nil {
			t.Errorf("writing test temporary file failed\n%s", err)
		}
		_ = f.Close()

		defer func() {
			if err := syscall.Unlink(f.Name()); err != nil {
				t.Errorf("error deleting test file\n%v\n%s", err, f.Name())
			}
		}()
		_ = ioutil.WriteFile(f.Name(), []byte(test.config), 0644)
		c, err := parseConfig(f.Name())
		if err != nil {
			t.Errorf("test config produced unexpected error\n%v\n%s", err, test.config)
		}
		if test.title != c.Title {
			t.Errorf("test config produced unexpected title\n%v\n%s", err, c.Title)
		}
		if test.bdPath != c.BDPath {
			t.Errorf("test config produced unexpected bd_path\n%v\n%s", err, c.BDPath)
		}
	}
	_, err := parseConfig("missing_file_here.asahsahsahsahs")
	if err == nil {
		t.Errorf("parseConfig must return an error when the config file is missing")
	}
	_, err = parseConfig("/etc/passwd")
	if err == nil {
		t.Errorf("parseConfig must return n error with an invalid config file")
	}
}
