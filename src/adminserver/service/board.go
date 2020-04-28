package service

import (
	"bytes"
	"fmt"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/tools/secureShell"
	"os"
	"os/exec"
	"path"
)

//Start Board without loading cfg.
func Start(host *models.Account) error {
	shell, err := SSHtoHost(host)
	if err != nil {
		return err
	}

	cmdComposeUp := fmt.Sprintf("docker-compose -f %s up -d", models.Boardcompose)
	err = shell.ExecuteCommand(cmdComposeUp)
	if err != nil {
		return err
	}

	UUIDpath := "/go/secrets/initialAdminPassword"
	if _, err = os.Stat(UUIDpath); !os.IsNotExist(err) {
		if err = dao.RemoveUUIDToken(); err != nil {
			return err
		}
		os.Remove(UUIDpath)
	}

	return nil
}

//Applycfg restarts Board with applying of cfg.
func Applycfg(host *models.Account) error {

	cfgPath := path.Join("/go", "/cfgfile/board.cfg")
	err := os.Rename(cfgPath, cfgPath+".bak1")
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			return err
		}
	}
	err = os.Rename(cfgPath+".tmp", cfgPath)
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			return err
		}
	}
	_, err = Execute(fmt.Sprintf("cp %s %s.tmp", cfgPath, cfgPath))
	if err != nil {
		return err
	}

	if err = StartBoard(host); err != nil {
		return err
	}

	if err = os.Remove(cfgPath + ".tmp"); err != nil {
		return err
	}

	return nil
}

//Shutdown Board.
func Shutdown(host *models.Account, uninstall bool) error {
	shell, err := SSHtoHost(host)
	if err != nil {
		return err
	}

	cmdCompose := fmt.Sprintf("docker-compose -f %s down", models.Boardcompose)
	err = shell.ExecuteCommand(cmdCompose)
	if err != nil {
		return err
	}

	if uninstall {
		cmdRm := fmt.Sprintf("rm -rf /data/board/* %s/board.cfg* && cp %s/adminserver/board.cfg %s/.", models.MakePath, models.MakePath, models.MakePath)
		err = shell.ExecuteCommand(cmdRm)
		if err != nil {
			return err
		}
	}
	return nil
}

func SSHtoHost(host *models.Account) (*secureShell.SecureShell, error) {
	var output bytes.Buffer
	var shell *secureShell.SecureShell

	HostIP, err := Execute("ip route | awk 'NR==1 {print $3}'|xargs echo -n")
	if err != nil {
		return nil, err
	}
	shell, err = secureShell.NewSecureShell(&output, HostIP, host.Username, host.Password)
	if err != nil {
		return nil, err
	}
	return shell, nil
}

//Execute command in container.
func Execute(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
