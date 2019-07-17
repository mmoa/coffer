package coffer

// Coffer
// API tests
// Copyright © 2019 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	//"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	// "github.com/claygod/coffer/domain"
	// "github.com/claygod/coffer/services"
	// "github.com/claygod/coffer/services/filenamer"
	// "github.com/claygod/coffer/services/journal"
	// "github.com/claygod/coffer/services/repositories/handlers"
	// "github.com/claygod/coffer/services/repositories/records"
	"github.com/claygod/coffer/services/journal"
	"github.com/claygod/coffer/services/resources"

	// "github.com/claygod/coffer/services/startstop"
	"github.com/claygod/coffer/usecases"
	// "github.com/claygod/tools/logger"
	// "github.com/claygod/tools/porter"
)

func TestNewCoffer(t *testing.T) {
	//defer forTestClearDir("./test/")
	jCnf := &journal.Config{
		BatchSize:              2000,
		LimitRecordsPerLogfile: 100000,
	}
	ucCnf := &usecases.Config{
		FollowPause:             1 * time.Second,
		ChagesByCheckpoint:      100,
		DirPath:                 "./test/", // "/home/ed/goPath/src/github.com/claygod/coffer/test",
		AllowStartupErrLoadLogs: true,
		MaxKeyLength:            100,
		MaxValueLength:          10000,
	}
	rcCnf := &resources.Config{
		LimitMemory: 1000 * megabyte, // minimum available memory (bytes)
		LimitDisk:   1000 * megabyte, // minimum free disk space
		DirPath:     "./test/",       // "/home/ed/goPath/src/github.com/claygod/coffer/test"
	}

	cnf := &Config{
		JournalConfig:       jCnf,
		UsecasesConfig:      ucCnf,
		ResourcesConfig:     rcCnf,
		MaxRecsPerOperation: 100,
		//MaxKeyLength:        100,
		//MaxValueLength:      10000,
	}
	cof, err := New(cnf)
	if err != nil {
		t.Error(err)
		return
	}
	if cof.Start() {
		defer cof.Stop()
	} else {
		t.Errorf("Failed to start")
		return
	}
	for i := 20; i < 30; i++ {
		if err := cof.Write("aasa"+strconv.Itoa(i), []byte("bbsb")); err != nil {
			t.Error(err)
		}
		time.Sleep(900 * time.Millisecond)
	}

}

func forTestClearDir(dir string) error {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		//fmt.Println(name)
		if strings.HasSuffix(name, ".log") || strings.HasSuffix(name, ".check") || strings.HasSuffix(name, ".checkpoint") {
			os.Remove(dir + name)
		}
		//		err = os.RemoveAll(filepath.Join(dir, name))
		//		if err != nil {
		//			return err
		//		}
	}
	return nil
}