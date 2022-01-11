package client

import (
	"bytes"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/pingcap/log"
	"github.com/pingcap/parser/ast"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/tikv"
	pd "github.com/tikv/pd/client"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap/zapcore"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"tikv-client/cdc"
	p "tikv-client/pd"
	kv "tikv-client/tikv"
	"tikv-client/util"
	"time"
)

const (
	GoodBye_  = "Goodbye!"
	back_     = "back"
	exit_     = "exit"
	Ticlient_ = "ticlient"
)

type Completer struct {
	GoClient             *tikv.RawKVClient
	EtcdClient           *clientv3.Client
	PlacementDriverGroup *p.PlacementDriverGroup
	cdc.Cdc
	kv.TiKV
	stmtNode []ast.StmtNode
	option   string
}

func NewCompleter(pdEndPoint string) (*Completer, error) {
	// PlacementDriver Group
	placementDriverGroup, err := p.PlacementDriverInfo(pdEndPoint)
	if err != nil {
		return nil, err
	}
	// client-go
	var pdHosts []string
	for _, placementDriver := range placementDriverGroup.Members {
		pdHost := strings.ReplaceAll(placementDriver.ClientUrls[0], "http://", "")
		pdHosts = append(pdHosts, pdHost)
	}
	c, err := newGoClient(pdHosts)
	if err != nil {
		return nil, err
	}
	ec, err := newEtcdClient(pdHosts)
	if err != nil {
		return nil, err
	}
	return &Completer{
		GoClient:             c,
		EtcdClient:           ec,
		PlacementDriverGroup: placementDriverGroup,
	}, nil
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {
	var suggest []prompt.Suggest
	args := specificationArgs(d.TextBeforeCursor())
	// For main menu.
	if c.option == "" && len(args) == 1 {
		suggest = []prompt.Suggest{
			{Text: kv.Tikv_, Description: "This is a tikv client, you can use simple SQL to query data and other operations in tikv."},
			{Text: cdc.Cdc_, Description: "You can use it for operation and maintenance management ticdc."},
			{Text: exit_},
		}
	}
	switch c.option {
	case kv.Tikv_:
		if d.TextBeforeCursor() == "" {
			return suggest
		}
		suggest = []prompt.Suggest{
			{Text: kv.Select_},
			{Text: kv.InsertInto_},
			{Text: kv.Delete_},
			{Text: kv.From_},
			{Text: kv.Where_},
			{Text: kv.Values_},
			{Text: kv.OrderBy_},
			{Text: kv.Limit_},
			{Text: kv.Truncate_},
			{Text: back_},
		}
		if len(args) > 1 {
			switch strings.ToLower(args[len(args)-2]) {
			case kv.From_, kv.Into_:
				suggest = []prompt.Suggest{
					{Text: kv.TikvTable_},
				}
			default:
			}
		}
	case cdc.Cdc_:
		// ticdc tree
		return c.Cdc.CdcCompleter(args, d)
	default:
	}
	return prompt.FilterHasPrefix(suggest, d.GetWordBeforeCursor(), true)
}

func (c *Completer) Executor(cmd string) {

	if strings.TrimSpace(strings.ToLower(cmd)) == exit_ {
		if err := c.GoClient.Close(); err != nil {
			fmt.Printf(util.Red("Exit failed, error: %s\n"), err)
			return
		}
		fmt.Println(GoodBye_)
		syscall.Exit(0)
	}

	cmd = strings.ToLower(strings.TrimSpace(cmd))

	// Note user has in self option.
	if c.inSelf(cmd) {
		return
	}
	if c.mainMenu(cmd) {
		return
	}

	// Run command...
	var err error
	switch c.option {
	case kv.Tikv_:
		err = kv.TiKV{
			PlacementDriverGroup: c.PlacementDriverGroup,
			GoClient:             c.GoClient,
		}.Run(cmd)
	case cdc.Cdc_:
		err = c.Cdc.Run(specificationArgs(cmd))
	default:
		fmt.Println(util.Red(fmt.Sprintf("Invalid command: %s", cmd)))
	}
	if err != nil {
		fmt.Println(util.Red(err.Error()))
	}
}

func (c *Completer) mainMenu(cmd string) bool {
	flag := false
	switch cmd {
	case "":
		flag = true
	case back_:
		// Back means clean the option
		c.cleanOption()
		flag = true
	case cdc.Cdc_:
		// Choose cdc, flash all cdc cluster info first
		Welcome(*c.PlacementDriverGroup, cmd+"-client")
		c.option = cmd
		c.newCdcCluster()
		flag = true
	case kv.Tikv_:
		Welcome(*c.PlacementDriverGroup, cmd+"-client")
		c.option = cmd
		flag = true
	}
	return flag
}

func Welcome(pd p.PlacementDriverGroup, option string) {
	if option == Ticlient_ {
		fmt.Printf("Welcome to %s! Commands `exit` to quit.\n"+"PD version: %s\n\n", util.Red(option), pd.Leader.BinaryVersion)
	} else {
		fmt.Printf("Enter %s!\nCommands `back` to return main menu. \n\n", util.Red(option))
	}
}

func (c *Completer) inSelf(s string) bool {
	if c.option == s && s != "" {
		fmt.Println(util.Red(fmt.Sprintf("You have in %s, please perform the corresponding operation or enter back to return to the main menu", s)))
		return true
	}
	return false
}

func (c *Completer) cleanOption() {
	c.option = ""
	fmt.Println(util.Green("Back to main menu..."))
}

func (c *Completer) newCdcCluster() {
	// Choose cdc, init all cdc info first.
	d, err := cdc.Init(*c.PlacementDriverGroup, c.EtcdClient)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Cdc = *d
}

func newGoClient(pdEndpoints []string) (*tikv.RawKVClient, error) {
	// To prevent go-client from printing logs
	log.SetLevel(zapcore.PanicLevel)
	c, err := tikv.NewRawKVClient(pdEndpoints, config.DefaultConfig().Security, pd.WithMaxErrorRetry(1))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newEtcdClient(pdEndpoints []string) (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints:   pdEndpoints,
		DialTimeout: 5 * time.Second,
	}
	var c *clientv3.Client
	var err error
	if c, err = clientv3.New(cfg); err != nil {
		return nil, err
	}
	return c, nil
}

func specificationArgs(d string) []string {
	args := strings.Split(d, " ")
	for i, v := range args {
		if v == " " {
			args = append(args[:i], args[i+1:]...)
		}
	}
	return args
}

func ExecuteAndGetResult(s string) (string, error) {
	out := &bytes.Buffer{}
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	r := string(out.Bytes())
	return r, nil
}
