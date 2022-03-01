package main

import (
	"context"
	"io/ioutil"
	"testing"
)

func TestGraphqlRequest(t *testing.T) {
	data, _ := ioutil.ReadFile("test.NEF")
	img, err := decodeImage("image/nef", data)
	if err != nil {
		t.Error(err)
	}

	r, err := resizeImageThumbnail(context.Background(), img, 1500, 0)
	if err != nil {
		t.Error(err)
	}

	data2, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("test.jpg", data2, 0700)
	if err != nil {
		t.Error(err)
	}
}
