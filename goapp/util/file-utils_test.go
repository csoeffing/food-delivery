package util

import "testing"

func TestGetLastPathComponent(t *testing.T) {
	testGetLastPathComponent(t, "a", "a")
	testGetLastPathComponent(t, "/a/b", "b")
	testGetLastPathComponent(t, "/a/b/c", "c")
}

func testGetLastPathComponent(t *testing.T, path, exp string) {
	act := GetLastPathComponent(path)

	if act != exp {
		t.Fatalf("[a] %s != %s [e]", act, exp)
	}
}

func TestGetFilenameAndExtensions(t *testing.T) {
	testGetFilenameAndExtensions(t, "file.ext", "file", "ext")
	testGetFilenameAndExtensions(t, "file", "file", "")
	testGetFilenameAndExtensions(t, ".ext", "", "ext")
}

func testGetFilenameAndExtensions(t *testing.T, path, expFilename, expExt string) {
	actFilename, actExt := GetFilenameAndExtension(path)

	if actFilename != expFilename {
		t.Fatalf("[a] %s != %s [e]", actFilename, expFilename)
	}

	if actExt != expExt {
		t.Fatalf("[a] %s != %s [e]", actExt, expExt)
	}
}
