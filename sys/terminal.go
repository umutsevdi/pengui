package sys

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func Exec(input string) string {
	cmd := exec.Command(input)
	err := cmd.Run()
	cmd.Stdout

}
