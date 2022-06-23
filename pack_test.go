package tcpserver

import (
	"bytes"
	"testing"
)

func TestPackTcp(t *testing.T) {
	cases := []struct {
		in   TcpData
		want []byte
	}{
		{&tcpdata{5, []byte("hello")}, []byte{5, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}},
		{&tcpdata{3, []byte("you")}, []byte{3, 0, 0, 0, 'y', 'o', 'u'}},
		{&tcpdata{9, []byte("你好毒")}, []byte{9, 0, 0, 0, 228, 189, 160, 229, 165, 189, 230, 175, 146}},
	}

	p := &packer{}
	for _, v := range cases {
		out, err := p.PackTcp(v.in)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%v", out)
		if bytes.Compare(v.want, out) != 0 {
			t.Fatalf("want: %v, result: %v", v.want, out)
		}
	}
}

func TestUnPackTcp(t *testing.T) {
	cases := []struct {
		want TcpData
		in   []byte
	}{
		{&tcpdata{5, []byte("hello")}, []byte{5, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}},
		{&tcpdata{3, []byte("you")}, []byte{3, 0, 0, 0, 'y', 'o', 'u'}},
		{&tcpdata{9, []byte("你好毒")}, []byte{9, 0, 0, 0, 228, 189, 160, 229, 165, 189, 230, 175, 146}},
	}

	p := &packer{}
	for _, v := range cases {
		out, err := p.UnPackTcp(bytes.NewBuffer(v.in))
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%v, %v", out.Size(), string(out.Data()))
		if v.want.Size() != out.Size() ||
			bytes.Compare(v.want.Data(), out.Data()) != 0 {
			t.Fatalf("want: %v, result: %v", v.want, out)
		}
	}
}

func TestPackMessage(t *testing.T) {
	cases := []struct {
		in   Message
		want []byte
	}{
		{&message{1, []byte("hello")}, []byte{1, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}},
		{&message{2, []byte("you")}, []byte{2, 0, 0, 0, 'y', 'o', 'u'}},
		{&message{3, []byte("你好毒")}, []byte{3, 0, 0, 0, 228, 189, 160, 229, 165, 189, 230, 175, 146}},
	}

	p := &packer{}
	for _, v := range cases {
		out, err := p.PackMessage(v.in)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%v", out)
		if bytes.Compare(v.want, out) != 0 {
			t.Fatalf("want: %v, result: %v", v.want, out)
		}
	}
}

func TestUnPackMessage(t *testing.T) {
	cases := []struct {
		want Message
		in   []byte
	}{
		{&message{1, []byte("hello")}, []byte{1, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}},
		{&message{2, []byte("you")}, []byte{2, 0, 0, 0, 'y', 'o', 'u'}},
		{&message{3, []byte("你好毒")}, []byte{3, 0, 0, 0, 228, 189, 160, 229, 165, 189, 230, 175, 146}},
	}

	p := &packer{}
	for _, v := range cases {
		out, err := p.UnPackMessage(bytes.NewBuffer(v.in))
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%v, %v", out.Id(), string(out.Data()))
		if v.want.Id() != out.Id() ||
			bytes.Compare(v.want.Data(), out.Data()) != 0 {
			t.Fatalf("want: %v, result: %v", v.want, out)
		}
	}
}

func TestPack(t *testing.T) {
	cases := []struct {
		in   Message
		want []byte
	}{
		{&message{1, []byte("hello")}, []byte{9, 0, 0, 0, 1, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}},
		{&message{2, []byte("you")}, []byte{7, 0, 0, 0, 2, 0, 0, 0, 'y', 'o', 'u'}},
		{&message{3, []byte("你好毒")}, []byte{13, 0, 0, 0, 3, 0, 0, 0, 228, 189, 160, 229, 165, 189, 230, 175, 146}},
	}

	p := &packer{}
	for _, v := range cases {
		out, err := p.Pack(v.in)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%v", out)
		if bytes.Compare(v.want, out) != 0 {
			t.Fatalf("want: %v, result: %v", v.want, out)
		}
	}
}

func TestUnPack(t *testing.T) {
	cases := []struct {
		want Message
		in   []byte
	}{
		{&message{1, []byte("hello")}, []byte{9, 0, 0, 0, 1, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}},
		{&message{2, []byte("you")}, []byte{7, 0, 0, 0, 2, 0, 0, 0, 'y', 'o', 'u'}},
		{&message{3, []byte("你好毒")}, []byte{13, 0, 0, 0, 3, 0, 0, 0, 228, 189, 160, 229, 165, 189, 230, 175, 146}},
	}

	p := &packer{}
	for _, v := range cases {
		out, err := p.UnPack(bytes.NewBuffer(v.in))
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%v, %v", out.Id(), string(out.Data()))
		if v.want.Id() != out.Id() ||
			bytes.Compare(v.want.Data(), out.Data()) != 0 {
			t.Fatalf("want: %v, result: %v", v.want, out)
		}
	}
}
