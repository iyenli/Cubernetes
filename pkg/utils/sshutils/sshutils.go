package sshutils

import (
	"bytes"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
)

func Exec(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		log.Println("[FATAL] Fail to create new session on ssh server, err: ", err)
		return "", err
	}
	defer func() { _ = session.Close() }()

	result, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Println("[FATAL] Fail to execute on ssh server, err: ", err)
		return "", err
	}

	return string(result), nil
}

func UploadFile(client *ssh.Client, localFile string, remoteFile string, mkdir string) error {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Println("fail to initialize sftp client, err: ", err)
		return err
	}
	defer func() { _ = sftpClient.Close() }()

	if mkdir != "" {
		_ = sftpClient.Mkdir(mkdir)
	}

	dstFile, err := sftpClient.Create(remoteFile)
	if err != nil {
		log.Println("fail to create file on sftp client, err: ", err)
		return err
	}
	defer func() { _ = dstFile.Close() }()

	buf, err := ioutil.ReadFile(localFile)
	if err != nil {
		log.Println("fail to read local file, err: ", err)
		return err
	}

	_, err = dstFile.Write(buf)
	if err != nil {
		log.Println("fail to write file on sftp client, err: ", err)
		return err
	}

	return nil
}

func ReadFile(client *ssh.Client, filename string) ([]byte, error) {
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Println("fail to initialize sftp client, err: ", err)
		return nil, err
	}
	defer func() { _ = sftpClient.Close() }()

	f, err := sftpClient.Open(filename)
	if err != nil {
		log.Println("fail to open file on sftp client, err: ", err)
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var buf bytes.Buffer
	_, err = f.WriteTo(&buf)
	if err != nil {
		log.Println("fail to read file from sftp client, err: ", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
