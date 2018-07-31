package main

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

var loglines string = `type=SYSCALL msg=audit(1532024611.855:7230280): arch=c000003e syscall=60 a0=0 a1=7f32270bb270 a2=7f32259e6700 a3=7f32259e69d0 items=0 ppid=0 pid=1 auid=4294967295 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=(none) ses=4294967295 comm="systemd" exe="/usr/lib/systemd/systemd" key=(null)
type=PROCTITLE msg=audit(1532024611.855:7230280): proctitle=2F7573722F6C69622F73797374656D642F73797374656D64002D2D73776974636865642D726F6F74002D2D73797374656D002D2D646573657269616C697A65003231
type=SYSCALL msg=audit(1532024611.885:7230281): arch=c000003e syscall=231 a0=0 a1=0 a2=0 a3=fffffffffffffe80 items=0 ppid=12987 pid=13248 auid=741 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts0 ses=1373 comm="vim" exe="/usr/bin/vim" key=(null)
type=PROCTITLE msg=audit(1532024611.885:7230281): proctitle=76696D002E2E2F2E2E2F6574632F61756469742E636F6E6669672E6A736F6E
type=SYSCALL msg=audit(1532024613.105:7230282): arch=c000003e syscall=60 a0=0 a1=7f32270bb270 a2=7f32259e6700 a3=7f32259e69d0 items=0 ppid=0 pid=1 auid=4294967295 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=(none) ses=4294967295 comm="systemd" exe="/usr/lib/systemd/systemd" key=(null)
type=PROCTITLE msg=audit(1532024613.105:7230282): proctitle=2F7573722F6C69622F73797374656D642F73797374656D64002D2D73776974636865642D726F6F74002D2D73797374656D002D2D646573657269616C697A65003231
type=SYSCALL msg=audit(1532024614.355:7230283): arch=c000003e syscall=60 a0=0 a1=7f32270bb270 a2=7f32259e6700 a3=7f32259e69d0 items=0 ppid=0 pid=1 auid=4294967295 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=(none) ses=4294967295 comm="systemd" exe="/usr/lib/systemd/systemd" key=(null)
type=PROCTITLE msg=audit(1532024614.355:7230283): proctitle=2F7573722F6C69622F73797374656D642F73797374656D64002D2D73776974636865642D726F6F74002D2D73797374656D002D2D646573657269616C697A65003231`

var expectedMsg1 string = `type=SYSCALL msg=audit(1532024611.855:7230280): arch=c000003e syscall=60 a0=0 a1=7f32270bb270 a2=7f32259e6700 a3=7f32259e69d0 items=0 ppid=0 pid=1 auid=4294967295 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=(none) ses=4294967295 comm="systemd" exe="/usr/lib/systemd/systemd" key=(null)`

var expectedMsg2 string = `type=SYSCALL msg=audit(1532024611.885:7230281): arch=c000003e syscall=231 a0=0 a1=0 a2=0 a3=fffffffffffffe80 items=0 ppid=12987 pid=13248 auid=741 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts0 ses=1373 comm="vim" exe="/usr/bin/vim" key=(null)`

func TestNewFileClient(t *testing.T) {
	_, elb := hookLogger()
	defer resetLogger()

	buf := bytes.NewBufferString(loglines)
	scanner := bufio.NewScanner(buf)
	f, err := NewFileClient(scanner)
	assert.NoError(t, err)
	if f == nil {
		t.Fatal("Expected a file client but had an unknown error instead")
	} else {
		assert.True(t, (f.s != nil), "No scanner")
		assert.Equal(t, "", elb.String(), "Did not expect any error messages")
	}
}

func TestFileClientReceive(t *testing.T) {
	_, elb := hookLogger()
	defer resetLogger()

	buf := bytes.NewBufferString(loglines)
	scanner := bufio.NewScanner(buf)
	f, err := NewFileClient(scanner)
	assert.NoError(t, err)
	if f == nil {
		t.Fatal("Expected a file client but had an unknown error instead")
	} else {
		msg, err := f.Receive()
		assert.NoError(t, err)
		assert.Equal(t, "", elb.String(), "Did not expect any error messages")
		assert.NotNil(t, msg)
		assert.Equal(t, expectedMsg1, string(msg.Data))

		// advance to 3rd message
		msg, err = f.Receive()
		msg, err = f.Receive()
		assert.NoError(t, err)
		assert.Equal(t, "", elb.String(), "Did not expect any error messages")
		assert.Equal(t, expectedMsg2, string(msg.Data))
	}
}
