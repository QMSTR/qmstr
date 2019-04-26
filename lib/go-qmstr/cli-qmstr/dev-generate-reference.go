package cliqmstr

import (
	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
)

func init() {
	cmd := common.CreateGenerateReferenceCmd(rootCmd)
	rootCmd.AddCommand(cmd)
}
