package btcwallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenDomido/Bitcoin-Wallet"
	"github.com/OpenDomido/Bitcoin-Wallet/api"
	"github.com/OpenDomido/Bitcoin-Wallet/cli"
	"github.com/OpenDomido/Bitcoin-Wallet/db"
	wi "github.com/OpenDomido/wallet-interface"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/jessevdk/go-flags"
	"github.com/natefinch/lumberjack"
	"github.com/op/go-logging"
	"github.com/yawning/bulb"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path"
	"time"
	//"github.com/btcsuite/btcd/txscript"
	//"github.com/golang/protobuf/ptypes"
	//"encoding/hex"
	"encoding/hex"
	//"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/txscript"
)

var parser = flags.NewParser(nil, flags.Default)

type Start struct {
	DataDir            string `short:"d" long:"datadir" description:"specify the data directory to be used"`
	Testnet            bool   `short:"t" long:"testnet" description:"use the test network"`
	Regtest            bool   `short:"r" long:"regtest" description:"run in regression test mode"`
	Mnemonic           string `short:"m" long:"mnemonic" description:"specify a mnemonic seed to use to derive the keychain"`
	WalletCreationDate string `short:"w" long:"walletcreationdate" description:"specify the date the seed was created. if omitted the wallet will sync from the oldest checkpoint."`
	TrustedPeer        string `short:"i" long:"trustedpeer" description:"specify a single trusted peer to connect to"`
	Tor                bool   `long:"tor" description:"connect via a running Tor daemon"`
	FeeAPI             string `short:"f" long:"feeapi" description:"fee API to use to fetch current fee rates. set as empty string to disable API lookups." default:"https://bitcoinfees.21.co/api/v1/fees/recommended"`
	MaxFee             uint64 `short:"x" long:"maxfee" description:"the fee-per-byte ceiling beyond which fees cannot go" default:"2000"`
	LowDefaultFee      uint64 `short:"e" long:"economicfee" description:"the default low fee-per-byte" default:"140"`
	MediumDefaultFee   uint64 `short:"n" long:"normalfee" description:"the default medium fee-per-byte" default:"160"`
	HighDefaultFee     uint64 `short:"p" long:"priorityfee" description:"the default high fee-per-byte" default:"180"`
	Gui                bool   `long:"gui" description:"launch an experimental GUI"`
}
type Version struct{}

var start Start
var version Version
var wl *btcwallet.SPVWallet

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("SPVwl. shutting down...")
			wl.Close()
			os.Exit(1)
		}
	}()


	fmt.Println(os.Args)
	if len(os.Args) == 1 {
		start.Gui = true
		start.Execute([]string{"defaultSettings"})
	} else {
		parser.AddCommand("start",
			"start the wl.",
			"The start command starts the wl. daemon",
			&start)
		parser.AddCommand("version",
			"print the version number",
			"Print the version number and exit",
			&version)
		cli.SetupCli(parser)
		if _, err := parser.Parse(); err != nil {
			os.Exit(1)
		}
	}
}

func (x *Version) Execute(args []string) error {
	fmt.Println(btcwallet.WALLET_VERSION)
	return nil
}

