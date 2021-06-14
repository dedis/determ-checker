package checker

import (
	"flag"
	"testing"

	"github.com/dedis/deter-checker/checker"
	"github.com/stretchr/testify/require"
)

var pfile string
var tfile string
var sfile string

func init() {
	flag.StringVar(&pfile, "pkg", "", "List of whitelisted packages")
	flag.StringVar(&tfile, "types", "", "List of blacklisted types")
	flag.StringVar(&sfile, "src", "", "Source file")
}

func Test_Checker(t *testing.T) {
	wlistPkg := checker.ReadList(&pfile)
	blistTypes := checker.ReadList(&tfile)
	require.NotNil(t, wlistPkg)
	require.NotNil(t, blistTypes)
	res := checker.AnalyzeSource(&sfile, wlistPkg, blistTypes)
	require.True(t, res)
}
