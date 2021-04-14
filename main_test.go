package main_test

import (
	"io/ioutil"
	"syscall"
	"testing"

	main "golift.io/turbovanityurls"
)

func TestParseFlags(t *testing.T) {
	t.Parallel()

	test := []string{"-l", "127.0.0.1:456", "-c", "config.file", "-v"}
	flags := main.ParseFlags(test)

	if flags.ListenAddr != test[1] {
		t.Errorf("test flag was not parsed properly: %v", flags.ListenAddr)
	}

	if flags.ConfigPath != test[3] {
		t.Errorf("test flag was not parsed properly: %v", flags.ConfigPath)
	}

	if !flags.ShowVer {
		t.Errorf("test flag was not parsed properly: ShowVer=%v", flags.ShowVer)
	}

	flags = main.ParseFlags([]string{})

	if flags.ListenAddr != ":8080" {
		t.Errorf("default flag value not correct: %v", flags.ListenAddr)
	}

	if flags.ConfigPath != main.DefaultConfFile {
		t.Errorf("default flag value not correct: %v", flags.ConfigPath)
	}
}

//nolint:funlen
func TestParseConfig(t *testing.T) {
	t.Parallel()

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

		_ = ioutil.WriteFile(f.Name(), []byte(test.config), 0600)
		c := &main.Config{}

		err = c.ParseConfig(f.Name())
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

	c := &main.Config{}

	err := c.ParseConfig("missing_file_here.asahsahsahsahs")
	if err == nil {
		t.Errorf("parseConfig must return an error when the config file is missing")
	}

	err = c.ParseConfig("/etc/passwd")
	if err == nil {
		t.Errorf("parseConfig must return n error with an invalid config file")
	}
}