func (x *Start) Execute(args []string) error {
	var err error
	// Create a new config
	config := btcwallet.NewDefaultConfig()

	basepath := config.RepoPath

	config.Params = &chaincfg.TestNet3Params
	config.RepoPath = path.Join(config.RepoPath, "testnet")


	_, ferr := os.Stat(config.RepoPath)
	if os.IsNotExist(ferr) {
		os.Mkdir(config.RepoPath, os.ModePerm)
	}
	if x.Mnemonic != "" {
		config.Mnemonic = x.Mnemonic
	}
	if x.TrustedPeer != "" {
		addr, err := net.ResolveTCPAddr("tcp", x.TrustedPeer)
		if err != nil {
			return err
		}
		config.TrustedPeer = addr
	}
	if x.Tor {
		var conn *bulb.Conn
		conn, err = bulb.Dial("tcp4", "127.0.0.1:9151")
		if err != nil {
			conn, err = bulb.Dial("tcp4", "127.0.0.1:9151")
			if err != nil {
				return errors.New("Tor daemon not found")
			}
		}
		dialer, err := conn.Dialer(nil)
		if err != nil {
			return err
		}
		config.Proxy = dialer
	}
	if x.FeeAPI != "" {
		u, err := url.Parse(x.FeeAPI)
		if err != nil {
			return err
		}
		config.FeeAPI = *u
	}
	if len(args) == 0 {
		config.MaxFee = x.MaxFee
		config.LowFee = x.LowDefaultFee
		config.MediumFee = x.MediumDefaultFee
		config.HighFee = x.HighDefaultFee
	}

	// Make the logging a little prettier
	var fileLogFormat = logging.MustStringFormatter(`%{time:15:04:05.000} [%{shortfunc}] [%{level}] %{message}`)
	w := &lumberjack.Logger{
		Filename:   path.Join(config.RepoPath, "logs", "bitcoin.log"),
		MaxSize:    10, // Megabytes
		MaxBackups: 3,
		MaxAge:     30, // Days
	}
	bitcoinFile := logging.NewLogBackend(w, "", 0)
	bitcoinFileFormatter := logging.NewBackendFormatter(bitcoinFile, fileLogFormat)
	config.Logger = logging.MultiLogger(logging.MultiLogger(bitcoinFileFormatter))

	// Select wl. datastore
	sqliteDatastore, _ := db.Create(config.RepoPath)
	config.DB = sqliteDatastore

	mn, _ := sqliteDatastore.GetMnemonic()
	if mn != "" {
		config.Mnemonic = mn
	}

	// Write version file
	f, err := os.Create(path.Join(basepath, "version"))
	if err != nil {
		return err
	}
	f.Write([]byte("1"))
	f.Close()

	// Load settings
	type Fees struct {
		Priority uint64 `json:"priority"`
		Normal   uint64 `json:"normal"`
		Economic uint64 `json:"economic"`
		FeeAPI   string `json:"feeAPI"`
	}

	type Settings struct {
		FiatCode      string `json:"fiatCode"`
		FiatSymbol    string `json:"fiatSymbol"`
		FeeLevel      string `json:"feeLevel"`
		SelectBox     string `json:"selectBox"`
		BitcoinUnit   string `json:"bitcoinUnit"`
		DecimalPlaces int    `json:"decimalPlaces"`
		TrustedPeer   string `json:"trustedPeer"`
		Proxy         string `json:"proxy"`
		Fees          Fees   `json:"fees"`
	}

	var settings Settings
	s, err := ioutil.ReadFile(path.Join(basepath, "settings.json"))
	if err != nil {
		settings = Settings{
			FiatCode:      "USD",
			FiatSymbol:    "$",
			FeeLevel:      "priority",
			SelectBox:     "bitcoin",
			BitcoinUnit:   "BTC",
			DecimalPlaces: 5,
			Fees: Fees{
				Priority: config.HighFee,
				Normal:   config.MediumFee,
				Economic: config.LowFee,
				FeeAPI:   config.FeeAPI.String(),
			},
		}
		f, err := os.Create(path.Join(basepath, "settings.json"))
		if err != nil {
			return err
		}
		s, err := json.MarshalIndent(&settings, "", "    ")
		if err != nil {
			return err
		}
		f.Write(s)
		f.Close()
	} else {
		err := json.Unmarshal([]byte(s), &settings)
		if err != nil {
			return err
		}
	}
	if settings.TrustedPeer != "" {
		var tp net.Addr
		tp, err = net.ResolveTCPAddr("tcp", settings.TrustedPeer)
		if err != nil {
			return err
		}
		config.TrustedPeer = tp
	}

	if settings.Proxy != "" {
		tbProxyURL, err := url.Parse("socks5://" + settings.Proxy)
		if err != nil {
			return err
		}
		tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
		if err != nil {
			return err
		}
		config.Proxy = tbDialer
	}
	feeApi, _ := url.Parse(settings.Fees.FeeAPI)
	config.FeeAPI = *feeApi
	config.HighFee = settings.Fees.Priority
	config.MediumFee = settings.Fees.Normal
	config.LowFee = settings.Fees.Economic

	creationDate := time.Time{}
	if x.WalletCreationDate  != "" {
		creationDate, err = time.Parse(time.RFC3339, x.WalletCreationDate)
		if err != nil {
			return errors.New("wl. creation date timestamp must be in RFC3339 format")
		}
	}
	config.CreationDate = creationDate

	// Create the wl.
	wl, err = btcwallet.NewSPVWallet(config)
	if err != nil {
		return err
	}

	if err := sqliteDatastore.SetMnemonic(config.Mnemonic); err != nil {
		return err
	}
	if err := sqliteDatastore.SetCreationDate(config.CreationDate); err != nil {
		return err
	}



	listener := func(tx wi.TransactionCallback) {

		scriptHex := hex.EncodeToString(tx.Outputs[0].ScriptPubKey)
		script, err := hex.DecodeString(scriptHex)
		if err != nil {
			return
		}

		// Extract and print details from the script.
		scriptClass, addresses, reqSigs, err := txscript.ExtractPkScriptAddrs(
			script, &chaincfg.TestNet3Params)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Script Class:", scriptClass)
		fmt.Println("Addresses:", addresses)
		fmt.Println("Required Signatures:", reqSigs)

	}

	wl.AddTransactionListener(listener)

	go api.ServeAPI(wl)

	// Start it!
	printSplashScreen()


	wl.Start()

	return nil
}

func printSplashScreen() {
	fmt.Println("")
	fmt.Println("Bitcoin wl. v=" + btcwallet.WALLET_VERSION + " ...")
	fmt.Println("[Press Ctrl+C to exit]")
}
